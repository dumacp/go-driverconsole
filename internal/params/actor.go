package params

import (
	"github.com/asynkron/protoactor-go/actor"
)

type Actor struct {
	props *actor.Props
}

func (a *Actor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:

	}
}
