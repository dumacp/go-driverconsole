package counterpass

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/eventstream"
	"github.com/dumacp/go-driverconsole/internal/pubsub"
	"github.com/dumacp/go-logs/pkg/logs"
	psub "github.com/dumacp/pubsub"
)

// Actor actor to listen events
type Actor struct {
	ctx            actor.Context
	evs            *eventstream.EventStream
	lastCounterMap *CounterMap
}

func NewActor() actor.Actor {
	return &Actor{}
}

func parseCounters(msg []byte) interface{} {

	event := new(CounterMap)
	if err := json.Unmarshal(msg, event); err != nil {
		fmt.Printf("error parseEvents = %s\n", err)
		return err
	}

	return event
}

func parseEvents(msg []byte) interface{} {

	message := &psub.Message{
		// Timestamp: float64(time.Now().UnixNano()) / 1000000000,
		// Type:      "COUNTERSDOOR",
	}

	// fmt.Printf("********* event msg: %s\n", msg)

	val := struct {
		// Coord    string  `json:"coord"`
		ID       int32   `json:"id"`
		State    uint    `json:"state"`
		Counters []int64 `json:"counters"`
		Type     string  `json:"type,omitempty"`
	}{}
	message.Value = &val

	if err := json.Unmarshal(msg, message); err != nil {
		fmt.Printf("error parseEvents = %s\n", err)
		return err
	}

	// fmt.Printf("********* parse event: %v, value: %v\n", message, message.Value)

	event := new(CounterEvent)

	if val.Counters != nil && len(val.Counters) > 1 {
		event.Inputs = int(val.Counters[0])
		event.Outputs = int(val.Counters[1])
	}

	return event
}

func parseExtraEvents(msg []byte) interface{} {

	message := &psub.Message{
		// Timestamp: float64(time.Now().UnixNano()) / 1000000000,
		// Type:      "TAMPERING",
	}

	val := struct {
		// Coord    string  `json:"coord"`
		ID       int32   `json:"id"`
		State    uint    `json:"state"`
		Counters []int64 `json:"counters"`
		Type     string  `json:"type,omitempty"`
	}{}
	message.Value = val

	if err := json.Unmarshal(msg, message); err != nil {
		fmt.Printf("error parseExtraEvents = %s\n", err)
		return err
	}

	if !strings.Contains(message.Type, "TAMPERING") {
		return fmt.Errorf("extraEvent not configured, type: %s", message.Type)
	}

	event := new(CounterExtraEvent)
	event.Text = []byte(fmt.Sprintf("Evento: %s", message.Type))

	return event
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
		case *CounterMap:
			return true
		}
		return false
	})
}

// Receive func Receive in actor
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
		logs.LogInfo.Printf("started \"%s\", %v", ctx.Self().GetId(), ctx.Self())
		if err := pubsub.Subscribe("COUNTERSMAPDOOR", ctx.Self(), parseCounters); err != nil {
			time.Sleep(3 * time.Second)
			logs.LogError.Panic(err)
		}
		if err := pubsub.Subscribe("EVENTS/backcounter", ctx.Self(), parseEvents); err != nil {
			time.Sleep(3 * time.Second)
			logs.LogError.Panic(err)
		}
		if err := pubsub.Subscribe("EVENTS/counterevents", ctx.Self(), parseExtraEvents); err != nil {
			time.Sleep(3 * time.Second)
			logs.LogError.Panic(err)
		}
	case *actor.Stopping:
		logs.LogWarn.Printf("\"%s\" - Stopped actor, reason -> %v", ctx.Self(), msg)
	case *actor.Restarting:
		logs.LogWarn.Printf("\"%s\" - Restarting actor, reason -> %v", ctx.Self(), msg)
	case *actor.Terminated:
		logs.LogWarn.Printf("\"%s\" - Terminated actor, reason -> %v", ctx.Self(), msg)
	case *CounterMap:
		a.lastCounterMap = msg
		if a.evs != nil {
			a.evs.Publish(msg)
		}
	case *CounterEvent:
		if a.evs != nil {
			a.evs.Publish(msg)
		}
		if ctx.Parent() != nil {
			ctx.Send(ctx.Parent(), msg)
		}
	case *CounterExtraEvent:
		if a.evs != nil {
			a.evs.Publish(msg)
		}
		if ctx.Parent() != nil {
			ctx.Send(ctx.Parent(), msg)
		}
	case *MsgSubscribe:
		if ctx.Sender() == nil {
			break
		}
		if a.evs == nil {
			a.evs = eventstream.NewEventStream()
		}
		subscribe(ctx, a.evs)
		if a.lastCounterMap != nil {
			counters := new(CounterMap)
			*counters = *a.lastCounterMap
			ctx.Respond(counters)
		}
	case *MsgRequestStatus:
		if ctx.Sender() != nil {
			break
		}
		ctx.Respond(&MsgStatus{State: true})
	}
}
