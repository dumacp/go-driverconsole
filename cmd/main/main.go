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

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/dumacp/go-driverconsole/internal/app"
	"github.com/dumacp/go-driverconsole/internal/buttons"
	"github.com/dumacp/go-driverconsole/internal/device"
	"github.com/dumacp/go-driverconsole/internal/display"
	"github.com/dumacp/go-logs/pkg/logs"
)

var port string
var baud int

func init() {
	flag.StringVar(&port, "port", "/dev/ttyUSB0", "path to port serial in OS")
	flag.IntVar(&baud, "baud", 19200, "serial port speed in baudios")
}

func main() {

	flag.Parse()

	sys := actor.NewActorSystem()
	// root := sys.Root

	//	decider := func(reason interface{}) actor.Directive {
	//		fmt.Println("handling failure for child")
	//		return actor.StopDirective
	//	}
	//
	//	strategy := actor.NewAllForOneStrategy(100, 30*time.Second, decider)

	props := actor.PropsFromFunc(func(ctx actor.Context) {

		switch ctx.Message().(type) {
		case *actor.Started:

			propsDevice := actor.PropsFromFunc(device.NewActor(port, baud, 3*time.Second).Receive)

			propsDisplay := actor.PropsFromFunc(display.NewActor().Receive)

			propsButtons := actor.PropsFromFunc(buttons.NewActor().Receive)

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

			pidApp, err := ctx.SpawnNamed(propsApp, "app")
			if err != nil {
				log.Fatalln(err)
			}

			ctx.RequestWithCustomSender(pidDevice, &device.Subscribe{}, pidButtons)
			ctx.RequestWithCustomSender(pidDevice, &device.Subscribe{}, pidDisplay)
			ctx.RequestWithCustomSender(pidButtons, &buttons.MsgSubscribe{}, pidApp)
			ctx.RequestWithCustomSender(pidApp, &app.MsgSubscribe{}, pidDisplay)

			routes := map[int]string{
				10: "RUTA CARAJILLO",
				20: "RUTA ORIENTAL",
				30: "RUTA OCCIDENTAL",
				40: "RUTA NORTE",
				50: "RUTA SUR",
			}

			ctx.Send(pidApp, &app.MsgSetRoutes{Routes: routes})

		case *actor.Stopped:
			log.Print("finished driver console")
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
	r.Register("driverconsole", props)
	r.Start()

	log.Printf("kinds: %v", r.GetKnownKinds())

	// log.Printf("kind: %s", remote.NewKind("driverconsole", props).Kind)

	finish := make(chan os.Signal, 1)
	signal.Notify(finish, syscall.SIGINT)
	signal.Notify(finish, syscall.SIGTERM)
	signal.Notify(finish, os.Interrupt)

	// go func() {

	// 	tick1 := time.Tick(10 * time.Second)
	// 	tick2 := time.Tick(20 * time.Second)
	// 	tick3 := time.Tick(30 * time.Second)

	// 	valuePercent := 0
	// 	valueDoor := [][2]int{{0, 1}, {0, 0}, {1, 0}, {1, 1}}
	// 	idxValueDoor := 0

	// 	for {
	// 		select {
	// 		case <-tick1:
	// 			root.Send(pidApp, &app.MsgAppPaso{Value: 1})
	// 		case <-tick2:
	// 			if valuePercent > 100 {
	// 				valuePercent = 0
	// 			} else {
	// 				valuePercent += 5
	// 			}
	// 			root.Send(pidApp, &app.MsgAppPercentRecorrido{Data: valuePercent})
	// 			if idxValueDoor >= len(valueDoor) {
	// 				idxValueDoor = 0
	// 			}
	// 			root.Send(pidApp, &app.MsgDoors{Value: valueDoor[idxValueDoor]})

	// 			idxValueDoor++
	// 		case <-tick3:
	// 			// root.Send(pidApp, &app.MsgConfirmationText{
	// 			// 	Text: []byte(fmt.Sprintf("texto de prueba\nTIME: %s", time.Now().Format("2006/01/02 15:04:05"))),
	// 			// })
	// 			// go func() {
	// 			// 	time.Sleep(3 * time.Second)
	// 			// 	root.Send(pidApp, &app.MsgMainScreen{})
	// 			// }()
	// 		}
	// 	}

	// }()

	for range finish {
		// time.Sleep(300 * time.Millisecond)
		// root.Poison(pidApp)
		// root.Poison(pidButtons)
		// root.Poison(pidDevice)
		// time.Sleep(300 * time.Millisecond)
		log.Print("finish")
		return
	}
}
