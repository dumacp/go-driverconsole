package parameters

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/dumacp/go-params/pkg/params"
)

func NewActor(id string) actor.Actor {
	return params.Actor(id, false, params.NatsActor(id, &DiscoveryActor{}))
}
