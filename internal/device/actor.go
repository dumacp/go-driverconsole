package device

import (
	"context"
	"fmt"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/eventstream"
	"github.com/dumacp/go-logs/pkg/logs"
	"github.com/looplab/fsm"
)

type Actor struct {
	// TODO: ctx???
	ctx           actor.Context
	fmachinae     *fsm.FSM
	evts          *eventstream.EventStream
	dev           Device
	contxt        context.Context
	lastError     time.Time
	cancel        context.CancelFunc
	actvie        bool
	retryInterval time.Duration // Intervalo progresivo para reintentos
	lastRetry     time.Time     // Último intento de reconexión
}

func NewActor(dev Device) actor.Actor {

	a := &Actor{}
	a.contxt = context.TODO()
	a.dev = dev
	a.evts = &eventstream.EventStream{}
	a.Fsm()
	return a
}

func subscribe(ctx actor.Context, evs *eventstream.EventStream) {
	rootctx := ctx.ActorSystem().Root
	pid := ctx.Sender()
	self := ctx.Self()

	fn := func(evt interface{}) {
		rootctx.RequestWithCustomSender(pid, evt, self)
	}
	evs.SubscribeWithPredicate(fn,
		func(evt interface{}) bool {
			switch evt.(type) {
			case *MsgDevice:
				return true
			}
			return false
		})

}

func (a *Actor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *tickMsg:
	default:
		fmt.Printf("message device-actor: %q --> %q, %T\n", func() string {
			if ctx.Sender() == nil {
				return ""
			} else {
				return ctx.Sender().GetId()
			}
		}(), ctx.Self().GetId(), ctx.Message())
	}
	a.ctx = ctx

	switch msg := ctx.Message().(type) {
	case *actor.Started:

		a.retryInterval = 3 * time.Second // Inicializar con 3 segundos
		contxt, cancel := context.WithCancel(context.Background())
		a.cancel = cancel
		go tick(contxt, ctx, 30*time.Second)

		ctx.Send(ctx.Self(), &StartDevice{})
	case *actor.Stopping:
		if a.cancel != nil {
			a.cancel()
		}
		a.fmachinae.Event(a.contxt, eError)
	case *StartDevice:
		if err := a.fmachinae.Event(a.contxt, eStarted); err != nil {
			a.actvie = false
			if time.Since(a.lastError) > 3*time.Minute {
				a.lastError = time.Now()
				logs.LogError.Printf("open device errorn: %s", err)
			}
			// No enviar inmediatamente, esperar al próximo tick
			break
		}
		a.actvie = true
		a.retryInterval = 3 * time.Second // Reset interval on success
		a.lastError = time.Time{}
		fmt.Printf("open device successfully\n")
	case *MsgDevice:
		a.fmachinae.Event(a.contxt, eOpenned)
		a.evts.Publish(msg)
		if ctx.Parent() != nil {
			ctx.Request(ctx.Parent(), msg)
		}

	case *StopDevice:
		a.fmachinae.Event(a.contxt, eClosed)
		a.fmachinae.Event(a.contxt, eStop)
	case *Subscribe:
		if ctx.Sender() == nil {
			break
		}
		subscribe(ctx, a.evts)
	case error:
		fmt.Printf("error device actor: %s\n", msg)
	case *tickMsg:
		if !a.actvie {
			// Solo intentar si ha pasado el intervalo requerido
			if time.Since(a.lastRetry) >= a.retryInterval {
				a.lastRetry = time.Now()
				ctx.Send(ctx.Self(), &StartDevice{})

				// Incrementar progresivamente: 3s -> 6s -> 10s -> 30s -> 60s -> 120s -> 180s (max)
				switch a.retryInterval {
				case 3 * time.Second:
					a.retryInterval = 6 * time.Second
				case 6 * time.Second:
					a.retryInterval = 10 * time.Second
				case 10 * time.Second:
					a.retryInterval = 30 * time.Second
				case 30 * time.Second:
					a.retryInterval = 60 * time.Second
				case 60 * time.Second:
					a.retryInterval = 120 * time.Second
				case 120 * time.Second:
					a.retryInterval = 180 * time.Second
					// Si ya está en 180s, se mantiene ahí
				}
			}
		}
	}
}

type tickMsg struct{}

func tick(contxt context.Context, ctx actor.Context, timeout time.Duration) {
	// Usar un ticker con intervalo mínimo (1 segundo) y controlar desde el actor
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	tick1 := time.NewTicker(timeout)
	defer tick1.Stop()

	for {
		select {
		case <-ticker.C:
			ctx.Send(ctx.Self(), &tickMsg{})
		case <-tick1.C:
			ctx.Send(ctx.Self(), &tickMsg{})
		case <-contxt.Done():
			return
		}
	}
}
