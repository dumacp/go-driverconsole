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
	"github.com/dumacp/go-driverconsole/internal/buttons"
	"github.com/dumacp/go-driverconsole/internal/constant"
	"github.com/dumacp/go-driverconsole/internal/counterpass"
	"github.com/dumacp/go-driverconsole/internal/pubsub"
	"github.com/dumacp/go-driverconsole/internal/ui"
	msgdriverterminal "github.com/dumacp/go-driverconsole/pkg/messages"
	"github.com/dumacp/go-fareCollection/pkg/messages"
	"github.com/dumacp/go-logs/pkg/logs"
	"github.com/dumacp/go-schservices/api/services"
)

const (
	dbpath               = "/SD/boltdbs/driverdb"
	databaseName         = "driverdb"
	collectionAppData    = "appfaredata"
	collectionDriverData = "terminaldata"

	TIMEOUT = 30 * time.Second
)

type App struct {
	updateTime           time.Time
	lastErrVerifyDisplay errorDisplay
	lastErrButtons       errorDisplay
	totalCountInput      int32
	totalCountOutput     int32
	countInput           int32
	countOutput          int32
	cashInput            int32
	electInput           int32
	timeLapse            int
	driver               int
	route                int
	brightness           int
	routeString          string
	routeCode            string
	driverCode           string
	displayVendor        string
	evts                 *eventstream.EventStream
	subs                 map[string]*eventstream.Subscription
	routes               map[int32]string
	shcservices          map[string]*services.ScheduleService
	notif                []string
	// evt2evtApp       map[buttons.KeyCode]EventLabel
	uix             ui.UI
	ctx             actor.Context
	db              *actor.PID
	pidApp          *actor.PID
	contxt          context.Context
	buttonDevice    buttons.ButtonDevice
	cancel          func()
	cancelStep      func()
	renewStep       func()
	cancelPop       func()
	gps             bool
	network         bool
	enableStep      bool
	isDisplayEnable bool
}

