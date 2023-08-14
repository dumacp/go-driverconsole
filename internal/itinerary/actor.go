package itinerary

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/dumacp/go-itinerary/pkg/route"
	"github.com/dumacp/go-params/pkg/params"
)

func NewActor(id string) actor.Actor {
	return params.Actor(id, false, route.NatsActor(id, &DiscoveryActor{}))
}
