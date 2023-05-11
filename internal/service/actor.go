package service

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/dumacp/go-schservices/pkg/services"
)

func NewActor(id string) actor.Actor {
	return services.Actor(id, services.NatsActor(id, &DiscoveryActor{}))
}
