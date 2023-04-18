package display

import (
	"fmt"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/dumacp/go-driverconsole/internal/device"
	"github.com/dumacp/go-logs/pkg/logs"
)

type DisplayActor struct {
	display   Display
	dev       device.Device
	pidDevice *actor.PID
}

func NewDisplayActor(dev device.Device) *DisplayActor {

	d := &DisplayActor{}
	d.dev = dev
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
		if d.dev != nil {
			var err error
			propsDevice := actor.PropsFromFunc(device.NewActor(d.dev).Receive)
			d.pidDevice, err = ctx.SpawnNamed(propsDevice, "device-actor")
			if err != nil {
				time.Sleep(3 * time.Second)
				logs.LogError.Panicf("create device-actor error: %s", err)
			}
		}
	case *device.MsgDevice:
		dev, err := New(msg.Device)
		if err != nil {
			logs.LogWarn.Println(err)
			break
		}
		d.display = dev
	case *InitMsg:
		err := d.display.Init()
		ctx.Respond(&AckMsg{Error: err})
	case *CloseMsg:
		err := d.display.Close()
		ctx.Respond(&AckMsg{Error: err})
	case *SwitchScreenMsg:
		err := d.display.SwitchScreen(msg.Num)
		ctx.Respond(&AckMsg{Error: err})
	case *WriteTextMsg:
		err := d.display.WriteText(msg.Label, msg.Text...)
		ctx.Respond(&AckMsg{Error: err})
	case *WriteNumberMsg:
		err := d.display.WriteNumber(msg.Label, msg.Num)
		ctx.Respond(&AckMsg{Error: err})
	case *PopupMsg:
		err := d.display.Popup(msg.Label, msg.Text...)
		ctx.Respond(&AckMsg{Error: err})
	case *BeepMsg:
		err := d.display.Beep(msg.Repeat, msg.Timeout)
		ctx.Respond(&AckMsg{Error: err})
	case *VerifyMsg:
		err := d.display.Verify()
		// Manejar el error
		ctx.Respond(&AckMsg{Error: err})
	case *ScreenMsg:
		num, err := d.display.Screen()
		ctx.Respond(&ScreenResponseMsg{
			Screen: num,
			Error:  err})
	case *ResetMsg:
		err := d.display.Reset()
		ctx.Respond(&AckMsg{Error: err})
	case *LedMsg:
		err := d.display.Led(msg.Label, msg.State)
		ctx.Respond(&AckMsg{Error: err})
	case *KeyNumMsg:
		num, err := d.display.KeyNum(msg.Prompt)
		ctx.Respond(&KeyNumResponseMsg{
			Num:   num,
			Error: err})
	case *KeyboardMsg:
		text, err := d.display.Keyboard(msg.Prompt)

		ctx.Respond(&KeyboardResponseMsg{
			Text:  text,
			Error: err})
	case *BrightnessMsg:
		err := d.display.Brightness(msg.Percent)
		ctx.Respond(&AckMsg{Error: err})

	// ... (Manejar los demás mensajes para los demás métodos)

	default:
		// Opcional: manejar mensajes desconocidos o inesperados
		logs.LogWarn.Printf("DisplayActor recibió un mensaje desconocido: %v", msg)
	}
}
