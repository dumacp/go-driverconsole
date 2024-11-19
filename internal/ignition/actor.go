package ignition

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/eventstream"
	"github.com/dumacp/go-driverconsole/internal/pubsub"
	ignexternal "github.com/dumacp/go-ignition/pkg/ignition"
	ignmessages "github.com/dumacp/go-ignition/pkg/messages"
	"github.com/dumacp/go-logs/pkg/logs"
)

var (
	timeout = 30 * time.Second
)

type actorIgnition struct {
	timeoutDevices    time.Duration
	lastIgnitionEvent *ignmessages.PowerEvent
	lastEvent         *IgnitionEvent
	evs               *eventstream.EventStream
	pidIgn            *actor.PID
	cancel            func()
}

func NewIgnition() actor.Actor {
	a := &actorIgnition{}
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
		case *IgnitionEvent:
			return true
		}
		return false
	})
}

func (a *actorIgnition) Receive(ctx actor.Context) {
	logs.LogBuild.Printf("Message arrived in ignitionActor: %s, %T, %s",
		ctx.Message(), ctx.Message(), ctx.Sender())
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		ctx.Send(ctx.Self(), &DiscoverService{})
		ctx.Send(ctx.Self(), &ignmessages.PowerStateRequest{})

		conx, cancel := context.WithCancel(context.TODO())
		a.cancel = cancel
		go tick(conx, ctx, timeout)
		if ctx.Parent() != nil {
			ctx.Request(ctx.Parent(), &MsgStarted{})
		}
	case *actor.Stopping:
		if a.cancel != nil {
			a.cancel()
		}
	case *actor.Terminated:
		logs.LogWarn.Printf("\"%s\" - Terminated actor, reason -> %v", ctx.Self(), msg)
		a.pidIgn = nil
	case *ignmessages.PowerEvent:
		a.lastIgnitionEvent = msg
		var state State
		if a.lastIgnitionEvent.GetValue() == ignmessages.StateType_DOWN {
			state = DOWN
		} else {
			state = UP
		}
		timestamp := time.UnixMilli(a.lastIgnitionEvent.GetTimestamp())
		if timestamp.IsZero() || timestamp.UnixMilli() == 0 {
			timestamp = time.Now()
		}

		e := &IgnitionEvent{
			Timestamp:  timestamp,
			StateEvent: state,
		}
		fmt.Printf("ignition event: %v\n", e)
		a.lastEvent = e
		ctx.Send(ctx.Self(), &MsgPublish{})
	case *MsgGetIgnitionEvent:
		// logs.LogInfo.Printf("Ignition Event: %s", msg)
		if ctx.Sender() != nil && a.lastEvent != nil {
			e := &IgnitionEvent{
				Timestamp:  a.lastEvent.Timestamp,
				StateEvent: a.lastEvent.StateEvent,
			}
			ctx.Respond(e)
		}
	case *MsgSubscribe:
		if ctx.Sender() == nil {
			break
		}
		if a.evs == nil {
			a.evs = eventstream.NewEventStream()
		}
		subscribe(ctx, a.evs)
		if a.lastEvent != nil && ctx.Sender() != nil {
			e := &IgnitionEvent{
				Timestamp:  a.lastEvent.Timestamp,
				StateEvent: a.lastEvent.StateEvent,
			}
			ctx.Request(ctx.Sender(), e)
		}
	case *MsgPublish:
		if a.lastEvent == nil {
			break
		}
		e := &IgnitionEvent{
			Timestamp:  a.lastEvent.Timestamp,
			StateEvent: a.lastEvent.StateEvent,
		}
		if a.evs != nil {
			a.evs.Publish(e)
		}
		if ctx.Parent() != nil {
			ctx.Send(ctx.Parent(), e)
		}
	case *MsgTick:
		if a.pidIgn == nil {
			ctx.Send(ctx.Self(), &DiscoverService{})
			break
		}
		if a.lastIgnitionEvent == nil {
			ctx.Request(a.pidIgn, &ignmessages.PowerStateRequest{})
		}
	case *DiscoverService:
		discv := &ignmessages.DiscoverIgnition{
			Id:   ctx.Self().GetId(),
			Addr: ctx.Self().GetAddress(),
		}
		data, err := json.Marshal(discv)
		if err != nil {
			time.Sleep(10 * time.Second)
			log.Panicln(err)
		}
		pubsub.Publish(ignexternal.DISCV_TOPIC, data)
	case *ignmessages.DiscoverResponseIgnition:
		if ctx.Sender() == nil {
			break
		}
		fmt.Printf("%T message: %s\n", msg, msg)
		a.pidIgn = ctx.Sender()
		ctx.Watch(a.pidIgn)
		a.timeoutDevices = time.Duration(msg.GetTimeout()) * time.Millisecond

		ctx.Request(a.pidIgn, &ignmessages.IgnitionPowerSubscription{})
	}
}

func tick(ctx context.Context, ctxactor actor.Context, timeout time.Duration) {
	rootctx := ctxactor.ActorSystem().Root
	self := ctxactor.Self()
	// t0_0 := time.After(300 * time.Millisecond)
	// t0_1 := time.After(600 * time.Millisecond)
	t1 := time.NewTicker(timeout)
	defer t1.Stop()
	for {
		select {
		// case <-t0_0:
		// 	rootctx.Send(self, &MsgTick{})
		// case <-t0_1:
		// 	rootctx.Send(self, &MsgTick{})
		case <-t1.C:
			rootctx.Send(self, &MsgTick{})
		case <-ctx.Done():
			return
		}
	}
}
