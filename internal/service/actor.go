package service

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/dumacp/go-driverconsole/internal/pubsub"
	"github.com/dumacp/go-gwiot/pkg/gwiot"
	"github.com/dumacp/go-params/pkg/params"
	"github.com/dumacp/go-schservices/pkg/services"
)

func NewActor(id string) actor.Actor {
	return services.Actor(id, services.NatsActor(id, gwiot.NewDiscoveryActor(params.TOPIC_REPLY,
		pubsub.Subscribe,
		pubsub.Publish)))
}
