package display

import (
	"fmt"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/dumacp/go-driverconsole/internal/device"
	"github.com/dumacp/go-logs/pkg/logs"
)

type DisplayActor struct {
	display     Display
	propsDevice *actor.Props
	pidDevice   *actor.PID
}

func NewDisplayActor(dev *actor.Props) *DisplayActor {

	d := &DisplayActor{}
	d.propsDevice = dev
	return d
}

func (d *DisplayActor) Receive(ctx actor.Context) {
	fmt.Printf("message: %q --> %q, %T\n", func() string {
		if ctx.Sender() == nil {
			return ""
		} else {
			return ctx.Sender().GetId()
		}
	}(), ctx.Self().GetId(), ctx.Message())
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		if d.propsDevice != nil {
			var err error
			d.pidDevice, err = ctx.SpawnNamed(d.propsDevice, "device-actor")
			if err != nil {
				time.Sleep(3 * time.Second)
				logs.LogError.Panicf("create device-actor error: %s", err)
			}
		}
	case *device.MsgDevice:
		
	case *InitMsg:
		err := d.display.Init()
		if err != nil {
			// Manejar el error, por ejemplo, enviando un mensaje de error al remitente.
			ctx.Respond(err)
		}

	case *CloseMsg:
		err := d.display.Close()
		if err != nil {
			// Manejar el error
			ctx.Respond(err)
		}

	case *SwitchScreenMsg:
		err := d.display.SwitchScreen(msg.Num)
		if err != nil {
			// Manejar el error
			ctx.Respond(err)
		}

	case *WriteTextMsg:
		err := d.display.WriteText(msg.Label, msg.Text...)
		if err != nil {
			// Manejar el error
			ctx.Respond(err)
		}

	// ... (Manejar los demás mensajes para los demás métodos)

	default:
		// Opcional: manejar mensajes desconocidos o inesperados
		logs.LogWarn.Printf("DisplayActor recibió un mensaje desconocido: %v", msg)
	}
}
