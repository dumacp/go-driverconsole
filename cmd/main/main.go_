package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"
	"github.com/dumacp/go-driverconsole/internal/app"
	"github.com/dumacp/go-driverconsole/internal/buttons"
	"github.com/dumacp/go-driverconsole/internal/counterpass"
	"github.com/dumacp/go-driverconsole/internal/device"
	"github.com/dumacp/go-driverconsole/internal/display"
	"github.com/dumacp/go-driverconsole/internal/pubsub"
	"github.com/dumacp/go-logs/pkg/logs"
)

var port string
var baud int
var standalone bool
var showVersion bool

const version = "1.0.2"

func init() {
	flag.StringVar(&port, "port", "/dev/ttyUSB0", "path to port serial in OS")
	flag.IntVar(&baud, "baud", 19200, "serial port speed in baudios")
	flag.BoolVar(&standalone, "standalone", false, "standalone running (without appfare supervision)")
	flag.BoolVar(&showVersion, "version", false, "show version")
}

func main() {

	flag.Parse()
	if showVersion {
		fmt.Printf("version: %s\n", version)
		os.Exit(2)
	}

	sys := actor.NewActorSystem()

	var pidApp *actor.PID
	props := actor.PropsFromFunc(func(ctx actor.Context) {

		switch ctx.Message().(type) {
		case *actor.Started:

			pubsub.Init(ctx.ActorSystem().Root)

			propsDevice := actor.PropsFromFunc(device.NewActor(port, baud, 3*time.Second).Receive)

			propsDisplay := actor.PropsFromFunc(display.NewActor().Receive)

			propsButtons := actor.PropsFromFunc(buttons.NewActor().Receive)

			propsCounter := actor.PropsFromFunc(counterpass.NewActor().Receive)

			propsApp := actor.PropsFromFunc(app.NewActor().Receive)

			pidDevice, err := ctx.SpawnNamed(propsDevice, "device")
			if err != nil {
				log.Fatalln(err)
			}

			pidDisplay, err := ctx.SpawnNamed(propsDisplay, "display")
			if err != nil {
				log.Fatalln(err)
			}

			pidButtons, err := ctx.SpawnNamed(propsButtons, "buttons")
			if err != nil {
				log.Fatalln(err)
			}

			pidCounter, err := ctx.SpawnNamed(propsCounter, "counter")
			if err != nil {
				log.Fatalln(err)
			}

			pidApp, err = ctx.SpawnNamed(propsApp, "app")
			if err != nil {
				log.Fatalln(err)
			}

			ctx.RequestWithCustomSender(pidDevice, &device.Subscribe{}, pidButtons)
			ctx.RequestWithCustomSender(pidDevice, &device.Subscribe{}, pidDisplay)
			ctx.RequestWithCustomSender(pidCounter, &counterpass.MsgSubscribe{}, pidApp)
			ctx.RequestWithCustomSender(pidButtons, &buttons.MsgSubscribe{}, pidApp)
			ctx.RequestWithCustomSender(pidApp, &app.MsgSubscribe{}, pidDisplay)

			// routes := map[int32]string{
			// 	10: "RUTA CARAJILLO",
			// 	20: "RUTA ORIENTAL",
			// 	30: "RUTA OCCIDENTAL",
			// 	40: "RUTA NORTE",
			// 	50: "RUTA SUR",
			// }

			// ctx.Send(pidApp, &app.MsgSetRoutes{Routes: routes})
		case *actor.Stopping:
			log.Print("stopping driver console")
		case *actor.Stopped:
			log.Print("stopped driver console")
		case *actor.Terminated:
			log.Print("terminated driver console")
		case *actor.Restarting:
			log.Print("restarting driver console")
		case *actor.Restart:
			log.Print("restart driver console")
		default:
			if pidApp != nil {
				ctx.RequestWithCustomSender(pidApp, ctx.Message(), ctx.Sender())
			}

		}
		//}).WithSupervisor(strategy)
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

	if !standalone {
		r.Register("driverconsole", props)
		r.Start()
		log.Printf("kinds: %v", r.GetKnownKinds())
	} else {
		var err error
		pidMain, err = sys.Root.SpawnNamed(props, " driverconsole")
		if err != nil {
			log.Fatalln(err)
		}
		r.Start()
	}

	finish := make(chan os.Signal, 1)
	signal.Notify(finish, syscall.SIGINT)
	signal.Notify(finish, syscall.SIGTERM)
	signal.Notify(finish, os.Interrupt)

	<-finish
	if standalone && pidMain != nil {
		sys.Root.Poison(pidMain)
	}
	log.Print("finish")
}
