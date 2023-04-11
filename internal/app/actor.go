package app

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/eventstream"

	"github.com/dumacp/go-driverconsole/internal/buttons"
	"github.com/dumacp/go-driverconsole/internal/counterpass"
	"github.com/dumacp/go-driverconsole/internal/database"
	"github.com/dumacp/go-fareCollection/pkg/messages"
	"github.com/dumacp/go-logs/pkg/logs"
)

const (
	dbpath               = "/SD/boltdbs/driverdb"
	databaseName         = "driverdb"
	collectionAppData    = "appfaredata"
	collectionDriverData = "terminaldata"
)

type actorApp struct {
	countUsosParcial int
	countUsosDriver  int
	countUsosAppFare int
	timeLapse        int
	evts             *eventstream.EventStream
	subs             map[string]*eventstream.Subscription
	routes           map[int32]string
	driver           int
	route            int
	routeString      string
	db               *actor.PID
}

func NewActor() actor.Actor {
	a := &actorApp{}
	a.evts = eventstream.NewEventStream()
	a.routes = make(map[int32]string)
	return a
}

func subscribe(ctx actor.Context, evs *eventstream.EventStream) *eventstream.Subscription {
	rootctx := ctx.ActorSystem().Root
	pid := ctx.Sender()
	self := ctx.Self()

	fn := func(evt interface{}) {
		fmt.Printf("%s -> %T -> %s\n", self.GetId(), evt, pid.GetId())
		rootctx.RequestWithCustomSender(pid, evt, self)
	}
	return evs.Subscribe(fn)
}

func subscribeExternal(ctx actor.Context, evs *eventstream.EventStream) *eventstream.Subscription {
	rootctx := ctx.ActorSystem().Root
	pid := ctx.Sender()
	self := ctx.Self()

	fn := func(evt interface{}) {
		rootctx.RequestWithCustomSender(pid, evt, self)
	}
	return evs.SubscribeWithPredicate(fn, func(evt interface{}) bool {
		switch evt.(type) {
		case *messages.MsgSetRoute, *messages.MsgDriverPaso, *messages.MsgGetParams:
			return true
		}
		return false
	})
}

