//go:build levis
// +build levis

package buttons

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/dumacp/go-levis"
	"github.com/dumacp/go-logs/pkg/logs"
)

const (
	addrSelectPaso  = 0
	addrEnterPaso   = 1
	addrEnterRuta   = 2
	addrEnterDriver = 3

	addrScreenAlarms = 7

	addrNoRoute    = 120
	addrNameRoute  = 100
	addrNoDriver   = 160
	addrNameDriver = 140

	addrAddBright = 21
	addrSubBright = 22
)

func ListenButtons(dev interface{}, ctx actor.Context, mem <-chan *MsgMemory, quit <-chan int) error {

	devv, ok := dev.(levis.Device)
	if !ok {
		return fmt.Errorf("dev is not LEVIS device")
	}

	devv.Conf().SetButtonMem(0, 32)

	if err := devv.AddButton(addrSelectPaso); err != nil {
		return err
	}
	devv.SetIndicator(addrSelectPaso, false)

	if err := devv.AddButton(addrEnterPaso); err != nil {
		return err
	}
	devv.SetIndicator(addrEnterPaso, false)

	if err := devv.AddButton(addrEnterRuta); err != nil {
		return err
	}
	devv.SetIndicator(addrEnterRuta, false)

	if err := devv.AddButton(addrEnterDriver); err != nil {
		return err
	}
	devv.SetIndicator(addrEnterDriver, false)

	if err := devv.AddButton(addrScreenAlarms); err != nil {
		return err
	}
	devv.SetIndicator(addrScreenAlarms, false)

	if err := devv.AddButton(addrAddBright); err != nil {
		return err
	}
	devv.SetIndicator(addrAddBright, false)

	if err := devv.AddButton(addrSubBright); err != nil {
		return err
	}
	devv.SetIndicator(addrSubBright, false)

	ch := devv.ListenButtons()

	go func() {
		// defer close(ch)

		lastStep := time.Now()
		enableStep := time.NewTimer(5 * time.Second)
		activeStep := false

		for {
			select {
			case <-quit:
				logs.LogWarn.Println("ListenButtons Levis device is closed")
				return
			case <-enableStep.C:
				diff := time.Since(lastStep)
				if diff < 3*time.Second && activeStep {
					enableStep.Reset(diff)
					break
				}
				if activeStep {
					fmt.Println("reset addrSelectPaso")
					activeStep = false
					if err := devv.SetIndicator(addrSelectPaso, false); err != nil {
						fmt.Println(err)
					}
				}
			case button, ok := <-ch:
				if !ok {
					return
				}
				fmt.Printf("button: %v\n", button)
				switch button.Addr {
				case addrSelectPaso:
					if !enableStep.Stop() {
						select {
						case <-enableStep.C:
						default:
						}
					}
					if button.Value == 0 {
						activeStep = false
						break
					}
					if activeStep {
						activeStep = false
						if err := devv.SetIndicator(addrSelectPaso, false); err != nil {
							fmt.Println(err)
						}
						break
					}
					enableStep.Reset(10 * time.Second)
					activeStep = true
					ctx.Send(ctx.Self(), &MsgSelectPaso{})
				case addrEnterPaso:
					if button.Value == 0 {
						break
					}
					if !activeStep {
						break
					}
					if time.Since(lastStep) < 300*time.Millisecond {
						break
					}
					lastStep = time.Now()
					ctx.Send(ctx.Self(), &MsgEnterPaso{})
					go func() {
						<-time.After(2 * time.Second)
					}()
				case addrEnterRuta:
					if button.Value == 0 {
						break
					}
					if err := devv.SetIndicator(addrEnterRuta, false); err != nil {
						fmt.Println(err)
						break
					}
					if route, err := route(devv); err != nil {
						logs.LogWarn.Println(err)
					} else {
						ctx.Send(ctx.Self(), &MsgEnterRuta{Route: route})
					}
				case addrEnterDriver:
					if button.Value == 0 {
						break
					}
					if err := devv.SetIndicator(addrEnterDriver, false); err != nil {
						fmt.Println(err)
						break
					}
					if route, err := driver(devv); err != nil {
						logs.LogWarn.Println(err)
					} else {
						ctx.Send(ctx.Self(), &MsgEnterDriver{Driver: route})
					}
				case addrScreenAlarms:
					if button.Value == 0 {
						break
					}
					if err := devv.SetIndicator(addrScreenAlarms, false); err != nil {
						fmt.Println(err)
						break
					}
					ctx.Send(ctx.Self(), &MsgShowAlarms{})
				case addrAddBright:
					if button.Value == 0 {
						break
					}
					if err := devv.SetIndicator(addrAddBright, false); err != nil {
						fmt.Println(err)
						break
					}
					ctx.Send(ctx.Self(), &MsgBrightnessPlus{})
				case addrSubBright:
					if button.Value == 0 {
						break
					}
					if err := devv.SetIndicator(addrSubBright, false); err != nil {
						fmt.Println(err)
						break
					}
					ctx.Send(ctx.Self(), &MsgBrightnessMinus{})
				}
			}
		}
	}()
	return nil
}
func driver(dev levis.Device) (int, error) {
	res, err := dev.ReadBytesRegister(addrNoDriver, 10)
	if err != nil {
		// fmt.Println(err)
		return -1, err
	}

	fmt.Printf("debug driver: %v\n", res)
	fmt.Printf("driver: %s\n", levis.EncodeToChars(res))

	if len(res) <= 0 {
		return -1, fmt.Errorf("response is empty")
	}

	driverID, err := strconv.Atoi(strings.ReplaceAll(fmt.Sprintf("%s", levis.EncodeToChars(res)), "\x00", ""))
	if err != nil {
		fmt.Println(err)

	}
	return driverID, nil

}

func route(dev levis.Device) (int, error) {
	res, err := dev.ReadBytesRegister(addrNoRoute, 2)
	if err != nil {
		// fmt.Println(err)
		return -1, err
	}

	fmt.Printf("debug route: %v\n", res)

	if len(res) <= 0 {
		return -1, fmt.Errorf("response is empty")
	}

	routeID, _ := strconv.Atoi(strings.ReplaceAll(fmt.Sprintf("%s", levis.EncodeToChars(res)), "\x00", ""))

	return routeID, nil

}
