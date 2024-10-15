package app

import (
	"testing"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

func TestApp_Receive(t *testing.T) {
	type args struct {
		ctx *actor.RootContext
	}
	tests := []struct {
		name string
		args args
	}{

		{
			name: "test_resetCounters",
			args: args{
				ctx: actor.NewActorSystem().Root,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			a := &App{
				countInput:  10,
				countOutput: 30,
			}
			// TODO: comment out for test
			// tn := time.Now()
			// tRefg = time.Date(tn.Year(), tn.Month(), tn.Day(), tn.Hour(), tn.Minute(), tn.Second()+5, 0, tn.Location())
			pid, err := tt.args.ctx.SpawnNamed(actor.PropsFromFunc(a.Receive), tt.name)
			if err != nil {
				t.Error(err)
			}

			time.Sleep(10 * time.Second)
			tt.args.ctx.PoisonFuture(pid).Wait()
			time.Sleep(1 * time.Second)

		})
	}
}
