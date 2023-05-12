package app

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/eventstream"

	"github.com/dumacp/go-actors/database"
	"github.com/dumacp/go-driverconsole/internal/buttons"
	"github.com/dumacp/go-driverconsole/internal/device"
	"github.com/dumacp/go-driverconsole/internal/ui"
	"github.com/dumacp/go-logs/pkg/logs"
	"github.com/dumacp/go-schservices/api/services"
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
	shcservices      map[string]*services.ScheduleService
	driver           int
	route            int
	routeString      string
	evt2evtApp       map[buttons.KeyCode]EventLabel
	uix              ui.UI
	db               *actor.PID
	contxt           context.Context
	buttonDevice     buttons.ButtonDevice
	buttonCancel     func()
	cancel           func()
}

func NewActor(uix ui.UI, button2buttonApp map[buttons.KeyCode]EventLabel) actor.Actor {
	a := &actorApp{}
	a.uix = uix
	a.evts = eventstream.NewEventStream()
	a.routes = make(map[int32]string)
	a.shcservices = make(map[string]*services.ScheduleService)
	a.evt2evtApp = button2buttonApp
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
		// switch evt.(type) {
		// case *messages.MsgSetRoute, *messages.MsgDriverPaso, *messages.MsgGetParams:
		// 	return true
		// }
		// return false
		return true
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
		db, err := database.Open(ctx, dbpath)
		if err != nil {
			logs.LogWarn.Printf("open database  err: %s\n", err)
		}
		if db != nil {
			a.db = db.PID()
			ctx.Request(a.db, &database.MsgQueryData{
				Buckets:  []string{databaseName, collectionAppData},
				PrefixID: "",
				Reverse:  false,
			})
			ctx.Request(a.db, &database.MsgQueryData{
				Buckets:  []string{databaseName, collectionDriverData},
				PrefixID: "",
				Reverse:  false,
			})
		}
	case *actor.Stopping:
		if a.db != nil {
			ctx.Send(a.db, &database.MsgCloseDB{})
		}
		if a.cancel != nil {
			a.cancel()
		}
	case *database.MsgQueryResponse:
		if err := func() error {
			if ctx.Sender() != nil {
				ctx.Send(ctx.Sender(), &database.MsgQueryNext{})
			}
			if len(msg.Data) <= 0 {
				return nil
			}
			if len(msg.Buckets) < 1 {
				return nil
			}
			data := make([]byte, len(msg.Data))
			copy(data, msg.Data)
			coll := msg.Buckets[len(msg.Buckets)-1]
			switch coll {
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

	case *MsgDoors:
		if len(msg.Value) <= 0 {
			break
		}
		a.evts.Publish(msg)
	case *device.MsgDevice:
		contxt, cancel := context.WithCancel(context.TODO())
		a.contxt = contxt
		a.cancel = cancel

		a.buttonDevice.Init(msg.Device)
		ch, err := a.buttonDevice.ListenButtons(contxt)
		if err != nil {
			logs.LogError.Printf("listenButtons error = %s", err)
			break
		}
		go func(ctx actor.Context) {
			rootctx := ctx.ActorSystem().Root
			self := ctx.Self()

			for v := range ch {
				rootctx.Request(self, v)
			}
		}(ctx)
	case *buttons.InputEvent:
		label, ok := a.evt2evtApp[msg.KeyCode]
		if !ok {
			break
		}
		if a.buttonCancel != nil {
			a.buttonCancel()
		}
		if err := func() error {
			switch label {
			case PROGRAMATION_DRIVER:
				if err := a.uix.ShowProgDriver(); err != nil {
					return fmt.Errorf("event ShowProgDriver error: %s", err)
				}
			case PROGRAMATION_VEH:
				if err := a.uix.ShowProgVeh(); err != nil {
					return fmt.Errorf("event ShowProgVeh error: %s", err)
				}
			case SHOW_NOTIF:
				if err := a.uix.ShowNotifications(); err != nil {
					return fmt.Errorf("event ShowNotifications error: %s", err)
				}
			case STATS:
				if err := a.uix.ShowStats(); err != nil {
					return fmt.Errorf("event ShowStats error: %s", err)
				}
			case ROUTE:
				contxt, cancel := context.WithCancel(a.contxt)
				a.buttonCancel = cancel
				num, err := a.uix.KeyNum(contxt, "ingrese el número de ruta:")
				if err != nil {
					return fmt.Errorf("route keyNum error: %s", err)
				}
				go func() {
					defer cancel()
					self := ctx.Self()
					rootctx := ctx.ActorSystem().Root

					select {
					case <-contxt.Done():
					case v := <-num:
						rootctx.Send(self, &MsgChangeRoute{
							ID: v,
						})
					}
				}()
			case DRIVER:
				contxt, cancel := context.WithCancel(a.contxt)
				a.buttonCancel = cancel
				num, err := a.uix.KeyNum(contxt, "ingrese el ID del conductor:")
				if err != nil {
					return fmt.Errorf("driver keyNum error: %s", err)
				}
				go func() {
					defer cancel()
					self := ctx.Self()
					rootctx := ctx.ActorSystem().Root

					select {
					case <-contxt.Done():
					case v := <-num:
						rootctx.Send(self, &MsgChangeDriver{
							ID: v,
						})
					}
				}()
			}
			return nil
		}(); err != nil {
			logs.LogWarn.Println(err)
		}

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
	// case *messages.MsgSubscribeConsole:

	// 	if ctx.Sender() == nil {
	// 		break
	// 	}
	// 	if a.evts == nil {
	// 		a.evts = eventstream.NewEventStream()
	// 	}
	// 	fmt.Printf("sender = %s\n", ctx.Sender())
	// 	if a.subs == nil {
	// 		a.subs = make(map[string]*eventstream.Subscription)
	// 	}
	// 	if s, ok := a.subs[ctx.Sender().GetId()]; ok {
	// 		a.evts.Unsubscribe(s)
	// 	}
	// 	a.subs[ctx.Sender().GetId()] = subscribeExternal(ctx, a.evts)
	case *MsgSetRoutes:
		a.routes = msg.Routes
	case *MsgConfirmationText:
		if err := a.uix.TextConfirmation(string(msg.Text)); err != nil {
			logs.LogWarn.Printf("textConfirmation error: %s", err)
		}
	case *MsgConfirmationTextMainScreen:
		if err := a.uix.TextConfirmationPopup(3*time.Second, string(msg.Text)); err != nil {
			logs.LogWarn.Printf("textConfirmation error: %s", err)
		}
	case *MsgWarningText:
		if err := a.uix.TextWarning(string(msg.Text)); err != nil {
			logs.LogWarn.Printf("textConfirmation error: %s", err)
		}
	case *services.UpdateServiceMsg:
		svc := msg.GetUpdate()
		if len(svc.GetState()) > 0 {
			if v, ok := a.shcservices[svc.GetId()]; ok {
				fmt.Printf("////// route: %v\n", v)
				UpdateService(v, msg.GetUpdate())
			} else {
				a.shcservices[svc.GetId()] = msg.GetUpdate()
			}
		}
		fmt.Printf("////////////// ROUTE: %+v\n", a.shcservices[svc.GetId()].GetRoute())
		if a.shcservices[svc.GetId()].GetRoute() != nil &&
			len(a.shcservices[svc.GetId()].GetRoute().Name) > 0 {
			a.uix.Route(a.shcservices[svc.GetId()].GetRoute().Name)
		}

		if svc.GetCheckpointTimingState() != nil && len(svc.GetCheckpointTimingState().GetState()) > 0 {
			state := int(services.TimingState_value[svc.GetCheckpointTimingState().GetState()])
			promtp := fmt.Sprintf("%s (%d)", svc.GetCheckpointTimingState().GetName(), svc.GetCheckpointTimingState().GetTimeDiff())
			fmt.Printf("///// state: %d\n", state)
			if err := a.uix.ServiceCurrentState(state, promtp); err != nil {
				logs.LogWarn.Printf("textConfirmation error: %s", err)
			}
		}
	case *services.ServiceMsg:
		svc := msg.GetUpdate()
		if len(svc.GetState()) > 0 {
			a.shcservices[svc.GetId()] = svc
		} else {
			if v, ok := a.shcservices[svc.GetId()]; ok {
				fmt.Printf("////// route: %v\n", v)
				UpdateService(v, svc)
			} else {
				a.shcservices[svc.GetId()] = svc
			}
		}
	case *services.ServiceAllMsg:
		svcs := msg.GetUpdates()
		a.shcservices = make(map[string]*services.ScheduleService)
		for _, svc := range svcs {
			if v, ok := a.shcservices[svc.GetId()]; ok {
				fmt.Printf("////// route: %v\n", v)
				UpdateService(v, svc)
			} else {
				a.shcservices[svc.GetId()] = svc
			}
		}
		arr := make([]string, 0)
		for k, _ := range a.shcservices {
			arr = append(arr, k)
		}
		fmt.Printf("//////////////// services (ori: %d): %v\n", len(msg.GetUpdates()), arr)
	case *services.RemoveServiceMsg:
		svc := msg.GetUpdate()
		delete(a.shcservices, svc.GetId())
	case *MsgScreen:
		if err := a.uix.Screen(msg.ID, msg.Switch); err != nil {
			logs.LogWarn.Printf("textConfirmation error: %s", err)
		}

	}
}
