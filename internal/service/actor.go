package service

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/dumacp/go-driverconsole/internal/pubsub"
	"github.com/dumacp/go-gwiot/pkg/gwiot"
	"github.com/dumacp/go-schservices/pkg/services"
)

func NewActor(id, url string) actor.Actor {
	return services.Actor(id, url, services.NatsActor(id, gwiot.NewDiscoveryActor(services.TOPIC_REPLY,
		pubsub.Subscribe,
		pubsub.Publish)))
}
