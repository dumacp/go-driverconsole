package buttons

import (
	"fmt"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/eventstream"
	"github.com/dumacp/go-driverconsole/internal/device"
	"github.com/dumacp/go-logs/pkg/logs"
)

type Actor struct {
	quit      chan int
	evts      *eventstream.EventStream
	device    interface{}
	mem       chan *MsgMemory
	pidDevice *actor.PID
}

func NewActor() actor.Actor {
	a := &Actor{}
	a.evts = eventstream.NewEventStream()
	a.mem = make(chan *MsgMemory)
	return a
}

func subscribe(ctx actor.Context, evs *eventstream.EventStream) {
	rootctx := ctx.ActorSystem().Root
	pid := ctx.Sender()
	self := ctx.Self()

	fn := func(evt interface{}) {
		rootctx.RequestWithCustomSender(pid, evt, self)
	}
	evs.Subscribe(fn)
}

func (a *Actor) Receive(ctx actor.Context) {
	fmt.Printf("message: %q --> %q, %T\n", func() string {
		if ctx.Sender() == nil {
			return ""
		} else {
			return ctx.Sender().GetId()
		}
	}(), ctx.Self().GetId(), ctx.Message())
	switch msg := ctx.Message().(type) {
	case *actor.Stopping:
		if a.quit != nil {
			select {
			case <-a.quit:
			default:
				close(a.quit)
			}
		}
	case *actor.Started:
		logs.LogInfo.Printf("started \"%s\", %v", ctx.Self().GetId(), ctx.Self())
	case *MsgSubscribe:
		if ctx.Sender() == nil {
			break
		}
		subscribe(ctx, a.evts)
	case *device.MsgDevice:
		if a.quit != nil {
			select {
			case <-a.quit:
			default:
				close(a.quit)
			}
			time.Sleep(300 * time.Millisecond)
		}

		a.quit = make(chan int)

		if err := ListenButtons(msg.Device, ctx, a.mem, a.quit); err != nil {
			logs.LogError.Printf("listenButtons error = %s", err)
			break
		}
		a.device = msg.Device
		a.pidDevice = ctx.Sender()

	case *MsgInitRecorrido, *MsgStopRecorrido:

		a.evts.Publish(msg)
	case *MsgChangeRuta:
		if a.evts == nil {
			break
		}
		a.evts.Publish(msg)
	case *MsgChangeDriver:
		if a.evts == nil {
			break
		}
		a.evts.Publish(msg)
	case *MsgMainScreen, *MsgConfirmation, *MsgWarning, *MsgReturnFromAlarms, *MsgReturnFromVehicle:
		if a.evts == nil {
			break
		}
		a.evts.Publish(msg)
	case *MsgEnterRuta:
		fmt.Printf("message -> \"%v\"\n", msg)
		if a.evts == nil {
			break
		}
		a.evts.Publish(msg)
	case *MsgEnterDriver:
		fmt.Printf("message -> \"%v\"\n", msg)
		if a.evts == nil {
			break
		}
		a.evts.Publish(msg)
	case *MsgSelectPaso:
	case *MsgEnterPaso:
		if a.evts == nil {
			break
		}
		a.evts.Publish(msg)
	case *MsgEnableEnterPaso:
	case *MsgFatal:
		ctx.Poison(ctx.Self())
	case *MsgDeviceError:
		if a.pidDevice != nil {
			ctx.Request(a.pidDevice, &device.StopDevice{})
			ctx.Request(a.pidDevice, &device.StartDevice{})
		}
	case *MsgMemory:
		select {
		case a.mem <- msg:
		case <-time.After(60 * time.Millisecond):
		}
	case *MsgShowAlarms:
		if a.evts != nil {
			a.evts.Publish(msg)
		}
	case *MsgBrightnessMinus:
		if a.evts != nil {
			a.evts.Publish(msg)
		}
	case *MsgBrightnessPlus:
		if a.evts != nil {
			a.evts.Publish(msg)
		}
	}
}
