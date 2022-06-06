package display

import (
	"fmt"
	"strings"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/dumacp/go-driverconsole/internal/app"
	"github.com/dumacp/go-driverconsole/internal/device"
	"github.com/dumacp/go-logs/pkg/logs"
)

type actorDisplay struct {
	dev              Display
	countUsosParcial int
	countLapse       int
	// varInputs        map[string]string
	cacheTextInput    string
	stopTimeRecorrido chan int
	stopTimeDate      chan int
	ctx               actor.Context
}

func NewActor() actor.Actor {
	a := &actorDisplay{}

	return a
}

func (a *actorDisplay) Receive(ctx actor.Context) {
	a.ctx = ctx
	fmt.Printf("message -> \"%s\", %T\n", ctx.Self().GetId(), ctx.Message())
	switch msg := ctx.Message().(type) {
	case *actor.Stopping:
		if a.stopTimeDate != nil {
			select {
			case <-a.stopTimeDate:
			default:
				close(a.stopTimeDate)
			}
		}
	case *actor.Started:
		if a.stopTimeDate != nil {
			select {
			case _, ok := <-a.stopTimeDate:
				if ok {
					close(a.stopTimeDate)
				}
			default:
				close(a.stopTimeDate)
			}
			time.Sleep(100 * time.Millisecond)
		}
		a.stopTimeDate = make(chan int)
		initTimeDate(ctx, a.stopTimeDate)
	case *device.MsgDevice:
		display, err := NewDisplay(msg.Device)
		if err != nil {
			logs.LogError.Println(err)
			break
		}
		a.dev = display
		a.dev.init()
		a.dev.mainScreen()
	case *Reset:
		a.dev.mainScreen()
		// a.dev.init()
		// a.dev.ingresos(msg.CountManual, msg.CountAppFare, msg.CountParcial)
		// a.dev.updateRuta(msg.Route, msg.Itininerary)
	case *ResetCounter:
		a.countUsosParcial = 0
		a.dev.ingresosPartial(0)
	case *app.MsgAppPercentRecorrido:
		if err := a.dev.recorridoPercent(msg.Data); err != nil {
			logs.LogWarn.Println(err)
		}
	case *app.MsgDoors:
		if err := a.dev.doors(msg.Value); err != nil {
			logs.LogWarn.Println(err)
		}
	case *app.MsgCounters:
		fmt.Printf("message -> \"%v\"\n", msg)
		a.dev.ingresos(msg.Efectivo, msg.App, msg.Parcial)

	case *app.MsgScreen:
		if err := a.dev.switchScreen(msg.ID, msg.Switch); err != nil {
			logs.LogWarn.Println(err)
		}
	case *app.MsgChangeRoute:
		if msg.ID > 0 {
			a.dev.keyNum(fmt.Sprintf("Ruta: %d", msg.ID))
		} else {
			a.dev.keyNum("Ruta: ")
		}
	case *app.MsgChangeDriver:
		if msg.ID > 0 {
			a.dev.keyNum(fmt.Sprintf("Conductor: %d", msg.ID))
		} else {
			a.dev.keyNum("Conductor: ")
		}
	case *app.MsgDriver:
		text := fmt.Sprintf("\n\n\n\n Conductor: %s", msg.Driver)
		a.dev.textConfirmation(text)
	case *app.MsgRoute:
		a.dev.textConfirmation(fmt.Sprintf(`



Ruta: %s`, msg.Route))
		if err := a.dev.route(msg.Route); err != nil {
			logs.LogWarn.Println(err)
		}
		// a.dev.mainScreen()
	case *app.MsgStopRecorrido:

		if a.stopTimeRecorrido != nil {
			select {
			case _, ok := <-a.stopTimeRecorrido:
				if ok {
					close(a.stopTimeRecorrido)
				}
			default:
				close(a.stopTimeRecorrido)
			}
			time.Sleep(300 * time.Millisecond)
		}

	case *app.MsgInitRecorrido:

		if a.stopTimeRecorrido != nil {
			select {
			case _, ok := <-a.stopTimeRecorrido:
				if ok {
					close(a.stopTimeRecorrido)
				}
			default:
				close(a.stopTimeRecorrido)
			}
			time.Sleep(300 * time.Millisecond)
		}

		a.stopTimeRecorrido = make(chan int)

		a.countUsosParcial = 0
		a.dev.ingresosPartial(0)
		a.countLapse = 0
		// a.dev.textInput("Ingrese el ID de conductor: ")

		go func() {
			timeLapse := 0
			a.dev.timeRecorrido(timeLapse)
			tick := time.NewTicker(60 * time.Second)
			defer tick.Stop()
			for {
				select {
				case <-a.stopTimeRecorrido:
					fmt.Println("stop timeLapse")
					a.dev.timeRecorrido(-1)
					return
				case <-tick.C:
					timeLapse++
					a.dev.timeRecorrido(timeLapse)
				}
			}

		}()
	case *app.MsgMainScreen:
		a.dev.mainScreen()
		// a.dev.timeRecorrido(msg.timeLapse)
		// a.dev.ingresos(msg.CountManual, msg.CountManual, msg.CountParcial)
		// a.dev.updateRuta(msg.Route, msg.Itininerary)
	case *UpdateDate:
		a.dev.updateDate(0)
	case *Route:
		a.dev.updateRuta(msg.Route, msg.Itininerary)
	case *TimeLapse:
		a.dev.timeRecorrido(msg.Data)
	case *DisplayCount:
		a.dev.ingresos(msg.CountManual, msg.CountAppFare, msg.CountParcial)
	case *DisplayError:
		a.dev.screenError(msg.Data)
	case *UpVoc:
	case *DisableSelectPASO:
	case *EnterPASO:
	case *EnterVoc:
	case *DelText:
		if len(cacheTextInput) > msg.Count {
			a.cacheTextInput = cacheTextInput[0 : len(cacheTextInput)-msg.Count]
		} else {
			a.cacheTextInput = ""
		}
		a.dev.textInput(cacheTextInput)
	case *AddText:
		a.cacheTextInput = strings.Join([]string{cacheTextInput, msg.Text}, "")
		a.dev.addTextInput(a.cacheTextInput)
	case *EnterText:
	case *app.MsgConfirmationText:
		a.dev.textConfirmation(string(msg.Text))
	case *app.MsgWarningText:
		a.dev.textError(string(msg.Text))
	case *app.MsgConfirmationButton, *app.MsgWarningButton:
		a.dev.mainScreen()
	}
}

func initTimeDate(ctx actor.Context, stop <-chan int) {

	go func(ctx *actor.RootContext, self *actor.PID) {
		tick := time.NewTicker(30 * time.Second)
		defer tick.Stop()
		for {

			select {
			case <-stop:
				return
			case <-tick.C:
				ctx.Send(self, &UpdateDate{})
			}
		}
	}(ctx.ActorSystem().Root, ctx.Self())
}
