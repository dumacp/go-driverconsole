package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"

	"github.com/dumacp/go-driverconsole/internal/buttons"
	"github.com/dumacp/go-driverconsole/internal/counterpass"
	"github.com/dumacp/go-driverconsole/internal/ignition"
	app "github.com/dumacp/go-driverconsole/internal/sibus"
	"github.com/dumacp/go-driverconsole/internal/ui"
	"github.com/dumacp/go-driverconsole/internal/utils"

	"github.com/dumacp/go-driverconsole/internal/device"
	"github.com/dumacp/go-driverconsole/internal/display"
	"github.com/dumacp/go-driverconsole/internal/pubsub"

	"github.com/dumacp/go-logs/pkg/logs"
)

var port string
var baud int
var id string
var debug bool
var logStd bool
var showversion bool

const version = "1.1.4_sibus"

func init() {
	flag.StringVar(&id, "id", "", "device ID")
	flag.StringVar(&port, "port", "/dev/ttyUSB0", "path to port serial in OS")
	flag.IntVar(&baud, "baud", 38400, "serial port speed in baudios")
	flag.BoolVar(&debug, "debug", false, "debug")
	flag.BoolVar(&logStd, "logStd", false, "send logs to stdout")
	flag.BoolVar(&showversion, "version", false, "show version")
}

