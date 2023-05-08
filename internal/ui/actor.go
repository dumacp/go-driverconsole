package ui

import (
	"fmt"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/dumacp/go-driverconsole/internal/display"
	"github.com/dumacp/go-logs/pkg/logs"
)

type ActorUI struct {
	ui           UI
	actorDisplay actor.Actor
	actorDevice  actor.Actor
	pidDisplay   *actor.PID
	pidDevice    *actor.PID
	pidInputs    *actor.PID
	screen       int
}

func NewActor(dev, disp actor.Actor) actor.Actor {

	a := &ActorUI{}
	a.actorDisplay = disp
	return a
}

func (a *ActorUI) Receive(ctx actor.Context) {
	fmt.Printf("message: %q --> %q, %T\n", func() string {
		if ctx.Sender() == nil {
			return ""
		} else {
			return ctx.Sender().GetId()
		}
	}(), ctx.Self().GetId(), ctx.Message())
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		propsDev := actor.PropsFromFunc(a.actorDevice.Receive)
		pidDev, err := ctx.SpawnNamed(propsDev, "display-actor")
		if err != nil {
			time.Sleep(3 * time.Second)
			logs.LogError.Panicf("%q error:", ctx.Self().GetId(), err)
		}
		propsDisplay := actor.PropsFromFunc(a.actorDisplay.Receive)
		pidDisplay, err := ctx.SpawnNamed(propsDisplay, "device-actor")
		if err != nil {
			time.Sleep(3 * time.Second)
			logs.LogError.Panicf("%q error:", ctx.Self().GetId(), err)
		}
		a.pidDevice = pidDev
		a.pidDisplay = pidDisplay
	case *display.DeviceMsg:
		if a.pidDisplay != nil {
			ctx.Request(a.pidDisplay, msg)
		}
		if a.pidInputs != nil {
			ctx.Request(a.pidInputs, msg)
		}
	case *InitUIMsg:
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.SwitchScreenMsg{
			Num: 0,
		}, 1*time.Second))
		if ctx.Sender() != nil {
			ctx.Respond(AckMsg{Error: result})
		}
	case *MainScreenMsg:
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.SwitchScreenMsg{
			Num: 0,
		}, time.Second*1))
		if ctx.Sender() != nil {
			ctx.Respond(AckMsg{Error: result})
		}
	case *TextWarningMsg:
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.WriteTextMsg{
			Label: WARNING_TEXT,
			Text:  msg.Text,
		}, time.Second*1))
		if ctx.Sender() != nil {
			ctx.Respond(AckMsg{Error: result})
		}
	case *TextConfirmationMsg:
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.WriteTextMsg{
			Label: CONFIRMATION_TEXT,
			Text:  msg.Text,
		}, 1*time.Second))
		if ctx.Sender() != nil {
			ctx.Respond(AckMsg{Error: result})
		}
	case *TextConfirmationPopupMsg:
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.PopupMsg{
			Label: POPUP_TEXT,
			Text:  msg.Text,
		}, 1*time.Second))
		if ctx.Sender() != nil {
			ctx.Respond(AckMsg{Error: result})
		}
	case *TextWarningPopupMsg:
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.PopupMsg{
			Label: POPUP_WARN_TEXT,
			Text:  msg.Text,
		}, 1*time.Second))
		if ctx.Sender() != nil {
			ctx.Respond(AckMsg{Error: result})
		}
	case *InputsMsg:
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.WriteNumberMsg{
			Label: INPUTS_TEXT,
			Num:   int64(msg.In),
		}, 1*time.Second))
		if ctx.Sender() != nil {
			ctx.Respond(AckMsg{Error: result})
		}
	case *OutputsMsg:
		ctx.Request(a.pidDisplay, &display.WriteNumberMsg{
			Label: OUTPUTS_TEXT,
			Num:   int64(msg.Out),
		})

	case *DeviationInputsMsg:
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.WriteNumberMsg{
			Label: DEVIATION_TEXT,
			Num:   int64(msg.Dev),
		}, 1*time.Second))
		if ctx.Sender() != nil {
			ctx.Respond(AckMsg{Error: result})
		}
	case *RouteMsg:
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.WriteTextMsg{
			Label: ROUTE_TEXT,
			Text:  msg.Route,
		}, 1*time.Second))
		if ctx.Sender() != nil {
			ctx.Respond(AckMsg{Error: result})
		}
	case *DriverMsg:
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.WriteTextMsg{
			Label: ROUTE_TEXT,
			Text:  []string{msg.Data},
		}, 1*time.Second))
		if ctx.Sender() != nil {
			ctx.Respond(AckMsg{Error: result})
		}
	case *BeepMsg:
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.BeepMsg{
			Repeat:  3,
			Timeout: 1 * time.Second,
		}, 1*time.Second))
		if ctx.Sender() != nil {
			ctx.Respond(AckMsg{Error: result})
		}
	case *DateMsg:
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.WriteTextMsg{
			Label: DATE_TEXT,
			Text:  []string{msg.Date.Format("2006/01/02 15:04:05")},
		}, 1*time.Second))
		if ctx.Sender() != nil {
			ctx.Respond(AckMsg{Error: result})
		}
	case *ScreenMsg:
		res, err := ctx.RequestFuture(a.pidDisplay, &display.SwitchScreenMsg{
			Num: msg.Num,
		}, 1*time.Second).Result()
		if err != nil {
			if ctx.Sender() != nil {
				ctx.Respond(&AckMsg{Error: err})
			}
			logs.LogWarn.Printf("Get ScreenMsg error: %s", err)
			break
		}
		if v, ok := res.(*display.ScreenResponseMsg); ok {
			a.screen = v.Screen
			if ctx.Sender() != nil {
				ctx.Respond(&AckMsg{Error: v.Error})
			}
		} else {
			if ctx.Sender() != nil {
				ctx.Respond(&AckMsg{Error: fmt.Errorf("unkown error")})
			}
		}
	case *GetScreenMsg:
		if ctx.Sender() == nil {
			break
		}
		res, err := ctx.RequestFuture(a.pidDisplay, &display.ScreenMsg{}, 1*time.Second).Result()
		if err != nil {
			logs.LogWarn.Printf("Get ScreenMsg error: %s", err)
		}
		if v, ok := res.(*display.ScreenResponseMsg); ok {
			if v.Error != nil {
				a.screen = v.Screen
			}
			if ctx.Sender() != nil {
				ctx.Respond(&ScreenResponseMsg{Error: v.Error, Num: v.Screen})
			}
		} else {
			if ctx.Sender() != nil {
				ctx.Respond(&ScreenResponseMsg{Error: fmt.Errorf("unkown error")})
			}
		}
	case *KeyNumMsg:
		if ctx.Sender() == nil {
			break
		}
		res, err := ctx.RequestFuture(a.pidDisplay, &display.KeyNumMsg{Prompt: msg.Prompt}, 1*time.Second).Result()
		if err != nil {
			logs.LogWarn.Printf("Get ScreenMsg error: %s", err)
		}
		if v, ok := res.(*display.KeyNumResponseMsg); ok {
			if ctx.Sender() != nil {
				ctx.Respond(&KeyNumResponseMsg{Error: v.Error, Num: v.Num})
			}
		} else {
			if ctx.Sender() != nil {
				ctx.Respond(&KeyNumResponseMsg{Error: fmt.Errorf("unkown error")})
			}
		}
	case *KeyboardMsg:
		if ctx.Sender() == nil {
			break
		}
		res, err := ctx.RequestFuture(a.pidDisplay, &display.KeyboardMsg{Prompt: msg.Prompt}, 1*time.Second).Result()
		if err != nil {
			logs.LogWarn.Printf("Get ScreenMsg error: %s", err)
		}
		if v, ok := res.(*display.KeyboardResponseMsg); ok {
			if ctx.Sender() != nil {
				ctx.Respond(&KeyboarResponsedMsg{Error: v.Error, Text: v.Text})
			}
		} else {
			if ctx.Sender() != nil {
				ctx.Respond(&KeyboarResponsedMsg{Error: fmt.Errorf("unkown error")})
			}
		}
	case *DoorsMsg:

	case *GpsMsg:

	case *NetworkMsg:

	case *AddNotificationsMsg:

	case *ShowNotificationsMsg:

	case *ShowProgDriverMsg:
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.SwitchScreenMsg{
			Num: 3,
		}, time.Second*1))
		if ctx.Sender() != nil {
			ctx.Respond(AckMsg{Error: result})
		}
		result = AckResponse(ctx.RequestFuture(a.pidDisplay, &display.WriteTextMsg{
			Text:  msg.Text,
			Label: PROGRAMATION_DRIVER_TEXT,
		}, time.Second*1))
		if ctx.Sender() != nil {
			ctx.Respond(AckMsg{Error: result})
		}
	case *ShowProgVehMsg:
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.SwitchScreenMsg{
			Num: 4,
		}, time.Second*1))
		if ctx.Sender() != nil {
			ctx.Respond(AckMsg{Error: result})
		}
		result = AckResponse(ctx.RequestFuture(a.pidDisplay, &display.WriteTextMsg{
			Text:  msg.Text,
			Label: PROGRAMATION_VEH_TEXT,
		}, time.Second*1))
		if ctx.Sender() != nil {
			ctx.Respond(AckMsg{Error: result})
		}

	case *ShowStatsMsg:

	case *BrightnessMsg:
	case *ServiceCurrentStateMsg:
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.WriteTextMsg{
			Label: SERVICE_CURRENT_STATE,
			Text:  []string{msg.Prompt},
		}, 1*time.Second))
		if ctx.Sender() != nil {
			ctx.Respond(AckMsg{Error: result})
		}
	case *AddInputsHandlerMsg:
		propsDev := actor.PropsFromFunc(msg.handler.Receive)
		pidDev, err := ctx.SpawnNamed(propsDev, "inputs-actor")
		if err != nil {
			time.Sleep(3 * time.Second)
			logs.LogError.Printf("%q error:", ctx.Self().GetId(), err)
		}
		a.pidInputs = pidDev

	default:
		ctx.Respond(fmt.Errorf("unhandled message type: %T", msg))
	}
}