func (a *actorApp) Receive(ctx actor.Context) {

	fmt.Printf("message: %q --> %q, %T (%s)\n", func() string {
		if ctx.Sender() == nil {
			return ""
		} else {
			return ctx.Sender().GetId()
		}
	}(), ctx.Self().GetId(), ctx.Message(), ctx.Message())

	switch msg := ctx.Message().(type) {
	case *actor.Started:
		db, err := database.Open(ctx.ActorSystem().Root, dbpath)
		if err != nil {
			logs.LogWarn.Printf("open database  err: %s\n", err)
		}
		if db != nil {
			a.db = db.PID()
			ctx.Request(a.db, &database.MsgQueryData{
				Database:   databaseName,
				Collection: collectionAppData,
				PrefixID:   "",
				Reverse:    false,
			})
			ctx.Request(a.db, &database.MsgQueryData{
				Database:   databaseName,
				Collection: collectionDriverData,
				PrefixID:   "",
				Reverse:    false,
			})
		}
	case *actor.Stopping:
		if a.db != nil {
			ctx.Send(a.db, &database.MsgCloseDB{})
		}
	case *database.MsgQueryResponse:
		if err := func() error {
			if ctx.Sender() != nil {
				ctx.Send(ctx.Sender(), &database.MsgQueryNext{})
			}
			if len(msg.Data) <= 0 {
				return nil
			}
			data := make([]byte, len(msg.Data))
			copy(data, msg.Data)
			switch msg.Collection {
			case collectionAppData:
				res := &ValidationData{}
				if err := json.Unmarshal(data, res); err != nil {
					return err
				}
				tn := time.Now()
				fmt.Printf("%s - (%d) %s\n", time.Unix(res.Time, 0), tn.Hour(), tn.Add(-time.Duration(tn.Hour())*time.Hour).Truncate(1*time.Hour))
				if time.Unix(res.Time, 0).Before(tn.Add(-time.Duration(tn.Hour()) * time.Hour).Truncate(1 * time.Hour)) {
					break
				}
				a.countUsosAppFare = res.Counter
				a.evts.Publish(&MsgCounters{
					Parcial:  a.countUsosParcial,
					Efectivo: a.countUsosDriver,
					App:      a.countUsosAppFare,
				})
			case collectionDriverData:
				res := &ValidationData{}
				if err := json.Unmarshal(data, res); err != nil {
					return err
				}
				tn := time.Now()
				if time.Unix(res.Time, 0).Before(tn.Add(-time.Duration(tn.Hour()) * time.Hour).Truncate(1 * time.Hour)) {
					break
				}
				a.countUsosDriver = res.Counter
				a.evts.Publish(&MsgCounters{
					Parcial:  a.countUsosParcial,
					Efectivo: a.countUsosDriver,
					App:      a.countUsosAppFare,
				})
			}
			return nil
		}(); err != nil {
			logs.LogWarn.Printf("error updating data from database: %s", err)
		}
	case *MsgAppPaso:
		a.countUsosParcial++
		a.countUsosAppFare += msg.Value
		a.evts.Publish(&MsgCounters{
			Parcial:  a.countUsosParcial,
			Efectivo: a.countUsosDriver,
			App:      a.countUsosAppFare,
		})
		data, err := json.Marshal(&ValidationData{
			Counter: a.countUsosAppFare,
			Time:    time.Now().Unix(),
		})
		if err != nil {
			logs.LogWarn.Printf("error persisting data counters: %s", err)
			break
		}
		ctx.Request(a.db, &database.MsgUpdateData{
			Database:   databaseName,
			Collection: collectionAppData,
			ID:         "validations",
			Data:       data,
		})
	case *messages.MsgAppPaso:
		a.evts.Publish(&MsgConfirmationTextMainScreen{
			Text: []byte(`Validación
Exitosa`),
			Timeout: 2 * time.Second,
		})
		a.countUsosParcial++
		a.countUsosAppFare += int(msg.GetValue())
		a.evts.Publish(&MsgCounters{
			Parcial:  a.countUsosParcial,
			Efectivo: a.countUsosDriver,
			App:      a.countUsosAppFare,
		})
		data, err := json.Marshal(&ValidationData{
			Counter: a.countUsosAppFare,
			Time:    time.Now().Unix(),
		})
		if err != nil {
			logs.LogWarn.Printf("error persisting data counters: %s", err)
			break
		}
		ctx.Request(a.db, &database.MsgUpdateData{
			Database:   databaseName,
			Collection: collectionAppData,
			ID:         "validations",
			Data:       data,
		})
	case *MsgDoors:
		if len(msg.Value) <= 0 {
			break
		}
		a.evts.Publish(msg)
	case *MsgAppPercentRecorrido:
		a.evts.Publish(msg)
	case *buttons.MsgEnterPaso:
		a.evts.Publish(&messages.MsgDriverPaso{
			Value: 1,
		})
	case *messages.MsgWritePayment:
		if ctx.Sender() != nil {
			ctx.Send(ctx.Sender(), &messages.MsgWritePaymentResponse{
				Uid:    msg.Uid,
				Type:   msg.Type,
				Raw:    make(map[string]string),
				Samuid: "",
				Seq:    msg.GetSeq(),
			})
		}
		a.evts.Publish(&buttons.MsgEnableEnterPaso{})
		a.countUsosDriver += 1
		a.countUsosParcial += 1
		a.evts.Publish(&MsgCounters{
			Parcial:  a.countUsosParcial,
			Efectivo: a.countUsosDriver,
			App:      a.countUsosAppFare,
		})
		a.evts.Publish(&MsgConfirmationTextMainScreen{
			Text:    []byte("paso en efectivo"),
			Timeout: 1 * time.Second,
		})
		data, err := json.Marshal(&ValidationData{
			Counter: a.countUsosDriver,
			Time:    time.Now().Unix(),
		})
		if err != nil {
			logs.LogWarn.Printf("error persisting data counters: %s", err)
			break
		}
		ctx.Request(a.db, &database.MsgUpdateData{
			Database:   databaseName,
			Collection: collectionDriverData,
			ID:         "validations",
			Data:       data,
		})
	case *messages.MsgResponseDriverPaso:
		a.evts.Publish(&buttons.MsgEnableEnterPaso{})
		a.countUsosDriver += int(msg.GetValue())
		a.countUsosParcial += int(msg.GetValue())
		a.evts.Publish(&MsgCounters{
			Parcial:  a.countUsosParcial,
			Efectivo: a.countUsosDriver,
			App:      a.countUsosAppFare,
		})
		a.evts.Publish(&MsgConfirmationTextMainScreen{
			Text:    []byte("paso en efectivo"),
			Timeout: 1 * time.Second,
		})
		data, err := json.Marshal(&ValidationData{
			Counter: a.countUsosDriver,
			Time:    time.Now().Unix(),
		})
		if err != nil {
			logs.LogWarn.Printf("error persisting data counters: %s", err)
			break
		}
		ctx.Request(a.db, &database.MsgUpdateData{
			Database:   databaseName,
			Collection: collectionDriverData,
			ID:         "validations",
			Data:       data,
		})
	case *messages.MsgAppError:
		a.evts.Publish(&MsgWarningTextInMainScreen{
			Text:    []byte(msg.Error),
			Timeout: 3 * time.Second,
		})
	case *messages.MsgRoute:
		fmt.Printf("message -> \"%v\"\n", msg)
		v := msg.GetItineraryName()
		if a.routeString == v {
			break
		}
		a.routeString = v
		a.evts.Publish(&MsgRoute{
			v,
		})
	case *buttons.MsgEnterRuta:
		// fmt.Printf("message -> \"%v\"\n", msg)
		if msg.Route >= 0 {
			a.route = msg.Route
		}
		if v, ok := a.routes[int32(a.route)]; ok {
			a.evts.Publish(&MsgSetRoute{
				v,
			})
			a.evts.Publish(&messages.MsgSetRoute{
				Code: int32(a.route),
			})
			a.route = msg.Route
			// a.changeDriver = true
		} else {
			logs.LogWarn.Printf("route \"%d\" not found", a.route)
			a.evts.Publish(&MsgWarningText{
				Text: []byte(`RUTA NO ENCONTRADA`),
			})
		}
	case *buttons.MsgEnterDriver:
		if msg.Driver >= 0 {
			a.driver = msg.Driver
		}
		if a.driver < 1_999_999 {
			logs.LogWarn.Printf("driver \"%d\" error", a.driver)
			a.evts.Publish(&MsgWarningText{
				Text: []byte(`ID inválido,

ingrese un número válido`),
			})
			break
		}
		a.evts.Publish(&MsgDriver{
			Driver: fmt.Sprintf("%d", a.driver),
		})
	case *buttons.MsgChangeRuta:
		if a.route > 0 && ctx.Sender() != nil {
			ctx.Send(ctx.Sender(), &buttons.MsgMemory{
				Key:   buttons.TextNumRoute,
				Value: fmt.Sprintf("%d", a.route),
			})
		}
		fmt.Printf("message -> \"%v\"\n", msg)
		a.evts.Publish(&MsgChangeRoute{ID: a.route})
	case *buttons.MsgChangeDriver:
		if a.route > 0 && ctx.Sender() != nil {
			ctx.Send(ctx.Sender(), &buttons.MsgMemory{
				Key:   buttons.TextNumRoute,
				Value: fmt.Sprintf("%d", a.driver),
			})
		}
		fmt.Printf("message -> \"%v\"\n", msg)
		a.evts.Publish(&MsgChangeDriver{ID: a.driver})
	case *buttons.MsgResetCounter:
		a.countUsosParcial = 0
	case *MsgMainScreen:
		if a.evts != nil {
			a.evts.Publish(msg)
		}
	case *MsgScreen:
		if a.evts != nil {
			a.evts.Publish(msg)
		}
	case *buttons.MsgMainScreen:
		if a.evts != nil {
			a.evts.Publish(&MsgMainScreen{})
		}
	case *buttons.MsgConfirmation, *buttons.MsgWarning, *buttons.MsgReturnFromAlarms, *buttons.MsgReturnFromVehicle:
		if a.evts != nil {
			a.evts.Publish(&MsgMainScreen{})
		}
	case *buttons.MsgStopRecorrido:
		if a.evts != nil {
			a.evts.Publish(&MsgStopRecorrido{})
		}
	case *buttons.MsgInitRecorrido:
		a.countUsosParcial = 0
		a.timeLapse = 0
		if a.evts != nil {
			a.evts.Publish(&MsgInitRecorrido{})
		}
	case *buttons.MsgInputText:
		log.Println(msg.Text)
	case *MsgSubscribe:
		if ctx.Sender() == nil {
			break
		}
		if a.evts == nil {
			a.evts = eventstream.NewEventStream()
		}
		if a.subs == nil {
			a.subs = make(map[string]*eventstream.Subscription)
		}
		if s, ok := a.subs[ctx.Sender().GetId()]; ok {
			a.evts.Unsubscribe(s)
		}
		a.subs[ctx.Sender().GetId()] = subscribe(ctx, a.evts)
	case *messages.MsgSubscribeConsole:

		if ctx.Sender() == nil {
			break
		}
		if a.evts == nil {
			a.evts = eventstream.NewEventStream()
		}
		fmt.Printf("sender = %s\n", ctx.Sender())
		if a.subs == nil {
			a.subs = make(map[string]*eventstream.Subscription)
		}
		if s, ok := a.subs[ctx.Sender().GetId()]; ok {
			a.evts.Unsubscribe(s)
		}
		a.subs[ctx.Sender().GetId()] = subscribeExternal(ctx, a.evts)
	case *MsgSetRoutes:
		a.routes = msg.Routes
	case *messages.MsgAppMapRoute:
		a.routes = msg.GetRoutes()
		if a.route <= 0 && ctx.Sender() != nil {
			ctx.Request(ctx.Sender(), &messages.MsgGetParams{})
		}
	case *MsgConfirmationText:
		if a.evts != nil {
			a.evts.Publish(msg)
		}
	case *MsgConfirmationTextMainScreen:
		if a.evts != nil {
			a.evts.Publish(msg)
		}
	case *MsgWarningText:
		if a.evts != nil {
			a.evts.Publish(msg)
		}
	case *counterpass.CounterMap:
		if a.evts != nil {
			a.evts.Publish(msg)
		}
	case *buttons.MsgDeviceError:
	case *messages.MsgGroundErr:
		if a.evts != nil {
			a.evts.Publish(&MsgNetDown{})
		}
	case *messages.MsgGroundOk:
		if a.evts != nil {
			a.evts.Publish(&MsgNetUP{})
		}
	case *messages.MsgGpsErr:
		if a.evts != nil {
			a.evts.Publish(&MsgGpsDown{})
		}
	case *messages.MsgGpsOk:
		if a.evts != nil {
			a.evts.Publish(&MsgGpsUP{})
		}
	case *messages.MsgAddAlarm:
		if a.evts != nil {
			a.evts.Publish(&MsgAddAlarm{Data: msg.GetAlarm()})
		}
	case *buttons.MsgShowAlarms:
		if a.evts != nil {
			a.evts.Publish(&MsgShowAlarms{})
		}
	case *buttons.MsgBrightnessMinus:
		if a.evts != nil {
			a.evts.Publish(&MsgBrightnessMinus{})
		}
	case *buttons.MsgBrightnessPlus:
		if a.evts != nil {
			a.evts.Publish(&MsgBrightnessPlus{})
		}
	}
}
