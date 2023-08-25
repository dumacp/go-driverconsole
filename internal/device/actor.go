package device

import (
	"context"
	"fmt"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/eventstream"
	"github.com/dumacp/go-logs/pkg/logs"
	"github.com/looplab/fsm"
)

type Actor struct {
	// TODO: ctx???
	ctx       actor.Context
	fmachinae *fsm.FSM
	evts      *eventstream.EventStream
	dev       Device
	contxt    context.Context
}

func NewActor(dev Device) actor.Actor {

	a := &Actor{}
	a.contxt = context.TODO()
	a.dev = dev
	a.evts = &eventstream.EventStream{}
	a.Fsm()
	return a
}

func subscribe(ctx actor.Context, evs *eventstream.EventStream) {
	rootctx := ctx.ActorSystem().Root
	pid := ctx.Sender()
	self := ctx.Self()

	fn := func(evt interface{}) {
		rootctx.RequestWithCustomSender(pid, evt, self)
	}
	evs.SubscribeWithPredicate(fn,
		func(evt interface{}) bool {
			switch evt.(type) {
			case *MsgDevice:
				return true
			}
			return false
		})

}

func (a *Actor) Receive(ctx actor.Context) {
	fmt.Printf("message: %q --> %q, %T\n", func() string {
		if ctx.Sender() == nil {
			return ""
		} else {
			return ctx.Sender().GetId()
		}
	}(), ctx.Self().GetId(), ctx.Message())
	a.ctx = ctx

	switch msg := ctx.Message().(type) {
	case *actor.Started:
		ctx.Send(ctx.Self(), &StartDevice{})
	case *actor.Stopping:
		a.fmachinae.Event(a.contxt, eError)
	case *StartDevice:
		if err := a.fmachinae.Event(a.contxt, eStarted); err != nil {
			logs.LogError.Printf("open device error: %s", err)
			time.Sleep(3 * time.Second)
			ctx.Send(ctx.Self(), &StartDevice{})
		}
		fmt.Printf("open device successfully\n")
	case *MsgDevice:
		a.fmachinae.Event(a.contxt, eOpenned)
		a.evts.Publish(msg)
		if ctx.Parent() != nil {
			ctx.Request(ctx.Parent(), msg)
		}

	case *StopDevice:
		a.fmachinae.Event(a.contxt, eClosed)
		a.fmachinae.Event(a.contxt, eStop)
	case *Subscribe:
		if ctx.Sender() == nil {
			break
		}
		subscribe(ctx, a.evts)
	case error:
		fmt.Printf("error device actor: %s\n", msg)
	}
}
