package app

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/eventstream"

	"github.com/dumacp/go-actors/database"
	"github.com/dumacp/go-driverconsole/internal/constant"
	"github.com/dumacp/go-driverconsole/internal/counterpass"
	"github.com/dumacp/go-driverconsole/internal/pubsub"
	"github.com/dumacp/go-driverconsole/internal/ui"
	"github.com/dumacp/go-driverconsole/internal/utils"
	msgdriverterminal "github.com/dumacp/go-driverconsole/pkg/messages"
	"github.com/dumacp/go-fareCollection/pkg/messages"
	"github.com/dumacp/go-itinerary/api/routes"
	"github.com/dumacp/go-logs/pkg/logs"
	"github.com/dumacp/go-params/api/params"
	"github.com/dumacp/go-schservices/api/services"
)

const (
	dbpath               = "/SD/boltdbs/driverdb"
	databaseName         = "driverdb"
	collectionAppData    = "appfaredata"
	collectionDriverData = "terminaldata"

	TIMEOUT = 10 * time.Second
)

type App struct {
	updateTime           time.Time
	lastErrVerifyDisplay errorDisplay
	lastErrButtons       errorDisplay
	totalCountInput      int32
	countInput           int32
	countOutput          int32
	cashInput            int32
	electInput           int32
	driver               *services.Driver
	route                int
	routeString          string
	companyId            string
	deviceId             string
	platformId           string
	evts                 *eventstream.EventStream
	subs                 map[string]*eventstream.Subscription
	routes               map[int32]string
	itineraries          map[int32]*routes.Itinerary
	shcservices          map[string]*services.ScheduleService
	companySchServices   map[string]*services.ScheduleService
	currentSchServices   map[string]*services.ScheduleService
	// vehicleSchServices   map[string]*services.ScheduleService
	currentService  *services.ScheduleService
	selectedService *services.ScheduleService
	// companyCurrentSchServices map[string]*CompanySchService
	companySchServicesShow []*CompanySchService
	vehicleSchServicesShow []*CompanySchService
	notif                  []string
	uix                    ui.UI
	ctx                    actor.Context
	db                     *actor.PID
	pidApp                 *actor.PID
	pidSvc                 *actor.PID
	cancel                 func()
	cancelStep             func()
	cancelPop              func()
	behavior               actor.Behavior
	gps                    bool
	network                bool
	isDisplayEnable        bool
}

func NewApp(uix ui.UI) *App {
	a := &App{}
	a.companyId = "95ac8afe-f793-4354-b522-9459e5a96725"
	a.uix = uix
	a.network = true
	a.gps = true
	a.evts = eventstream.NewEventStream()
	a.routes = make(map[int32]string)
	a.shcservices = make(map[string]*services.ScheduleService)
	a.companySchServices = make(map[string]*services.ScheduleService)
	a.behavior = actor.NewBehavior()
	a.behavior.Become(a.Starting)
	a.deviceId = utils.Hostname()
	return a
}

