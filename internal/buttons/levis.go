//+build levis

package buttons

import (
	"fmt"
	"strconv"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/dumacp/go-levis"
	"github.com/dumacp/go-logs/pkg/logs"
)

const (
	addrSelectPaso     = 1
	addrEnterPaso      = 2
	addrEnterRuta      = 3
	addrResetRecorrido = 4

	addrNoRoute   = 110
	addrNameRoute = 20
)

func ListenButtons(dev interface{}, ctx actor.Context, quit <-chan int) error {

	devv, ok := dev.(levis.Device)
	if !ok {
		return fmt.Errorf("dev is not LEVIS device")
	}

	devv.Conf().SetButtonMem(0, 15)

	if err := devv.AddButton(addrSelectPaso); err != nil {
		return err
	}
	devv.SetIndicator(addrSelectPaso, false)
	if err := devv.AddButton(addrEnterPaso); err != nil {
		return err
	}
	devv.SetIndicator(addrEnterPaso, false)
	if err := devv.AddButton(addrResetRecorrido); err != nil {
		return err
	}
	devv.SetIndicator(addrResetRecorrido, false)
	if err := devv.AddButton(addrEnterRuta); err != nil {
		return err
	}
	devv.SetIndicator(addrEnterRuta, false)

	ch := devv.ListenButtons()

	go func() {
		// defer close(ch)

		for {
			select {
			case <-quit:
				logs.LogWarn.Println("ListenButtons Levis device is closed")
				return
			case button := <-ch:
				switch button.Addr {
				case addrSelectPaso:
					if button.Value != 0 {
						ctx.Send(ctx.Self(), &MsgSelectPaso{})
					}
				case addrEnterPaso:
					if button.Value == 0 {
						ctx.Send(ctx.Self(), &MsgEnterPaso{})
					}
				case addrResetRecorrido:
					if button.Value != 0 {
						ctx.Send(ctx.Self(), &MsgInitRecorrido{})
					} else {
						ctx.Send(ctx.Self(), &MsgStopRecorrido{})
					}
				case addrEnterRuta:
					if button.Value != 0 {
						if err := devv.SetIndicator(addrEnterRuta, false); err != nil {
							fmt.Println(err)
							break
						}
						if route, err := route(devv); err != nil {
							logs.LogWarn.Println(err)
						} else {
							ctx.Send(ctx.Self(), &MsgEnterRuta{Route: route})
						}
					}
				}
			}
		}
	}()
	return nil
}

func route(dev levis.Device) (int, error) {
	res, err := dev.ReadBytesRegister(addrNoRoute, 1)
	if err != nil {
		// fmt.Println(err)
		return -1, err
	}

	fmt.Printf("debug route: %v\n", res)

	if len(res) <= 0 {
		return -1, fmt.Errorf("response is empty")
	}

	routeID, _ := strconv.Atoi(fmt.Sprintf("%s", levis.EncodeToChars(res)))

	return routeID, nil

}
