package ground

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/dumacp/go-driverconsole/internal/utils"
	"github.com/dumacp/go-logs/pkg/logs"
)

type groundActor struct {
	url      string
	lastPing time.Time
	client   *http.Client
	cancel   func()
	isOk     bool
}

func NewActor(url string) actor.Actor {
	a := &groundActor{}
	a.url = url
	return a

}

var timeout = 300 * time.Second

func (a *groundActor) Receive(ctx actor.Context) {
	logs.LogBuild.Printf("Message arrived in groundActor: %s, %T, %s",
		ctx.Message(), ctx.Message(), ctx.Sender())
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		if ctx.Parent() != nil {
			ctx.Request(ctx.Parent(), &MsgStarted{})
		}
		contxt, cancel := context.WithCancel(context.TODO())
		a.cancel = cancel
		go tick(contxt, ctx, timeout)
	case *actor.Stopping:
		if a.cancel != nil {
			a.cancel()
		}
	case *MsgRequestStatus:
		if ctx.Sender() == nil {
			break
		}
		if a.isOk {
			ctx.Respond(&MsgGroundOk{})
		} else {
			ctx.Respond(&MsgGroundErr{})
		}
	case *MsgVerify:
		if !a.isOk && !msg.Force {
			break
		}
		url := fmt.Sprintf("%s/auth/", a.url)
		if data, err := utils.PingHttp(a.client, url); err != nil {
			a.isOk = false
			log.Printf("ground error (url: %s): %s", url, err)
			if ctx.Sender() != nil {
				ctx.Respond(&MsgGroundErr{})
			}
		} else {
			logs.LogBuild.Printf("ground verify response: %d", data)
			a.isOk = true
			if ctx.Sender() != nil {
				ctx.Respond(&MsgGroundOk{})
			}
		}
	case *MsgTick:
		if !a.isOk || time.Since(a.lastPing) > 3*time.Minute {
			a.lastPing = time.Now()
			url := fmt.Sprintf("%s/auth/", a.url)
			if data, err := utils.PingHttp(a.client, url); err == nil {
				logs.LogBuild.Printf("ground verify response: %d", data)
				a.isOk = true
				if ctx.Parent() != nil {
					ctx.Send(ctx.Parent(), &MsgGroundOk{})
				}
			} else {
				a.isOk = false
				log.Printf("ground error (url: %s): %s", url, err)
				if ctx.Parent() != nil {
					ctx.Send(ctx.Parent(), &MsgGroundErr{})
				}
			}
		}
	case error:
		fmt.Printf("error message: %s (%s)\n", msg, ctx.Self().GetId())
	default:
		fmt.Printf("unhandled message type: %T (%s)\n", msg, ctx.Self().GetId())
	}
}

func tick(contxt context.Context, ctx actor.Context, timeout time.Duration) {
	rootctx := ctx.ActorSystem().Root
	self := ctx.Self()
	t0_0 := time.After(3 * time.Second)
	t0_1 := time.After(10 * time.Second)
	t0_2 := time.After(30 * time.Second)
	timeout_ := timeout
	t1 := time.NewTimer(timeout_)
	defer t1.Stop()
	for {
		select {
		case <-t0_0:
			rootctx.Send(self, &MsgTick{})
		case <-t0_1:
			rootctx.Send(self, &MsgTick{})
		case <-t0_2:
			rootctx.Send(self, &MsgTick{})
		case <-t1.C:
			rootctx.Send(self, &MsgTick{})
			t1.Reset(timeout_)
		case <-contxt.Done():
			return
		}
	}
}