func (a *App) RegisterActorService(pid *actor.PID) {
	a.pidSvc = pid
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

func (a *App) Receive(ctx actor.Context) {
	a.ctx = ctx
	switch msg := ctx.Message().(type) {
	case *params.Parameters:
		if msg != nil {
			if len(msg.COMPANY_ID) > 0 {
				a.companyId = msg.COMPANY_ID
			}
			// if len(msg.DEV_SERIAL) > 0 {
			// 	a.deviceId = msg.DEV_SERIAL
			// }
			if len(msg.DEV_PID) > 0 {
				a.platformId = msg.DEV_PID
			}
		}
	}
	a.behavior.Receive(ctx)
}

func (a *App) Starting(ctx actor.Context) {
	fmt.Printf("message: %q --> %q, %T\n", func() string {
		if ctx.Sender() == nil {
			return ""
		} else {
			return ctx.Sender().GetId()
		}
	}(), ctx.Self().GetId(), ctx.Message())

	switch ctx.Message().(type) {

	case *actor.Started:
		if a.cancel != nil {
			a.cancel()
		}

		contxt, cancel := context.WithCancel(context.Background())
		a.cancel = cancel
		go tick(contxt, ctx, TIMEOUT)
		ctx.Send(ctx.Self(), &tickMsg{})
	case *tickMsg:
		if a.uix != nil {
			time.Sleep(3 * time.Second)
			if err := a.uix.Init(); err != nil {
				logs.LogWarn.Printf("init error: %s", err)
				break
			}
		}
		a.behavior.Become(a.Runstate)
		ctx.Send(ctx.Self(), &startMsg{})

	case *actor.Stopping:
		if a.cancel != nil {
			a.cancel()
		}
		if a.cancelPop != nil {
			a.cancelPop()
		}
	}
}

func (a *App) Runstate(ctx actor.Context) {
	fmt.Printf("message: %q --> %q, %T\n", func() string {
		if ctx.Sender() == nil {
			return ""
		} else {
			return ctx.Sender().GetId()
		}
	}(), ctx.Self().GetId(), ctx.Message())

	switch msg := ctx.Message().(type) {
	case *startMsg:
		if err := func() error {
			db, err := database.Open(ctx, dbpath)
			if err != nil {
				return fmt.Errorf("open database  err: %s", err)
			}
			if db != nil {
				fmt.Println("database")
				a.db = db.PID()
				result, err := ctx.RequestFuture(a.db, &database.MsgGetData{
					Buckets: []string{collectionDriverData},
					ID:      "counters",
				}, 3000*time.Millisecond).Result()
				if err != nil {
					return fmt.Errorf("open database  err: %s", err)
				}
				switch v := result.(type) {
				case *database.MsgAckGetData:
					data := make([]byte, len(v.Data))
					copy(data, v.Data)
					res := &ValidationData{}
					if err := json.Unmarshal(data, res); err != nil {
						return fmt.Errorf("open database  err: %s", err)
					}
					fmt.Printf("recover database data: %+v\n", res)
					logs.LogInfo.Printf("recover database data: %s", res)

					if time.UnixMilli(res.Time).Day() != time.Now().Day() {
						a.countInput = 0
						a.countOutput = 0
						a.cashInput = 0
						a.electInput = 0
						logs.LogInfo.Printf("restart data: %v", &ValidationData{
							CountInputs: a.countInput, CashInputs: a.cashInput, CountOutputs: a.countOutput, ElectInputs: a.electInput})
					} else {
						a.countInput += res.CountInputs
						a.countOutput += res.CountOutputs
						a.cashInput += res.CashInputs
						a.electInput += res.ElectInputs
					}
					ctx.Send(ctx.Self(), &MsgShowCounters{})
				}
			}
			return nil
		}(); err != nil {
			logs.LogWarn.Println(err)
			ctx.Send(ctx.Self(), &ValidationData{})
		}
		fmt.Printf("********* /////////// subscribe to %q topic\n", constant.DISCOVERY_TOPIC)
		pubsub.Subscribe(constant.DISCOVERY_TOPIC, ctx.Self(), func(b []byte) interface{} {
			fmt.Printf("********* /////////// message arrive %q topic, msg: %s\n", constant.DISCOVERY_TOPIC, b)
			msg := new(msgdriverterminal.Discovery)
			if err := json.Unmarshal(b, msg); err != nil {
				logs.LogWarn.Printf("discovery parse error: %s", err)
				return err
			}
			return msg
		})

		contxt, cancel := context.WithCancel(context.Background())
		a.cancel = cancel
		go tick(contxt, ctx, TIMEOUT)

		a.isDisplayEnable = true

		a.uix.Gps(a.gps)
		a.uix.Network(a.network)

	case *actor.Stopping:
		if a.cancel != nil {
			a.cancel()
		}
		if a.cancelPop != nil {
			a.cancelPop()
		}
		if a.db != nil {
			data, err := json.Marshal(&ValidationData{
				CountInputs:  a.countInput,
				CountOutputs: a.countOutput,
				CashInputs:   a.cashInput,
				ElectInputs:  a.electInput,
				Time:         time.Now().UnixMilli(),
			})
			if err != nil {
				logs.LogWarn.Printf("database persistence error: %s", err)
				break
			}
			ctx.RequestFuture(a.db, &database.MsgUpdateData{
				ID:      "counters",
				Buckets: []string{collectionDriverData},
				Data:    data,
			}, time.Millisecond*3000).Wait()
			logs.LogInfo.Printf("backup database data: %s", data)
			fmt.Printf("backup database data: %s\n", data)
		}
		if a.db != nil {
			ctx.RequestFuture(a.db, &database.MsgCloseDB{}, 100*time.Millisecond).Wait()
		}
		if a.cancel != nil {
			a.cancel()
		}
		if a.cancelStep != nil {
			a.cancelStep()
		}
	case *msgdriverterminal.Discovery:
		if len(msg.GetAddress()) <= 0 || len(msg.GetId()) <= 0 {
			logs.LogWarn.Println("Discovery bad message")
			break
		}
		pid := actor.NewPID(msg.GetAddress(), msg.GetId())
		ctx.Request(pid, &msgdriverterminal.DiscoveryResponse{})
	case *MsgMainScreen:
		if err := a.mainScreen(); err != nil {
			logs.LogWarn.Println("mainScreen error: ", err)
		}
	case *StepMsg:
		if a.pidApp == nil && ctx.Parent() == nil {
			break
		}
		mss := &messages.MsgDriverPaso{}
		if a.pidApp != nil {
			ctx.Request(a.pidApp, mss)
		} else if ctx.Parent() != nil {
			ctx.Request(ctx.Parent(), mss)
		}
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
	case *messages.MsgAppPaso:

		if msg.GetCode() == messages.MsgAppPaso_CASH {
			a.cashInput += 1
			if a.uix != nil {
				if err := a.uix.CashInputs(int32(a.cashInput + a.electInput)); err != nil {
					logs.LogWarn.Printf("inputs error: %s", err)
				}
			}
		} else {
			a.electInput += 1
			if a.uix != nil {
				a.uix.Beep(3, 50, 600*time.Millisecond)
				if a.cancelPop != nil {
					a.cancelPop()
				}
				// if a.uix.GetScreen() == ui.MAIN_SCREEN {
				contxt, cancel := context.WithCancel(context.TODO())
				a.cancelPop = cancel
				if err := a.uix.TextConfirmationPopup(
					"entrada confirmada"); err != nil {
					logs.LogWarn.Printf("textConfirmation error: %s", err)
				}
				go func() {
					defer cancel()
					select {
					case <-contxt.Done():
					case <-time.After(4 * time.Second):
					}
					if err := a.uix.TextConfirmationPopupclose(); err != nil {
						logs.LogWarn.Printf("textConfirmation error: %s", err)
					}
				}()

				// }
				if err := a.uix.CashInputs(int32(a.electInput + a.cashInput)); err != nil {
					logs.LogWarn.Printf("inputs error: %s", err)
				}
			}
		}
	case *messages.MsgAppError:
		textBytes := make([]string, 0)
		v := msg.GetError()
		if len(v) > 26 {
			textBytes = append(textBytes, SplitHeader(v, 26)...)
		} else {
			textBytes = append(textBytes, v)
		}
		if a.uix != nil {
			a.uix.Beep(12, 90, 300*time.Millisecond)
			if a.cancelPop != nil {
				a.cancelPop()
			}
			// if a.uix.GetScreen() == ui.MAIN_SCREEN {
			contxt, cancel := context.WithCancel(context.TODO())
			a.cancelPop = cancel
			if err := a.uix.TextWarningPopup(textBytes...); err != nil {
				logs.LogWarn.Printf("textWarningPopup error: %s", err)
				break
			}
			go func() {
				defer cancel()
				select {
				case <-contxt.Done():
				case <-time.After(4 * time.Second):
				}
				if err := a.uix.TextWarningPopupClose(); err != nil {
					logs.LogWarn.Printf("textWarningPopupClose error: %s", err)
				}
			}()
			// }
			fmt.Printf("screen in warn: %v\n", a.uix.GetScreen())
		}
	case *tickResetCountersMsg:
		a.countInput = 0
		a.countOutput = 0
		a.cashInput = 0
		a.electInput = 0
		ctx.Send(ctx.Self(), &MsgShowCounters{})
	case *MsgShowCounters:
		if a.uix != nil {
			if err := a.uix.CashInputs(int32(a.cashInput + a.electInput)); err != nil {
				logs.LogWarn.Printf("inputs error: %s", err)
			}
		}
	case *ValidationData:
		a.countInput += msg.CountInputs
		a.countOutput += msg.CountOutputs
		a.cashInput += msg.CashInputs
		a.electInput += msg.ElectInputs
		if a.uix != nil {
			if err := a.uix.CashInputs(int32(a.cashInput + a.electInput)); err != nil {
				logs.LogWarn.Printf("inputs error: %s", err)
			}
		}

	case *MsgDoors:
		if len(msg.Value) <= 0 {
			break
		}
		a.evts.Publish(msg)
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
		a.pidApp = ctx.Sender()

	// case *messages.MsgAppMapRoute:
	// 	fmt.Printf("message -> \"%v\"\n", msg)
	// 	a.routes = msg.Routes
	case *routes.ItinerariesMsg:
		itis := make(map[int32]*routes.Itinerary)
		for _, v := range msg.Itineraries {
			itis[v.PaymentMediumCode] = v
		}
		if len(itis) > 0 {
			a.itineraries = itis
		}
	case *messages.MsgRoute:
		fmt.Printf("message -> \"%v\"\n", msg)
		v := msg.GetItineraryName()
		a.route = int(msg.GetItineraryCode())
		if a.routeString == v {
			break
		}
		a.routeString = v
		a.uix.Beep(3, 50, 600*time.Millisecond)

		// routeS := fmt.Sprintf("%d", a.route)

		if err := a.uix.Route(a.routeString); err != nil {
			logs.LogWarn.Printf("route error: %s", err)
		}
	case *MsgSetRoute:
		a.uix.Beep(3, 50, 600*time.Millisecond)
		a.route = msg.Route
		if a.pidApp != nil {
			mss := &messages.MsgSetRoute{
				Code: int32(a.route),
			}
			ctx.Request(a.pidApp, mss)
		}
		if len(msg.RouteName) > 0 {
			a.routeString = msg.RouteName
			if err := a.uix.Route(msg.RouteName); err != nil {
				logs.LogWarn.Printf("error Route: %s", err)
			}
		}
	case *MsgSetDriver:
		if err := a.setDriver(ctx, msg); err != nil {
			logs.LogWarn.Println("setDriver error: ", err)
		} else {
			if a.pidApp != nil && msg.Driver > 0 {
				mss := &messages.MsgSetDriver{
					Code: int32(msg.Driver),
				}
				ctx.Request(a.pidApp, mss)
			}
		}
	case *MsgSetRoutes:
		a.routes = msg.Routes
	case *MsgConfirmationText:
		if err := a.uix.TextConfirmation(string(msg.Text)); err != nil {
			logs.LogWarn.Printf("textConfirmation error: %s", err)
		}
	case *MsgConfirmationTextMainScreen:
		if a.uix != nil && a.uix.GetScreen() == ui.MAIN_SCREEN {
			a.uix.Beep(3, 50, 600*time.Millisecond)
			// if a.uix.GetScreen() == ui.MAIN_SCREEN {
			if a.cancelPop != nil {
				a.cancelPop()
			}
			contxt, cancel := context.WithCancel(context.TODO())
			a.cancelPop = cancel
			if err := a.uix.TextConfirmationPopup(
				"entrada confirmada"); err != nil {
				logs.LogWarn.Printf("textConfirmation error: %s", err)
			}
			go func() {
				defer cancel()
				select {
				case <-contxt.Done():
				case <-time.After(4 * time.Second):
					if err := a.uix.TextConfirmationPopupclose(); err != nil {
						logs.LogWarn.Printf("textConfirmation error: %s", err)
					}
				}
			}()
			// }
			fmt.Printf("screen in text: %v\n", a.uix.GetScreen())
		}
	case *MsgWarningText:
		if err := a.uix.TextWarning(string(msg.Text)); err != nil {
			logs.LogWarn.Printf("textConfirmation error: %s", err)
		}
	case *ListProgVeh:
		if err := a.listProg(msg); err != nil {
			logs.LogWarn.Println("listProg error: ", err)
			if err := a.uix.TextWarningPopup(err.Error()); err != nil {
				logs.LogWarn.Printf("textWarningPopup error: %s", err)
			}
			if a.cancelPop != nil {
				a.cancelPop()
			}
			contxt, cancel := context.WithCancel(context.Background())
			a.cancelPop = cancel
			go func() {
				defer cancel()
				select {
				case <-contxt.Done():
				case <-time.After(4 * time.Second):
				}
				if err := a.uix.TextWarningPopupClose(); err != nil {
					logs.LogWarn.Printf("textWarningPopupClose error: %s", err)
				}
			}()
		}
	case *RequestProgVeh:
		if err := a.requestProg(ctx, msg); err != nil {
			logs.LogWarn.Println("requestProg error: ", err)
		}
	case *ListProgDriver:
		if err := a.listDriverProg(msg); err != nil {
			if err := a.uix.TextWarningPopup(fmt.Sprintf("%s\n", err)); err != nil {
				logs.LogWarn.Printf("textWarningPopup error: %s", err)
			}
			if a.cancelPop != nil {
				a.cancelPop()
			}
			contxt, cancel := context.WithCancel(context.TODO())
			a.cancelPop = cancel
			go func() {
				defer cancel()
				select {
				case <-contxt.Done():
				case <-time.After(4 * time.Second):
				}
				if err := a.uix.TextWarningPopupClose(); err != nil {
					logs.LogWarn.Printf("textWarningPopupClose error: %s", err)
				}
			}()
			logs.LogWarn.Println("requestProg error: ", err)
		}
	case *RequestTakeService:
		if err := a.takeservice(); err != nil {
			if err := a.uix.TextWarningPopup(fmt.Sprintf("%s\n", err)); err != nil {
				logs.LogWarn.Printf("textWarningPopup error: %s", err)
			}
			if a.cancelPop != nil {
				a.cancelPop()
			}
			contxt, cancel := context.WithCancel(context.TODO())
			a.cancelPop = cancel
			go func() {
				defer cancel()
				select {
				case <-contxt.Done():
				case <-time.After(4 * time.Second):
				}
				if err := a.uix.TextWarningPopupClose(); err != nil {
					logs.LogWarn.Printf("textWarningPopupClose error: %s", err)
				}
			}()
			logs.LogWarn.Printf("takeservice error: %s", err)
		}
	case *services.StatusSch:
		fmt.Printf("******** (%T) %v **********\n", msg, msg)
		if msg.State == 0 && a.network {
			if err := a.uix.Network(true); err != nil {
				logs.LogWarn.Printf("network error: %s", err)
			}
		} else if msg.State == 1 && !a.network {
			if err := a.uix.Network(false); err != nil {
				logs.LogWarn.Printf("network error: %s", err)
			}
		}
		a.network = func() bool { return msg.State != 0 }()
	case *services.Snapshot:
		fmt.Printf("******** (%T) %v **********\n", msg, msg)
		if len(msg.ScheduledServices) > 0 {
			a.shcservices = make(map[string]*services.ScheduleService)
			for _, v := range msg.ScheduledServices {
				a.shcservices[v.GetId()] = v
			}
		}
	case *services.UpdateServiceMsg:
		fmt.Printf("******** (%T) %v **********\n", msg, msg)
		svc := msg.GetUpdate()
		a.showCurrentService(svc)
		if a.currentService != nil &&
			len(a.currentService.GetItinerary().GetName()) > 0 {
			if !strings.EqualFold(a.currentService.GetItinerary().GetName(), a.routeString) {
				a.ctx.Send(a.ctx.Self(), &MsgSetRoute{
					Route:     int(a.currentService.GetItinerary().GetId()),
					RouteName: a.currentService.GetItinerary().GetName(),
				})
				a.routeString = a.currentService.GetItinerary().GetName()
				if err := a.uix.Route(a.currentService.GetItinerary().GetName()); err != nil {
					logs.LogWarn.Printf("route error: %s", err)
				}
			}
		}
	case *services.ServiceMsg:
		fmt.Printf("******** (%T) %v **********\n", msg, msg)
		svc := msg.GetUpdate()
		if len(svc.GetState()) > 0 && svc.GetScheduleDateTime() > 0 {
			a.shcservices[svc.GetId()] = svc
		} else {
			if v, ok := a.shcservices[svc.GetId()]; ok {
				fmt.Printf("////// route: %v\n", v)
				UpdateService(v, svc)
			} else if svc.GetScheduleDateTime() > 0 {
				a.shcservices[svc.GetId()] = svc
			}
		}
	case *services.ServiceAllMsg:
		fmt.Printf("******** (%T) %v **********\n", msg, msg)
		a.showCurrentServiceWithAll(msg)
		if a.currentService != nil &&
			len(a.currentService.GetItinerary().GetName()) > 0 {
			if strings.EqualFold(a.currentService.GetItinerary().GetName(), a.routeString) {
				if err := a.uix.Route(a.currentService.GetItinerary().GetName()); err != nil {
					logs.LogWarn.Printf("route error: %s", err)
				}
			}
		}
	case *services.RemoveServiceMsg:
		svc := msg.GetUpdate()
		if svc.GetCheckpointTimingState() != nil && len(svc.GetCheckpointTimingState().GetState()) > 0 {
			state := int(services.TimingState_value[svc.GetCheckpointTimingState().GetState()])
			promtp := fmt.Sprintf("%s (%d / fin)", svc.GetCheckpointTimingState().GetName(), svc.GetCheckpointTimingState().GetTimeDiff())
			fmt.Printf("///// state: %d\n", state)
			if err := a.uix.ServiceCurrentState(0, promtp); err != nil {
				logs.LogWarn.Printf("textConfirmation error: %s", err)
			}
		} else {
			if err := a.uix.ServiceCurrentState(0, ""); err != nil {
				logs.LogWarn.Printf("textConfirmation error: %s", err)
			}
		}
		delete(a.shcservices, svc.GetId())
	case *MsgScreen:
		if err := a.uix.Screen(msg.ID, msg.Switch); err != nil {
			logs.LogWarn.Printf("msgScreen error: %s", err)
		}
	case *counterpass.CounterEvent:
		fmt.Printf("counter event: %v\n", msg)
		a.countInput += int32(msg.Inputs)
		if msg.Inputs > 0 {
			if err := a.uix.Inputs(int32(a.countInput)); err != nil {
				logs.LogWarn.Printf("inputs error: %s", err)
			}
		}
		a.countOutput += int32(msg.Outputs)
		if msg.Outputs > 0 {
			if err := a.uix.Outputs(int32(a.countOutput)); err != nil {
				logs.LogWarn.Printf("outputs error: %s", err)
			}
		}

	case *counterpass.CounterExtraEvent:
		if len(msg.Text) > 0 {
			if err := a.uix.TextWarningPopup(string(msg.Text)); err != nil {
				logs.LogWarn.Printf("textWarningPopup error: %s", err)
			}
			if a.cancelPop != nil {
				a.cancelPop()
			}
			contxt, cancel := context.WithCancel(context.TODO())
			a.cancelPop = cancel
			go func() {
				defer cancel()
				select {
				case <-contxt.Done():
				case <-time.After(4 * time.Second):
				}
				if err := a.uix.TextWarningPopupClose(); err != nil {
					logs.LogWarn.Printf("textWarningPopupClose error: %s", err)
				}
			}()
		}
	case *counterpass.CounterMap:
	case *counterpass.TurnstileRegisters:
		fmt.Printf("counter turnstile registers: %v\n", msg)
		if len(msg.Registers) < 3 {
			break
		}
		totalCountInput := int32(msg.Registers[0])
		if a.totalCountInput < totalCountInput {
			if a.totalCountInput == 0 {
				a.totalCountInput = totalCountInput
			}
			// fmt.Printf("%d += %d - %d\n", a.countInput, a.totalCountInput, totalCountInput)
			a.countInput += totalCountInput - a.totalCountInput
			fmt.Printf("input count: %d\n", a.countInput)
			if err := a.uix.Inputs(int32(a.countInput)); err != nil {
				logs.LogWarn.Printf("inputs error: %s", err)
			}
		}
		a.totalCountInput = totalCountInput
	case *messages.MsgGpsOk:
		a.gps = false
		if err := a.uix.Gps(false); err != nil {
			logs.LogWarn.Printf("gps error: %s", err)
		}
	case *messages.MsgGpsErr:
		a.gps = true
		if err := a.uix.Gps(true); err != nil {
			logs.LogWarn.Printf("gps error: %s", err)
		}
	case *messages.MsgGroundOk:
		a.network = false
		if err := a.uix.Network(false); err != nil {
			logs.LogWarn.Printf("network error: %s", err)
		}
	case *messages.MsgGroundErr:
		a.network = true
		if err := a.uix.Network(true); err != nil {
			logs.LogWarn.Printf("network error: %s", err)
		}
	case *MsgUpdateTime:
		tNow := time.Now()
		if a.updateTime.Minute() == tNow.Minute() && a.updateTime.Hour() == tNow.Hour() {
			break
		}
		a.updateTime = tNow
		if err := a.uix.Date(tNow); err != nil {
			logs.LogWarn.Printf("date error: %s", err)
		}
	case *TestTextProgDriver:
		if err := a.uix.ShowProgDriver(msg.Text...); err != nil {
			logs.LogWarn.Printf("textProgVehicle error: %s", err)
		}
	case *ErrorDisplay:
		if msg.Error == nil {
			break
		}
		if time.Since(a.lastErrButtons.Timestamp) > 3*time.Minute ||
			a.lastErrButtons.Error.Error() != msg.Error.Error() {
			a.lastErrButtons = errorDisplay{
				Timestamp: time.Now(),
				Error:     msg.Error,
			}
			logs.LogWarn.Printf("error display: %s", msg.Error)

		}
		if a.isDisplayEnable {
			a.isDisplayEnable = false
		}
		fmt.Println(msg.Error)
	case *tickMsg:
		if a.uix == nil {
			break
		}
		if err := a.uix.VerifyDisplay(); err != nil {
			if time.Since(a.lastErrVerifyDisplay.Timestamp) > 3*time.Minute ||
				a.lastErrVerifyDisplay.Error.Error() != err.Error() {
				a.lastErrVerifyDisplay = errorDisplay{
					Timestamp: time.Now(),
					Error:     err,
				}
				logs.LogWarn.Printf("error display: %s", err)
			}
			if a.isDisplayEnable {
				fmt.Printf("///////////////////// false")
				a.isDisplayEnable = false
			}
		} else if !a.isDisplayEnable {
			// fmt.Printf("///////////////////// true")
			a.isDisplayEnable = true
			ctx.Send(ctx.Self(), &MsgMainScreen{})
		}
	case error:
		fmt.Printf("error message: %s (%s)\n", msg, ctx.Self().GetId())
	default:
		fmt.Printf("unhandled message type: %T (%s)\n", msg, ctx.Self().GetId())
	}
}

type tickResetCountersMsg struct{}
type tickMsg struct{}
type startMsg struct{}

// TODO: comment out for test
// var tRefg time.Time

func tick(contxt context.Context, ctx actor.Context, timeout time.Duration) {

	self := ctx.Self()
	ctxroot := ctx.ActorSystem().Root

	go func() {

		tn := time.Now()
		var until time.Duration
		// TODO: comment out for test
		// fmt.Printf("////////// time: %s\n", tRefg.Sub(time.Time{}))
		t := time.Date(tn.Year(), tn.Month(), tn.Day(), 00, 01, 59, 0, tn.Location())
		until = time.Until(t)
		t2 := time.NewTicker(timeout)
		defer t2.Stop()
		t1 := time.NewTimer(until)
		defer t1.Stop()
		for {
			select {
			case <-contxt.Done():
				return
			case <-t1.C:
				ctxroot.Send(self, &tickResetCountersMsg{})
			case <-t2.C:
				ctxroot.Send(self, &tickMsg{})
			}
		}
	}()
}
