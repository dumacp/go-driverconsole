package parameters

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/dumacp/go-driverconsole/internal/pubsub"
	"github.com/dumacp/go-gwiot/pkg/gwiot"
	"github.com/dumacp/go-params/pkg/params"
)

func NewActor(id string, timeout time.Duration) actor.Actor {
	a := params.Actor(id, false, params.NatsActor(id, gwiot.NewDiscoveryActor("dconsolediscovery/params",
		pubsub.Subscribe,
		pubsub.Publish)))
	return a
}
