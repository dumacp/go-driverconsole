package itinerary

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/dumacp/go-driverconsole/internal/pubsub"
	"github.com/dumacp/go-gwiot/pkg/gwiot"
	"github.com/dumacp/go-itinerary/pkg/route"
)

func NewActor(id string) actor.Actor {
	return route.Actor(id, route.NatsActor(id, gwiot.NewDiscoveryActor(route.TOPIC_REPLY,
		pubsub.Subscribe,
		pubsub.Publish)))
}