func NewApp(uix ui.UI, displayVendor string) *App {
	a := &App{}
	a.uix = uix
	a.displayVendor = displayVendor
	a.network = true
	a.gps = true
	a.evts = eventstream.NewEventStream()
	a.routes = make(map[int32]string)
	a.shcservices = make(map[string]*services.ScheduleService)
	// a.evt2evtApp = button2buttonApp
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

func (a *App) Receive(ctx actor.Context) {

	a.ctx = ctx
	fmt.Printf("message: %q --> %q, %T\n", func() string {
		if ctx.Sender() == nil {
			return ""
		} else {
			return ctx.Sender().GetId()
		}
	}(), ctx.Self().GetId(), ctx.Message())

	switch msg := ctx.Message().(type) {
	case *actor.Started:
		db, err := database.Open(ctx, dbpath)
		if err != nil {
			logs.LogWarn.Printf("open database  err: %s\n", err)
		}
		if db != nil {
			fmt.Println("database")
			a.db = db.PID()
			ctx.Request(a.db, &database.MsgQueryData{
				Buckets:  []string{collectionAppData},
				PrefixID: "",
				Reverse:  false,
			})
			ctx.Request(a.db, &database.MsgQueryData{
				Buckets:  []string{collectionDriverData},
				PrefixID: "",
				Reverse:  false,
			})
		}

		// if err := a.uix.Inputs(int64(a.countInput)); err != nil {
		// 	logs.LogWarn.Printf("inputs error: %s", err)
		// }
		// if err := a.uix.Outputs(int64(a.countOutput)); err != nil {
		// 	logs.LogWarn.Printf("outputs error: %s", err)
		// }
		// if err := a.uix.DeviationInputs(int64(0)); err != nil {
		// 	logs.LogWarn.Printf("devitation error: %s", err)
		// }
		ctx.Send(ctx.Self(), &ValidationData{})

		if a.uix != nil {
			time.Sleep(1 * time.Second)
			if err := a.uix.Init(); err != nil {
				logs.LogWarn.Printf("init error: %s", err)
			}
			time.Sleep(3 * time.Second)
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
			}, time.Millisecond*100).Wait()
			fmt.Printf("backup database data: %s\n", data)
		}
		if a.db != nil {
			ctx.Send(a.db, &database.MsgCloseDB{})
		}
		if a.cancel != nil {
			a.cancel()
		}
		if a.cancelStep != nil {
			a.cancelStep()
		}
	case MsgMainScreen:
		ctx.Send(ctx.Self(), MsgMainScreenForce{})
	case *MsgMainScreenForce:
		fmt.Printf("**** screen: %d\n", a.uix.GetScreen())
		if !msg.Force && a.uix.GetScreen() == ui.MAIN_SCREEN {
			break
		}
		if err := func() error {
			if err := a.uix.MainScreen(); err != nil {
				return fmt.Errorf("main screen error: %s", err)
			}
			if err := a.uix.ElectronicInputs(int32(a.electInput)); err != nil {
				return fmt.Errorf("electInput error: %s", err)
			}
			if err := a.uix.CashInputs(int32(a.cashInput)); err != nil {
				return fmt.Errorf("cashInput error: %s", err)
			}
			if err := a.uix.Inputs(int32(a.countInput)); err != nil {
				return fmt.Errorf("countInput error: %s", err)
			}
			if err := a.uix.Outputs(int32(a.countOutput)); err != nil {
				return fmt.Errorf("countOutput error: %s", err)
			}
			if err := a.uix.DateWithFormat(a.updateTime, "2006/01/02 15:04"); err != nil {
				return fmt.Errorf("date error: %s", err)
			}
			if err := a.uix.Driver(fmt.Sprintf("%d", a.driver)); err != nil {
				return fmt.Errorf("driver error: %s", err)
			}
			if len(a.routeString) > 0 {
				routeS := func() string {
					if len(a.routeString) > 32 {
						return a.routeString[:32]
					}
					return a.routeString
				}()
				if err := a.uix.Route(routeS); err != nil {
					return fmt.Errorf("route error: %s", err)
				}
			}
			return nil
		}(); err != nil {
			logs.LogWarn.Println(err)
		}
	case *MsgSetDriver:
		if a.pidApp != nil {
			mss := &messages.MsgSetDriver{
				Code: int32(msg.Driver),
			}
			ctx.Request(a.pidApp, mss)
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
		if a.uix.GetScreen() != ui.MAIN_SCREEN {
			break
		}
		routeS := func() string {
			if len(a.routeString) > 32 {
				return a.routeString[:32]
			}
			return a.routeString
		}()

		if err := a.uix.Route(routeS); err != nil {
			logs.LogWarn.Printf("route error: %s", err)
		}
	case *MsgSetRoute:
		// if len(a.routes) <= 0 {
		// 	break
		// }
		// a.route = msg.Route
		// a.routeString = a.routes[int32(msg.Route)]
		// if err := a.uix.Route(a.routeString); err != nil {
		// 	logs.LogWarn.Printf("route error: %s", err)
		// }
		a.uix.Beep(3, 50, 600*time.Millisecond)
		if a.pidApp != nil {
			mss := &messages.MsgSetRoute{
				Code: int32(a.route),
			}
			ctx.Request(a.pidApp, mss)
		}
	case *msgdriverterminal.Discovery:
		if len(msg.GetAddress()) <= 0 || len(msg.GetId()) <= 0 {
			logs.LogWarn.Println("Discovery bad message")
			break
		}
		pid := actor.NewPID(msg.GetAddress(), msg.GetId())
		ctx.Request(pid, &msgdriverterminal.DiscoveryResponse{})
	case *StepMsg:
		if a.pidApp == nil {
			break
		}
		mss := &messages.MsgDriverPaso{}
		ctx.Request(a.pidApp, mss)
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
		fmt.Printf("**** screen: %d\n", a.uix.GetScreen())
		if msg.GetCode() == messages.MsgAppPaso_CASH {
			a.cashInput += 1
			if a.uix != nil {
				if err := a.uix.CashInputs(int32(a.cashInput)); err != nil {
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
				if a.uix.GetScreen() == ui.MAIN_SCREEN {
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
						if strings.EqualFold(GTT50, a.displayVendor) {
							if err := a.uix.ElectronicInputs(int32(a.electInput)); err != nil {
								logs.LogWarn.Printf("electInput error: %s", err)
							}
							if err := a.uix.CashInputs(int32(a.cashInput)); err != nil {
								logs.LogWarn.Printf("cashInput error: %s", err)
							}
							if err := a.uix.Inputs(int32(a.countInput)); err != nil {
								logs.LogWarn.Printf("countInput error: %s", err)
							}
							if err := a.uix.Outputs(int32(a.countOutput)); err != nil {
								logs.LogWarn.Printf("countOutput error: %s", err)
							}
						}
					}()

				}
				if !strings.EqualFold(GTT50, a.displayVendor) {
					if err := a.uix.ElectronicInputs(int32(a.electInput)); err != nil {
						logs.LogWarn.Printf("inputs error: %s", err)
					}
				}
			}
		}
	case *messages.MsgAppError:
		fmt.Printf("**** screen: %d\n", a.uix.GetScreen())
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
			if a.uix.GetScreen() == ui.MAIN_SCREEN {
				contxt, cancel := context.WithCancel(context.TODO())
				a.cancelPop = cancel
				if err := a.uix.TextWarningPopup(
					textBytes...); err != nil {
					logs.LogWarn.Printf("textWarningPopup error: %s", err)
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
					if strings.EqualFold(a.displayVendor, GTT50) {
						if err := a.uix.ElectronicInputs(int32(a.electInput)); err != nil {
							logs.LogWarn.Printf("electInput error: %s", err)
						}
						if err := a.uix.CashInputs(int32(a.cashInput)); err != nil {
							logs.LogWarn.Printf("cashInput error: %s", err)
						}
						if err := a.uix.Inputs(int32(a.countInput)); err != nil {
							logs.LogWarn.Printf("countInput error: %s", err)
						}
						if err := a.uix.Outputs(int32(a.countOutput)); err != nil {
							logs.LogWarn.Printf("countOutput error: %s", err)
						}
					}
				}()

			}
		}
	case *tickResetCountersMsg:
		if a.db != nil {
			if err := func() error {
				a.countInput = 0
				a.countOutput = 0
				data, err := json.Marshal(&ValidationData{
					CountInputs:  0,
					CountOutputs: 0,
					CashInputs:   0,
					ElectInputs:  0,
					Time:         time.Now().UnixMilli(),
				})
				if err != nil {
					return fmt.Errorf("database persistence error: %s", err)
				}
				ctx.RequestFuture(a.db, &database.MsgUpdateData{
					ID:      "counters",
					Buckets: []string{collectionDriverData},
					Data:    data,
				}, time.Millisecond*100).Wait()
				fmt.Printf("backup database data: %s\n", data)
				return nil
			}(); err != nil {
				logs.LogWarn.Printf("error updating data from database: %s", err)
			}
		}
	case *database.MsgQueryResponse:
		fmt.Printf("form database: %q\n", msg)
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
			fmt.Printf("form database: %q\n", data)
			coll := msg.Buckets[len(msg.Buckets)-1]
			switch coll {
			case collectionAppData:
			case collectionDriverData:
				fmt.Printf("***** data (%s) restore: %s\n", msg.ID, data)
				res := &ValidationData{}
				if err := json.Unmarshal(data, res); err != nil {
					return err
				}
				if time.UnixMilli(res.Time).Day() != time.Now().Day() {
					if a.db != nil {
						data, err := json.Marshal(&ValidationData{
							CountInputs:  a.countInput,
							CountOutputs: a.countOutput,
							CashInputs:   a.cashInput,
							ElectInputs:  a.electInput,
							Time:         time.Now().UnixMilli(),
						})
						if err != nil {
							return fmt.Errorf("database persistence error: %s", err)
						}
						ctx.RequestFuture(a.db, &database.MsgUpdateData{
							ID:      "counters",
							Buckets: []string{collectionDriverData},
							Data:    data,
						}, time.Millisecond*100).Wait()
						fmt.Printf("backup database data: %s\n", data)
					}
					break
				}
				fmt.Printf("recover database data: %v\n", res)

				ctx.Send(ctx.Self(), res)

			}
			return nil
		}(); err != nil {
			logs.LogWarn.Printf("error updating data from database: %s", err)
		}
	case *ValidationData:
		a.countInput += msg.CountInputs
		a.countOutput += msg.CountOutputs
		a.cashInput += msg.CashInputs
		a.electInput += msg.ElectInputs
		if a.uix != nil && a.uix.GetScreen() == ui.MAIN_SCREEN {
			if err := a.uix.CashInputs(int32(a.cashInput)); err != nil {
				logs.LogWarn.Printf("inputs error: %s", err)
			}
			if err := a.uix.ElectronicInputs(int32(a.electInput)); err != nil {
				logs.LogWarn.Printf("outputs error: %s", err)
			}
			if err := a.uix.Inputs(int32(a.countInput)); err != nil {
				logs.LogWarn.Printf("countInput error: %s", err)
			}
			if err := a.uix.Outputs(int32(a.countOutput)); err != nil {
				logs.LogWarn.Printf("countOutput error: %s", err)
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
	case *messages.MsgAppMapRoute:
		fmt.Printf("message -> \"%v\"\n", msg)
		a.routes = msg.Routes
	case *MsgSetRoutes:
		fmt.Printf("message -> \"%v\"\n", msg)
		a.routes = msg.Routes
	case *MsgConfirmationText:
		if err := a.uix.TextConfirmation(string(msg.Text)); err != nil {
			logs.LogWarn.Printf("textConfirmation error: %s", err)
		}
	case *MsgWarningTextInMainScreen:
	case *MsgConfirmationTextMainScreen:
		if a.uix == nil || a.uix.GetScreen() != ui.MAIN_SCREEN {
			break
		}
		if a.cancelPop != nil {
			a.cancelPop()
		}
		contxt, cancel := context.WithCancel(context.TODO())
		a.cancelPop = cancel
		if err := a.uix.TextConfirmationPopup(string(msg.Text)); err != nil {
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
				if strings.EqualFold(GTT50, a.displayVendor) {
					if err := a.uix.ElectronicInputs(int32(a.electInput)); err != nil {
						logs.LogWarn.Printf("electInput error: %s", err)
					}
					if err := a.uix.CashInputs(int32(a.cashInput)); err != nil {
						logs.LogWarn.Printf("cashInput error: %s", err)
					}
					if err := a.uix.Inputs(int32(a.countInput)); err != nil {
						logs.LogWarn.Printf("countInput error: %s", err)
					}
					if err := a.uix.Outputs(int32(a.countOutput)); err != nil {
						logs.LogWarn.Printf("countOutput error: %s", err)
					}
				}
			}
		}()
	case *MsgWarningText:
		if err := a.uix.TextWarning(string(msg.Text)); err != nil {
			logs.LogWarn.Printf("textWarning error: %s", err)
		}
	case *services.StatusSch:
		fmt.Printf("******** %v **********", msg)
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
	case *services.UpdateServiceMsg:

		svc := msg.GetUpdate()
		if len(svc.GetState()) > 0 {
			if v, ok := a.shcservices[svc.GetId()]; ok {
				// fmt.Printf("////// route: %v\n", v)
				UpdateService(v, msg.GetUpdate())
				fmt.Printf("////// update: %v\n", v)
			} else {
				a.shcservices[svc.GetId()] = msg.GetUpdate()
			}
			v := a.shcservices[svc.GetId()]
			data := strings.ToLower(fmt.Sprintf(" %s: (%d) %s (%s)", time.Now().Format("01/02 15:04"),
				v.GetItinerary().GetId(), v.GetItinerary().GetName(), v.GetState()))
			a.notif = append(a.notif, data)
			if len(a.notif) > 10 {
				copy(a.notif, a.notif[1:])
				a.notif = a.notif[:len(a.notif)-1]
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
		fmt.Printf("/////////////// services (ori: %d): %v\n", len(msg.GetUpdates()), arr)
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
		// a.countInput += int32(msg.Inputs)
		// if msg.Inputs > 0 {
		// 	if err := a.uix.Inputs(int32(a.countInput)); err != nil {
		// 		logs.LogWarn.Printf("inputs error: %s", err)
		// 	}
		// }
		// a.countOutput += int32(msg.Outputs)
		// if msg.Outputs > 0 {
		// 	if err := a.uix.Outputs(int32(a.countOutput)); err != nil {
		// 		logs.LogWarn.Printf("outputs error: %s", err)
		// 	}
		// }

	case *counterpass.CounterExtraEvent:
		if a.uix.GetScreen() != ui.MAIN_SCREEN {
			break
		}
		if len(msg.Text) > 0 {
			if a.cancelPop != nil {
				a.cancelPop()
			}
			contxt, cancel := context.WithCancel(context.TODO())
			a.cancelPop = cancel
			if err := a.uix.TextWarningPopup(
				string(msg.Text)); err != nil {
				logs.LogWarn.Printf("textWarningPopup error: %s", err)
			}
			go func() {
				defer cancel()
				select {
				case <-contxt.Done():
				case <-time.After(4 * time.Second):
					if err := a.uix.TextWarningPopupClose(); err != nil {
						logs.LogWarn.Printf("textWarningPopupClose error: %s", err)
					}
					if strings.EqualFold(a.displayVendor, GTT50) {
						if err := a.uix.ElectronicInputs(int32(a.electInput)); err != nil {
							logs.LogWarn.Printf("electInput error: %s", err)
						}
						if err := a.uix.CashInputs(int32(a.cashInput)); err != nil {
							logs.LogWarn.Printf("cashInput error: %s", err)
						}
						if err := a.uix.Inputs(int32(a.countInput)); err != nil {
							logs.LogWarn.Printf("countInput error: %s", err)
						}
						if err := a.uix.Outputs(int32(a.countOutput)); err != nil {
							logs.LogWarn.Printf("countOutput error: %s", err)
						}
					}
				}
			}()
		}
	case *counterpass.CounterMap:
		totalCountInput := int32(msg.Inputs0 + msg.Inputs1)
		if a.totalCountInput < totalCountInput {
			if a.totalCountInput == 0 {
				a.totalCountInput = totalCountInput
			}
			a.countInput += a.totalCountInput - totalCountInput
			if err := a.uix.Inputs(int32(a.countInput)); err != nil {
				logs.LogWarn.Printf("inputs error: %s", err)
			}
		}
		a.totalCountInput = totalCountInput
		totalCountOutput := int32(msg.Outputs0 + msg.Outputs1)
		if a.totalCountOutput < totalCountOutput {
			if a.totalCountOutput == 0 {
				a.totalCountOutput = totalCountInput
			}
			a.countOutput += a.totalCountOutput - totalCountOutput
			if err := a.uix.Outputs(a.countOutput); err != nil {
				logs.LogWarn.Printf("outputs error: %s", err)
			}
		}
		a.totalCountOutput = totalCountOutput
	case *messages.MsgGpsOk:
		a.gps = false
		if a.uix.GetScreen() != ui.MAIN_SCREEN && strings.EqualFold(a.displayVendor, GTT50) {
			break
		}
		if err := a.uix.Gps(false); err != nil {
			logs.LogWarn.Printf("gps error: %s", err)
		}
	case *messages.MsgGpsErr:
		a.gps = true
		if a.uix.GetScreen() != ui.MAIN_SCREEN && strings.EqualFold(a.displayVendor, GTT50) {
			break
		}
		if err := a.uix.Gps(true); err != nil {
			logs.LogWarn.Printf("gps error: %s", err)
		}
	case *messages.MsgGroundOk:
		a.network = false
		if a.uix.GetScreen() != ui.MAIN_SCREEN && strings.EqualFold(a.displayVendor, GTT50) {
			break
		}
		if err := a.uix.Network(false); err != nil {
			logs.LogWarn.Printf("network error: %s", err)
		}
	case *messages.MsgGroundErr:
		a.network = true
		if a.uix.GetScreen() != ui.MAIN_SCREEN && strings.EqualFold(a.displayVendor, GTT50) {
			break
		}
		if err := a.uix.Network(true); err != nil {
			logs.LogWarn.Printf("network error: %s", err)
		}
	case *MsgUpdateTime:
		tNow := time.Now()
		if a.updateTime.Minute() == tNow.Minute() && a.updateTime.Hour() == tNow.Hour() {
			break
		}
		a.updateTime = tNow
		if a.uix.GetScreen() != ui.MAIN_SCREEN && strings.EqualFold(a.displayVendor, GTT50) {
			break
		}
		if err := a.uix.DateWithFormat(tNow, "2006/01/02 15:04"); err != nil {
			logs.LogWarn.Printf("date error: %s", err)
		}
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
				fmt.Printf("/////////////// false\n")
				a.isDisplayEnable = false
			}
		} else if !a.isDisplayEnable {
			// fmt.Printf("/////////////// true\n")
			a.isDisplayEnable = true
			ctx.Send(ctx.Self(), &MsgMainScreenForce{Force: true})
		}
	case error:
		fmt.Printf("error message: %s (%s)\n", msg, ctx.Self().GetId())
	default:
		fmt.Printf("unhandled message type: %T (%s)\n", msg, ctx.Self().GetId())
	}
}

type tickResetCountersMsg struct{}
type tickMsg struct{}

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
		// if tRefg.Sub(time.Time{}) <= 24*time.Hour {
		t := time.Date(tn.Year(), tn.Month(), tn.Day(), 23, 59, 59, 0, tn.Location())
		until = time.Until(t)
		// } else {
		// 	until = time.Until(tRefg)
		// }
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
