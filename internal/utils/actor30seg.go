package utils

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

type MsgTick1 struct {
}
type MsgTickVerifyLive struct {
}

type TrickActor struct {
	timeout time.Duration
	quit    chan int
}

func (a *TrickActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:
		if a.quit != nil {
			select {
			case _, ok := <-a.quit:
				if ok {
					close(a.quit)
				}
			default:
				close(a.quit)
			}
			time.Sleep(100 * time.Millisecond)
		}
		a.quit = make(chan int)
		tick1 := time.NewTicker(a.timeout)
		tick2 := time.NewTicker(120 * time.Second)
		go func(ctx actor.Context) {
			defer tick1.Stop()
			defer tick2.Stop()
			for {
				select {
				case <-tick1.C:
					ctx.Send(ctx.Parent(), &MsgTick1{})
				case <-tick2.C:
					ctx.Send(ctx.Parent(), &MsgTickVerifyLive{})
				case <-a.quit:
					return
				}
			}
		}(ctx)
	case *actor.Terminated:
		close(a.quit)
	}
}

func NewTrickActor(timeout time.Duration) *TrickActor {

	return &TrickActor{timeout: timeout}
}
