package service

import (
	"encoding/json"
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/dumacp/go-driverconsole/internal/pubsub"
	"github.com/dumacp/go-gwiot/pkg/gwiotmsg"
	"github.com/dumacp/go-gwiot/pkg/gwiotmsg/gwiot"
	"github.com/dumacp/go-logs/pkg/logs"
	"github.com/dumacp/go-schservices/pkg/services"
)

type DiscoveryActor struct {
}

func Parse(m []byte) interface{} {
	ev := &gwiotmsg.DiscoveryResponse{}

	if err := json.Unmarshal(m, ev); err != nil {
		return fmt.Errorf("error: %s", err)
	}
	return ev
}

func (a *DiscoveryActor) Receive(ctx actor.Context) {
	logs.LogBuild.Printf("Message arrived in %s: %s, %T, %s",
		ctx.Self().GetId(), ctx.Message(), ctx.Message(), ctx.Sender())
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		if err := pubsub.Subscribe(services.TOPIC_REPLY, ctx.Self(), Parse); err != nil {
			logs.LogError.Printf("subscribe pubsub to %s error: %s", services.TOPIC_REPLY, err)
			break
		}
		logs.LogInfo.Printf("started \"%s\", %v", ctx.Self().GetId(), ctx.Self())
	case *actor.Stopping:
	case *actor.Terminated:
	case *gwiotmsg.Discovery:
		disc := &gwiotmsg.Discovery{
			Reply: services.TOPIC_REPLY,
		}
		if data, err := json.Marshal(disc); err != nil {
			logs.LogWarn.Printf("error discover request: %s", err)
		} else {
			pubsub.Publish(gwiot.SUBS_TOPIC, data)
		}
	case *gwiotmsg.DiscoveryResponse:
		if ctx.Parent() != nil {
			ctx.Send(ctx.Parent(), msg)
		}
	}
}