func main() {

	flag.Parse()
	if showversion {
		fmt.Printf("version: %s\n", version)
		os.Exit(2)
	}

	initLogs(debug, logStd)

	if len(id) <= 0 {
		id = utils.Hostname()
	} else {
		utils.SetHostname(id)
	}

	sys := actor.NewActorSystem()
	root := sys.Root

	// decider := func(reason interface{}) actor.Directive {
	// 	fmt.Println("handling failure for child")
	// 	return actor.RestartDirective
	// }

	// strategy := actor.NewAllForOneStrategy(100, 30*time.Second, decider)

	pubsub.Init(root)

	type Init struct{}
	var lastIgnitionEvent *ignition.IgnitionEvent
	var pidApp *actor.PID
	isAppActive := false
	var pidCounter *actor.PID
	var uii ui.UI

	props := actor.PropsFromFunc(func(ctx actor.Context) {

		switch msg := ctx.Message().(type) {
		case *actor.Started:

			if err := func() error {
				propsIgnition := actor.PropsFromFunc(ignition.NewIgnition().Receive)
				pidIgnition, err := ctx.SpawnNamed(propsIgnition, "ignition")
				if err != nil {
					fmt.Printf("ignition main-actor error: %s\n", err)
					return nil
				}

				res, err := ctx.RequestFuture(pidIgnition, &ignition.MsgGetIgnitionEvent{}, 2*time.Second).Result()
				if err != nil {
					fmt.Printf("ignition main-actor error: %s\n", err)
					return nil
				}
				switch e := res.(type) {
				case *ignition.IgnitionEvent:
					lastIgnitionEvent = e
					if e.StateEvent == ignition.DOWN {
						if err := exec.Command("/bin/sh", "-c", "echo 0 > /sys/class/leds/enable-hmi/brightness").Run(); err != nil {
							logs.LogWarn.Printf("error with command hmi off: %s", err)
						}
						return fmt.Errorf("ignition in DOWN state")
					}
				}
				return nil
			}(); err != nil {
				logs.LogError.Printf("ignition main-actor error: %s", err)
				break
			}
			if err := exec.Command("/bin/sh", "-c", "echo 1 > /sys/class/leds/enable-hmi/brightness").Run(); err != nil {
				logs.LogWarn.Printf("error with command hmi on: %s", err)
			}
			ctx.Send(ctx.Self(), &Init{})
		case *ignition.IgnitionEvent:
			switch {
			case lastIgnitionEvent != nil && lastIgnitionEvent.StateEvent == ignition.DOWN && msg.StateEvent == ignition.UP:
				ctx.Send(ctx.Self(), &Init{})
			case msg.StateEvent == ignition.DOWN:
				if pidCounter != nil {
					ctx.PoisonFuture(pidCounter)
				}
				isAppActive = false
				if pidApp != nil {
					ctx.PoisonFuture(pidApp).Wait()
				}

				if uii != nil {
					uii.Shutdown()
				}
			}
			lastIgnitionEvent = msg
		case *Init:
			if pidCounter != nil {
				ctx.PoisonFuture(pidCounter)
			}
			isAppActive = false
			if pidApp != nil {
				ctx.PoisonFuture(pidApp).Wait()
			}
			if uii != nil {
				uii.Shutdown()
			}
			var err error
			pidCounter, err = ctx.SpawnNamed(actor.PropsFromFunc(counterpass.NewActor().Receive), "counter-actor")
			if err != nil {
				log.Fatalf("counter actor error: %s", err)
			}

			var confDev device.Device
			var confButtons buttons.ButtonDevice
			var confDisplay display.Display

			confDev = device.NewPiDevice(port, baud)

			confButtons = buttons.NewConfPiButtons(0, 30, []int{
				app.AddrAddBright, app.AddrEnterDriver, app.AddrEnterPaso, app.AddrEnterRuta,
				app.AddrScreenAlarms, app.AddrSelectPaso, app.AddrSubBright, app.AddrScreenMore,
				app.AddrScreenProgDriver, app.AddrScreenProgVeh, app.AddrScreenSwitch,
				app.AddrSwitchStep, app.AddrSendStep},
			)
			confDisplay = display.NewPiDisplay(app.Label2DisplayRegister)

			uii, err = ui.New(ctx,
				device.NewActor(confDev),
				display.NewDisplayActor(confDisplay))

			if err != nil {
				log.Fatalf("newDisplayActor error: %s", err)
			}

			time.Sleep(3 * time.Second)

			appinstance := app.NewApp(uii)
			propsApp := actor.PropsFromFunc(appinstance.Receive)
			pidApp, err = ctx.SpawnNamed(propsApp, "app")
			isAppActive = true
			if err != nil {
				log.Fatalf("app-actor error: %s", err)
			}

			if err := uii.InputHandler(buttons.NewActor(confButtons), app.ButtonsPi(appinstance)); err != nil {
				log.Fatalf("inputHandler error: %s", err)
			}

		case *actor.Stopping:
			logs.LogInfo.Print("stopping driver console")
		case *actor.Stopped:
			logs.LogInfo.Print("stopped driver console")
		case *actor.Terminated:
			logs.LogInfo.Printf("terminated %q", msg.GetWho())
		case *actor.Restarting:
			logs.LogInfo.Print("restarting driver console")
		case *actor.Restart:
			logs.LogInfo.Print("restart driver console")
		default:
			fmt.Printf("main message: %q --> %q, %T (%s)\n", func() string {
				if ctx.Sender() == nil {
					return ""
				} else {
					return ctx.Sender().GetId()
				}
			}(), ctx.Self().GetId(), ctx.Message(), ctx.Message())
			if isAppActive && pidApp != nil {
				ctx.RequestWithCustomSender(pidApp, ctx.Message(), ctx.Sender())
			}

		}
		// }).WithSupervisor(strategy)
	})

	portlocal := 8099
	for {

		socket := fmt.Sprintf("127.0.0.1:%d", portlocal)
		testConn, err := net.DialTimeout("tcp", socket, 1*time.Second)
		if err != nil {
			break
		}
		testConn.Close()
		logs.LogWarn.Printf("socket busy -> \"%s\"", socket)
		time.Sleep(1 * time.Second)
		portlocal++
	}

	// kind := remote.NewKind("driverconsole", props)
	rconfig := remote.Configure("127.0.0.1", portlocal)
	//	remote.NewKind("driverconsole", props))
	r := remote.NewRemote(sys, rconfig)

	var pidMain *actor.PID

	var err error
	pidMain, err = sys.Root.SpawnNamed(props, " driverconsole")
	if err != nil {
		log.Fatalln(err)
	}
	r.Start()

	// log.Printf("kind: %s", remote.NewKind("driverconsole", props).Kind)

	finish := make(chan os.Signal, 1)
	signal.Notify(finish, syscall.SIGINT)
	signal.Notify(finish, syscall.SIGTERM)
	signal.Notify(finish, os.Interrupt)

	tickStart := time.NewTicker(1 * time.Second)
	defer tickStart.Stop()
	timerStart := time.NewTimer(10 * time.Second)
	defer timerStart.Stop()
	func() {
		for {
			select {
			case <-tickStart.C:
				if pidApp != nil {
					tickStart.Stop()
					return
				}
			case <-timerStart.C:
				return
			}
		}
	}()
	// sys.Root.Send(pidApp, &messages.MsgRoute{RouteCode: 10})

	// receiveSimulateDriverPaso := actor.PropsFromFunc(func(ctx actor.Context) {
	// 	fmt.Printf("message: %q --> %q, %T\n", func() string {
	// 		if ctx.Sender() == nil {
	// 			return ""
	// 		} else {
	// 			return ctx.Sender().GetId()
	// 		}
	// 	}(), ctx.Self().GetId(), ctx.Message())
	// 	switch msg := ctx.Message().(type) {
	// 	case *messages.MsgDriverPaso:
	// 		value := msg.GetValue()
	// 		if ctx.Sender() != nil {
	// 			ctx.Respond(&messages.MsgResponseDriverPaso{
	// 				Value: value,
	// 			})
	// 		}
	// 	}
	// })

	// pidPaso, err := root.SpawnNamed(receiveSimulateDriverPaso, "receiveSimulateDriverPaso")
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// root.RequestWithCustomSender(pidApp, &messages.MsgSubscribeConsole{}, pidPaso)

	if pidApp != nil {
		// routes := map[int32]string{
		// 	10: "RUTA CARAJILLO",
		// 	20: "RUTA ORIENTAL",
		// 	30: "RUTA OCCIDENTAL",
		// 	40: "RUTA NORTE",
		// 	50: "RUTA SUR",
		// }

		// pidFare, _ := root.SpawnNamed(actor.PropsFromFunc(func(ctx actor.Context) {
		// 	switch ctx.Message().(type) {
		// 	case *messages.MsgDriverPaso:
		// 		if ctx.Sender() == nil {
		// 			break
		// 		}
		// 		ctx.Respond(&messages.MsgAppPaso{
		// 			Code:  messages.MsgAppPaso_CASH,
		// 			Value: 1,
		// 		})
		// 	}
		// }), "appFareTest")

		// root.RequestWithCustomSender(pidApp, &messages.MsgSubscribeConsole{}, pidFare)

		// root.Request(pidApp, &app.MsgSetRoutes{
		// 	Routes: routes,
		// })
		tick0 := time.Tick(5 * time.Second)
		// tick1 := time.Tick(30 * time.Second)
		tick2 := time.Tick(3 * time.Second)
		tick3 := time.Tick(300 * time.Second)
		tick4 := time.Tick(330 * time.Second)

		// toggle := false
		for {
			select {
			case <-tick0:
				// toggle = !toggle
				// if toggle {
				// 	root.Request(pidApp, &messages.MsgGpsErr{})
				// 	root.Request(pidApp, &messages.MsgGroundErr{})
				// } else {
				// 	root.Request(pidApp, &messages.MsgGpsOk{})
				// 	root.Request(pidApp, &messages.MsgGroundOk{})
				// }
				// root.Send(pidApp, &counterpass.CounterMap{Inputs0: 20, Outputs1: 21})
			// case <-tick1:
			// root.Send(pidApp, &messages.MsgAppPaso{Value: 1})
			case <-tick2:
				// root.Send(pidApp, &messages.MsgAddAlarm{Alarm: fmt.Sprintf("%s: notif (( %d ))", time.Now().Format("2006-01-02 15:04"), countAlarm)})
				// countAlarm++

				if isAppActive && pidApp != nil {
					root.Send(pidApp, &app.MsgUpdateTime{})
				}
				// root.Send(pidApp, &counterpass.CounterEvent{Inputs: 1, Outputs: 1})

			case <-tick3:

				// root.Send(pidApp, &messages.MsgAppPaso{
				// 	Value: 1,
				// 	Code:  messages.MsgAppPaso_ELECTRONIC,
				// })
			case <-tick4:
			// root.Send(pidApp, &messages.MsgAppError{
			// 	Error: "entrada invalida",
			// })
			case <-finish:
				// TODO:
				sys.Root.PoisonFuture(pidMain).Wait()
				time.Sleep(400 * time.Millisecond)
				log.Print("Finish")

				// root.Poison(pidButtons)
				// root.Poison(pidDevice)
				// time.Sleep(300 * time.Millisecond)
				log.Print("finish")
				return

			}
		}
	}
}
