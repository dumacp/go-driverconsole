package app

import (
	"fmt"
	"log"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/eventstream"

	"github.com/dumacp/go-driverconsole/internal/buttons"
	"github.com/dumacp/go-driverconsole/pkg/messages"
	"github.com/dumacp/go-logs/pkg/logs"
)

type actorApp struct {
	pidAppFare       *actor.PID
	pidDisplay       *actor.PID
	countUsosParcial int
	countUsosDriver  int
	countUsosAppFare int
	timeLapse        int
	flagPaso         chan int
	evts             *eventstream.EventStream
	routes           map[int]string
	changeRoute      bool
	changeDriver     bool
	driver           int
	route            int
}

func NewActor() actor.Actor {
	a := &actorApp{}
	a.evts = eventstream.NewEventStream()
	a.routes = make(map[int]string)
	return a
}

func subscribe(ctx actor.Context, evs *eventstream.EventStream) {
	rootctx := ctx.ActorSystem().Root
	pid := ctx.Sender()
	self := ctx.Self()

	fn := func(evt interface{}) {
		rootctx.RequestWithCustomSender(pid, evt, self)
	}
	evs.Subscribe(fn)
	// sub.WithPredicate(func(evt interface{}) bool {
	// 	switch evt.(type) {
	// 	case *Msg, *MsgInitRecorrido, *MsgStopRecorrido:
	// 		return true
	// 	}
	// 	return false
	// })
}

func subscribeExternal(ctx actor.Context, evs *eventstream.EventStream) {
	rootctx := ctx.ActorSystem().Root
	pid := ctx.Sender()
	self := ctx.Self()

	fn := func(evt interface{}) {
		rootctx.RequestWithCustomSender(pid, evt, self)
	}
	evs.SubscribeWithPredicate(fn, func(evt interface{}) bool {
		switch evt.(type) {
		case *messages.MsgDriverPaso:
			return true
		}
		return false
	})
}
func (a *actorApp) Receive(ctx actor.Context) {

	fmt.Printf("message -> \"%s\", %T\n", ctx.Self().GetId(), ctx.Message())

	switch msg := ctx.Message().(type) {
	case *actor.Started:
		if a.flagPaso == nil {
			a.flagPaso = make(chan int)
		}
	case *MsgAppPaso:
		// a.countUsosDriver++
		a.countUsosParcial++
		a.countUsosAppFare += msg.Value
		a.evts.Publish(&MsgCounters{
			Parcial:  a.countUsosParcial,
			Efectivo: a.countUsosDriver,
			App:      a.countUsosAppFare,
		})
	case *messages.MsgAppPaso:
		a.countUsosParcial++
		a.countUsosAppFare += int(msg.GetValue())
		a.evts.Publish(&MsgCounters{
			Parcial:  a.countUsosParcial,
			Efectivo: a.countUsosDriver,
			App:      a.countUsosAppFare,
		})
	case *messages.MsgAppDoor:
		v := &MsgDoors{
			Value: [2]int{int(msg.GetId()), int(msg.GetValue())},
		}
		ctx.Send(ctx.Self(), v)
	case *MsgDoors:
		if len(msg.Value) <= 0 {
			break
		}
		a.evts.Publish(msg)
	case *MsgAppPercentRecorrido:
		a.evts.Publish(msg)
	case *buttons.MsgEnterPaso:
		a.countUsosDriver++
		a.countUsosParcial++
		a.evts.Publish(&MsgCounters{
			Parcial:  a.countUsosParcial,
			Efectivo: a.countUsosDriver,
			App:      a.countUsosAppFare,
		})
		a.evts.Publish(&messages.MsgDriverPaso{
			Value: 1,
		})
	case *messages.MsgAppRoute:
		fmt.Printf("message -> \"%v\"\n", msg)
		v := msg.GetName()
		a.evts.Publish(&MsgRoute{
			v,
		})
	case *buttons.MsgEnterRuta:
		fmt.Printf("message -> \"%v\"\n", msg)
		switch {
		case a.changeRoute:
			if v, ok := a.routes[msg.Route]; ok {
				a.evts.Publish(&MsgRoute{
					v,
				})
				a.route = msg.Route
				a.changeDriver = true
			} else {
				logs.LogWarn.Printf("route \"%d\" not found", msg.Route)
				a.evts.Publish(&MsgWarningText{
					Text: []byte(`




    RUTA NO ENCONTRADA`),
				})
			}
			a.changeRoute = false
		case a.changeDriver:
			if msg.Route < 10_999_999 {
				logs.LogWarn.Printf("driver \"%d\" error", msg.Route)
				a.evts.Publish(&MsgWarningText{
					Text: []byte(`




    ID inválido, 
	ingrese un número válido`),
				})
				break
			}
			a.evts.Publish(&MsgDriver{
				Driver: fmt.Sprintf("%d", msg.Route),
			})
			a.driver = msg.Route
			a.changeDriver = false

		default:
		}
	case *buttons.MsgChangeRuta:
		if a.route > 0 && ctx.Sender() != nil {
			ctx.Send(ctx.Sender(), &buttons.MsgMemory{
				Key:   "textNumLabel",
				Value: fmt.Sprintf("%d", a.route),
			})
		}
		a.changeRoute = true
		fmt.Printf("message -> \"%v\"\n", msg)
		a.evts.Publish(&MsgChangeRoute{ID: a.route})
	case *buttons.MsgResetCounter:
		//log.Println("main ResetCounter")
		a.countUsosParcial = 0
		// a.evts.Publish(&ResetRecorrido{})
	case *MsgMainScreen:
		a.evts.Publish(msg)
	case *buttons.MsgMainScreen:
		a.evts.Publish(&MsgMainScreen{})
	case *buttons.MsgConfirmation:
		switch {
		case a.changeDriver:
			if a.driver > 0 && ctx.Sender() != nil {
				ctx.Send(ctx.Sender(), &buttons.MsgMemory{
					Key:   buttons.TextNumLabel,
					Value: fmt.Sprintf("%d", a.driver),
				})
			}
			a.evts.Publish(&MsgChangeDriver{
				ID: a.driver,
			})
		default:
			a.evts.Publish(&MsgConfirmationButton{})
		}
	case *buttons.MsgWarning:
		a.evts.Publish(&MsgWarningButton{})
	case *buttons.MsgStopRecorrido:
		a.evts.Publish(&MsgStopRecorrido{})
	case *buttons.MsgInitRecorrido:
		//log.Println("main ResetRecorrido")
		a.countUsosParcial = 0
		a.timeLapse = 0
		a.evts.Publish(&MsgInitRecorrido{})
	case *buttons.MsgInputText:
		log.Println(msg.Text)
	case *MsgSubscribe:
		if ctx.Sender() == nil {
			break
		}
		subscribe(ctx, a.evts)
	case *messages.MsgSubscribe:
		if ctx.Sender() == nil {
			break
		}
		subscribeExternal(ctx, a.evts)
	case *MsgSetRoutes:
		a.routes = msg.Routes
	case *MsgConfirmationText:
		// a.evts.Publish(&MsgScreen{ID: 3, Switch: true})
		a.evts.Publish(msg)
	case *MsgWarningText:
		// a.evts.Publish(&MsgScreen{ID: 3, Switch: true})
		a.evts.Publish(msg)
	}
}
