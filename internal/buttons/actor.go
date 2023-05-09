package buttons

import (
	"context"
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/eventstream"
	"github.com/dumacp/go-driverconsole/internal/device"
	"github.com/dumacp/go-logs/pkg/logs"
)

type Actor struct {
	evts   *eventstream.EventStream
	dev    device.Device
	device ButtonDevice
	// mem       chan *MsgMemory
	pidDevice *actor.PID
	cancel    func()
}

func NewActor(dev ButtonDevice) actor.Actor {
	a := &Actor{}
	a.evts = eventstream.NewEventStream()
	a.device = dev
	// a.mem = make(chan *MsgMemory)
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
		if a.cancel != nil {
			a.cancel()
		}
	case *actor.Started:
		logs.LogInfo.Printf("started \"%s\", %v", ctx.Self().GetId(), ctx.Self())
	case *MsgSubscribe:
		if ctx.Sender() == nil {
			break
		}
		subscribe(ctx, a.evts)
	case *device.MsgDevice:
		contxt, cancel := context.WithCancel(context.TODO())
		a.cancel = cancel

		fmt.Println("///////////////////////")

		if a.device == nil {
			break
		}
		if err := a.device.Init(msg.Device); err != nil {
			logs.LogError.Printf("listenButtons error = %s", err)
			break
		}
		ch, err := a.device.ListenButtons(contxt)
		if err != nil {
			logs.LogError.Printf("listenButtons error = %s", err)
			break
		}

		go func(ctx actor.Context) {
			rootctx := ctx.ActorSystem().Root
			self := ctx.Self()

			for v := range ch {
				rootctx.Request(self, v)
			}
			logs.LogError.Printf("listenButtons close")
		}(ctx)

		// a.pidDevice = ctx.Sender()
	case *InputEvent:
		newmsg := *msg
		a.evts.Publish(&newmsg)
		if ctx.Parent() != nil {
			ctx.Send(ctx.Parent(), &newmsg)
		}
	}
}
