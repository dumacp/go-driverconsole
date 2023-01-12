package gps

import (
	"bytes"
	"fmt"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/eventstream"
	"github.com/dumacp/go-driverconsole/internal/pubsub"
	"github.com/dumacp/go-logs/pkg/logs"
	"github.com/dumacp/gpsnmea"
	"github.com/golang/geo/s2"
)

// Actor actor to listen events
type Actor struct {
	ctx       actor.Context
	evs       *eventstream.EventStream
	subs      map[string]*eventstream.Subscription
	lastFrame GPSData
}

func NewActor() actor.Actor {
	a := &Actor{}
	a.subs = make(map[string]*eventstream.Subscription)
	a.evs = eventstream.NewEventStream()
	return a
}

func parseEvents() func([]byte) interface{} {

	mem := new(GPSData)

	return func(msg []byte) interface{} {

		event := new(GPSData)
		switch {
		case bytes.Contains(msg, []byte("GPRMC")):
			rmc := gpsnmea.ParseRMC(string(msg))
			lat := gpsnmea.LatLongToDecimalDegree(rmc.Lat, rmc.LatCord)
			lon := gpsnmea.LatLongToDecimalDegree(rmc.Lat, rmc.LongCord)
			latlon := s2.LatLngFromDegrees(lat, lon)
			event.LatLon = latlon
			event.Speed = float32(rmc.Speed * 1.852)
			if mem.HDop >= 0 && time.Since(mem.Time) < 5*time.Second {
				event.HDop = mem.HDop
			} else {
				event.HDop = float32(-1)
			}
		case bytes.Contains(msg, []byte("GPGGA")):
			gga := gpsnmea.ParseGGA(string(msg))
			lat := gpsnmea.LatLongToDecimalDegree(gga.Lat, gga.LatCord)
			lon := gpsnmea.LatLongToDecimalDegree(gga.Lat, gga.LongCord)
			latlon := s2.LatLngFromDegrees(lat, lon)
			event.LatLon = latlon
			event.HDop = float32(gga.HDop)
			if mem.Speed >= 0 && time.Since(mem.Time) < 5*time.Second {
				event.Speed = mem.Speed
			} else {
				event.Speed = float32(-1)
			}
		default:
			return fmt.Errorf("unknown frame")
		}
		mem = event
		event.Time = time.Now()

		return event
	}
}

func subscribe(ctx actor.Context, evs *eventstream.EventStream) {
	rootctx := ctx.ActorSystem().Root
	pid := ctx.Sender()
	self := ctx.Self()

	fn := func(evt interface{}) {
		rootctx.RequestWithCustomSender(pid, evt, self)
	}
	evs.SubscribeWithPredicate(fn, func(evt interface{}) bool {
		switch evt.(type) {
		case *GPSData:
			return true
		}
		return false
	})
}

// Receive func Receive in actor
func (a *Actor) Receive(ctx actor.Context) {
	fmt.Printf("message: %q --> %q, %T\n", func() string {
		if ctx.Sender() == nil {
			return ""
		} else {
			return ctx.Sender().GetId()
		}
	}(), ctx.Self().GetId(), ctx.Message())

	a.ctx = ctx
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		logs.LogInfo.Printf("started \"%s\", %v", ctx.Self().GetId(), ctx.Self())
		if err := pubsub.Subscribe("GPS", ctx.Self(), parseEvents()); err != nil {
			time.Sleep(3 * time.Second)
			logs.LogError.Panic(err)
		}
	case *actor.Stopping:
		logs.LogWarn.Printf("\"%s\" - Stopped actor, reason -> %v", ctx.Self(), msg)
	case *actor.Restarting:
		logs.LogWarn.Printf("\"%s\" - Restarting actor, reason -> %v", ctx.Self(), msg)
	case *actor.Terminated:
		logs.LogWarn.Printf("\"%s\" - Terminated actor, reason -> %v", ctx.Self(), msg)
	case *GPSData:
		a.lastFrame = *msg

	case *MsgSubscribe:
		if ctx.Sender() == nil {
			break
		}
		if a.evs == nil {
			a.evs = eventstream.NewEventStream()
		}
		if v, ok := a.subs[ctx.Self().GetId()]; ok {
			a.evs.Unsubscribe(v)
		}
		subscribe(ctx, a.evs)
		if time.Since(a.lastFrame.Time) < 30*time.Second {
			ctx.Respond(&MsgGpsData{Data: a.lastFrame})
		}
	}
}
