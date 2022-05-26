package buttons

import (
	"fmt"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/eventstream"
	"github.com/dumacp/go-driverconsole/internal/device"
	"github.com/dumacp/go-logs/pkg/logs"
)

type Actor struct {
	// portSerialGtt   string
	// portSerialSpeed int
	// display         *gtt43a.Display
	quit      chan int
	evts      *eventstream.EventStream
	device    interface{}
	updateMem bool
	mem       chan *MsgMemory
}

func NewActor() actor.Actor {
	a := &Actor{}
	a.evts = eventstream.NewEventStream()
	a.mem = make(chan *MsgMemory, 0)
	return a
}

func subscribe(ctx actor.Context, evs *eventstream.EventStream) {
	rootctx := ctx.ActorSystem().Root
	pid := ctx.Sender()
	self := ctx.Self()

	fn := func(evt interface{}) {
		rootctx.RequestWithCustomSender(pid, evt, self)
	}
	evs.SubscribeWithPredicate(fn, func(evt interface{}) bool {
		switch evt.(type) {
		case *MsgEnterPaso, *MsgInitRecorrido,
			*MsgStopRecorrido, *MsgEnterRuta, *MsgMainScreen,
			*MsgChangeRuta, *MsgConfirmation, *MsgWarning:
			return true
		}
		return false
	})
}

func (a *Actor) Receive(ctx actor.Context) {
	fmt.Printf("message -> \"%s\", %T,\n", ctx.Self().GetId(), ctx.Message())
	switch msg := ctx.Message().(type) {
	case *actor.Stopping:
		select {
		case <-a.quit:
		default:
			if a.quit != nil {
				close(a.quit)
			}
		}
	case *actor.Started:

	// config := &gtt43a.PortOptions{Port: a.portSerialGtt, Baud: a.portSerialSpeed}
	// gtt := gtt43a.NewDisplay(config)
	// if ok := gtt.Open(); !ok {
	// 	logs.LogError.Fatal("Not connection to display")
	// }
	case *MsgSubscribe:
		if ctx.Sender() == nil {
			break
		}
		subscribe(ctx, a.evts)
	case *device.MsgDevice:
		if a.quit != nil {
			select {
			case _, ok := <-a.quit:
				if ok {
					close(a.quit)
				}
			default:
				close(a.quit)
			}
			time.Sleep(300 * time.Millisecond)
		}

		a.quit = make(chan int)

		if err := ListenButtons(msg.Device, ctx, a.mem, a.quit); err != nil {
			logs.LogError.Println(err)
			break
		}
		a.device = msg.Device
	case *MsgInitRecorrido, *MsgStopRecorrido:

		a.evts.Publish(msg)
	case *MsgChangeRuta:
		if a.evts == nil {
			break
		}
		a.evts.Publish(msg)
	case *MsgMainScreen, *MsgConfirmation, *MsgWarning:
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
	case *MsgSelectPaso:
	case *MsgEnterPaso:
		if a.evts == nil {
			break
		}
		a.evts.Publish(msg)
	case *MsgFatal:
		ctx.Poison(ctx.Self())
	case *MsgMemory:
		select {
		case a.mem <- msg:
		case <-time.After(60 * time.Millisecond):
		}
	}
}
