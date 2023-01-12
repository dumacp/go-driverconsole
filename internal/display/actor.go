package display

import (
	"container/list"
	"fmt"
	"strings"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/dumacp/go-driverconsole/internal/app"
	"github.com/dumacp/go-driverconsole/internal/counterpass"
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
	behavior          actor.Behavior
	queueMsgs         *list.List
	brightness        int
	quit              chan int
}

func NewActor() actor.Actor {
	a := &actorDisplay{}
	a.behavior = make(actor.Behavior, 0)
	a.behavior.Become(a.InitState)

	return a

}

func (a *actorDisplay) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *app.MsgBrightnessMinus:
		if a.brightness >= 20 {
			a.brightness -= 10
			a.dev.setBrightness(a.brightness)
		}
	case *app.MsgBrightnessPlus:
		if a.brightness <= 90 {
			a.brightness += 10
			a.dev.setBrightness(a.brightness)
		}
	default:
		a.behavior.Receive(ctx)
	}
}

func (a *actorDisplay) InitState(ctx actor.Context) {
	fmt.Printf("message (InitState): %q --> %q, %T\n", func() string {
		if ctx.Sender() == nil {
			return ""
		} else {
			return ctx.Sender().GetId()
		}
	}(), ctx.Self().GetId(), ctx.Message())
	switch msg := ctx.Message().(type) {
	case *actor.Stopping:
		if a.stopTimeDate != nil {
			select {
			case <-a.stopTimeDate:
			default:
				close(a.stopTimeDate)
			}
		}
		if a.quit != nil {
			select {
			case _, ok := <-a.quit:
				if ok {
					close(a.quit)
				}
			default:
				close(a.quit)
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
		a.queueMsgs = list.New()
	case *device.MsgDevice:
		a.brightness = 50
		display, err := NewDisplay(msg.Device)
		if err != nil {
			logs.LogError.Printf("newDisplay error = %s", err)
			break
		}
		fmt.Println("////////// 0")
		a.dev = display
		a.dev.init()
		fmt.Println("////////// 1")
		a.dev.mainScreen()
		fmt.Println("////////// 1")
		for {
			v := a.queueMsgs.Front()
			if v == nil {
				break
			}
			ctx.Send(ctx.Self(), v.Value)
			a.queueMsgs.Remove(v)
		}
		fmt.Println("////////// 2")

		if a.quit != nil {
			select {
			case _, ok := <-a.quit:
				if ok {
					close(a.quit)
				}
			default:
				close(a.quit)
			}
		}

		a.quit = make(chan int)

		a.dev.verifyReset(a.quit, ctx)
		fmt.Println("////////// 3")
		a.behavior.Become(a.Runstate)
	default:
		a.queueMsgs.PushBack(msg)
	}

}

func (a *actorDisplay) Runstate(ctx actor.Context) {
	a.ctx = ctx
	fmt.Printf("message (Runstate): %q --> %q, %T\n", func() string {
		if ctx.Sender() == nil {
			return ""
		} else {
			return ctx.Sender().GetId()
		}
	}(), ctx.Self().GetId(), ctx.Message())
	switch msg := ctx.Message().(type) {
	case *actor.Stopping:
		if a.stopTimeDate != nil {
			select {
			case <-a.stopTimeDate:
			default:
				close(a.stopTimeDate)
			}
		}
		if a.quit != nil {
			select {
			case _, ok := <-a.quit:
				if ok {
					close(a.quit)
				}
			default:
				close(a.quit)
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
		a.brightness = 50
		display, err := NewDisplay(msg.Device)
		if err != nil {
			logs.LogError.Printf("newDisplay error = %s", err)
			break
		}
		a.dev = display
		a.dev.init()
		a.dev.mainScreen()

		if a.quit != nil {
			select {
			case _, ok := <-a.quit:
				if ok {
					close(a.quit)
				}
			default:
				close(a.quit)
			}
		}
		a.quit = make(chan int)

		a.dev.verifyReset(a.quit, ctx)
	case *Reset:
		a.brightness = 50
		a.dev.mainScreen()
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
		if err := a.dev.ingresos(msg.Efectivo, msg.App, msg.Parcial); err != nil {
			fmt.Println(err)
		}

	case *app.MsgScreen:
		if err := a.dev.switchScreen(msg.ID, msg.Switch); err != nil {
			logs.LogWarn.Println(err)
		}
	case *app.MsgChangeRoute:
		if msg.ID > 0 {
			a.dev.inputValue(fmt.Sprintf("%d", msg.ID), SCREEN_INPUT_ROUTE)
		} else {
			a.dev.inputValue("", SCREEN_INPUT_ROUTE)
		}
	case *app.MsgChangeDriver:
		if msg.ID > 0 {
			a.dev.inputValue(fmt.Sprintf("%d", msg.ID), SCREEN_INPUT_DRIVER)
		} else {
			a.dev.inputValue("", SCREEN_INPUT_DRIVER)
		}
	case *app.MsgDriver:
		text := fmt.Sprintf(`Conductor:
%s`, msg.Driver)
		if err := a.dev.textConfirmation(text); err != nil {
			fmt.Println(err)
		}
		if err := a.dev.driver(msg.Driver); err != nil {
			logs.LogWarn.Println(err)
		}
	case *app.MsgSetRoute:
		if err := a.dev.textConfirmation(fmt.Sprintf(`Ruta:

%s`, msg.Route)); err != nil {
			fmt.Println(err)
		}
		if err := a.dev.route(msg.Route); err != nil {
			logs.LogWarn.Println(err)
		}
		// a.dev.mainScreen()
	case *app.MsgRoute:
		a.dev.textConfirmationMainScreen(3*time.Second, fmt.Sprintf(`Ruta:

%s`, msg.Route))
		if err := a.dev.route(msg.Route); err != nil {
			logs.LogWarn.Println(err)
		}
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
		if a.dev == nil {
			break
		}
		for range []int{1, 2, 3} {
			if err := a.dev.mainScreen(); err != nil {
				continue
			}
			break
		}
	case *UpdateDate:
		if a.dev == nil {
			break
		}
		a.dev.updateDate(0)
	case *Route:
		a.dev.updateRuta(msg.Route, msg.Itininerary)
	case *TimeLapse:
		a.dev.timeRecorrido(msg.Data)
	case *DisplayCount:
		a.dev.ingresos(msg.CountManual, msg.CountAppFare, msg.CountParcial)
	case *DisplayError:
		for range []int{1, 2, 3} {
			if err := a.dev.screenError(msg.Data); err != nil {
				continue
			}
			break
		}
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
		for range []int{1, 2, 3} {
			if err := a.dev.textConfirmation(string(msg.Text)); err != nil {
				fmt.Println(err)
				continue
			}
			break
		}
	case *app.MsgConfirmationTextMainScreen:
		a.dev.textConfirmationMainScreen(msg.Timeout, string(msg.Text))
	case *app.MsgWarningTextInMainScreen:
		a.dev.warningInMainScreen(msg.Timeout, string(msg.Text))
	case *app.MsgWarningText:
		for range []int{1, 2, 3} {
			if err := a.dev.textError(string(msg.Text)); err != nil {
				continue
			}
			break
		}
	case *counterpass.CounterMap:
		a.dev.counters(msg.Inputs0, msg.Outputs0, msg.Inputs1, msg.Outputs1)
	case *counterpass.CounterEvent:
		a.dev.eventCount(msg.Inputs, msg.Outputs)
	case *DisplayDeviceError:
		a.brightness = 50
		for range []int{1, 2, 3} {
			if err := a.dev.mainScreen(); err != nil {
				continue
			}
			break
		}
	case *app.MsgGpsDown:
		a.dev.gpsstate(0)
	case *app.MsgGpsUP:
		a.dev.gpsstate(1)
	case *app.MsgNetDown:
		a.dev.netstate(0)
	case *app.MsgNetUP:
		a.dev.netstate(1)
	case *app.MsgAddAlarm:
		if err := a.dev.addnotification(msg.Data); err != nil {
			fmt.Println(err)
		}
	case *app.MsgShowAlarms:
		if err := a.dev.shownotifications(); err != nil {
			fmt.Println(err)
		}
	}
}

func initTimeDate(ctx actor.Context, stop <-chan int) {

	go func(ctx *actor.RootContext, self *actor.PID) {
		tick := time.NewTicker(5 * time.Second)
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
