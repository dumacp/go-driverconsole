package display

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/dumacp/go-driverconsole/internal/device"
	"github.com/dumacp/go-logs/pkg/logs"
)

type DisplayActor struct {
	display Display
	// dev       device.Device
	pidDevice *actor.PID
	behavior  actor.Behavior
}

func NewDisplayActor(disp Display) actor.Actor {

	d := &DisplayActor{}
	d.display = disp
	d.behavior = actor.NewBehavior()
	d.behavior.Become(d.InitState)
	// d.dev = dev
	// d.actorButton = inputDevice
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
	switch ctx.Message().(type) {
	case *actor.Started:
	}
	d.behavior.Receive(ctx)

}

func (d *DisplayActor) InitState(ctx actor.Context) {

	switch msg := ctx.Message().(type) {
	case *actor.Started:
	case *device.MsgDevice:

		err := d.display.Init(msg.Device)
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: err})
		}
		logs.LogInfo.Printf("actor to runState")
		d.behavior.Become(d.RunState)
	case *InitMsg:
	default:
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: fmt.Errorf("actor in InitState")})
		}
	}
}

func (d *DisplayActor) RunState(ctx actor.Context) {

	switch msg := ctx.Message().(type) {
	case *actor.Started:
		// if d.dev != nil {
		// 	var err error
		// 	propsDevice := actor.PropsFromFunc(device.NewActor(d.dev).Receive)
		// 	d.pidDevice, err = ctx.SpawnNamed(propsDevice, "device-actor")
		// 	if err != nil {
		// 		time.Sleep(3 * time.Second)
		// 		logs.LogError.Panicf("create device-actor error: %s", err)
		// 	}
		// }
	case *device.MsgDevice:
		fmt.Println("new new")
		err := d.display.Init(msg.Device)
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: err})
		}
	case *InitMsg:

	case *CloseMsg:
		err := d.display.Close()
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: err})
		}
	case *SwitchScreenMsg:
		err := d.display.SwitchScreen(msg.Num)
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: err})
		}
	case *WriteTextMsg:
		err := d.display.WriteText(msg.Label, msg.Text...)
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: err})
		}
	case *ArrayPictMsg:
		fmt.Printf("arrayPictMsg: %v\n", msg)
		err := d.display.ArrayPict(msg.Label, msg.Num)
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: err})
		}
	case *WriteNumberMsg:
		err := d.display.WriteNumber(msg.Label, msg.Num)
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: err})
		}
	case *PopupMsg:
		err := d.display.Popup(msg.Label, msg.Text...)
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: err})
		}
	case *BeepMsg:
		err := d.display.Beep(msg.Repeat, msg.Timeout)
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: err})
		}
	case *VerifyMsg:
		err := d.display.Verify()
		// Manejar el error
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: err})
		}
	case *ScreenMsg:
		num, err := d.display.Screen()
		if ctx.Sender() != nil {
			ctx.Respond(&ScreenResponseMsg{
				Screen: num,
				Error:  err})
		}
	case *ResetMsg:
		err := d.display.Reset()
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: err})
		}
	case *LedMsg:
		err := d.display.Led(msg.Label, msg.State)
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: err})
		}
	case *KeyNumMsg:
		num, err := d.display.KeyNum(msg.Prompt)
		if ctx.Sender() != nil {
			ctx.Respond(&KeyNumResponseMsg{
				Num:   num,
				Error: err})
		}
	case *KeyboardMsg:
		text, err := d.display.Keyboard(msg.Prompt)
		if ctx.Sender() != nil {
			ctx.Respond(&KeyboardResponseMsg{
				Text:  text,
				Error: err})
		}
	case *BrightnessMsg:
		err := d.display.Brightness(msg.Percent)
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: err})
		}
	case *ReadBytesMsg:
		num, err := d.display.ReadBytes(msg.Label)
		if ctx.Sender() != nil {
			ctx.Respond(&ResponseBytesMsg{
				Label: msg.Label,
				Value: num,
				Error: err,
			})
		}

	// ... (Manejar los demás mensajes para los demás métodos)

	default:
		// Opcional: manejar mensajes desconocidos o inesperados
		logs.LogWarn.Printf("DisplayActor recibió un mensaje desconocido: %v", msg)
	}
}
