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
	propsDisplay *actor.Props
	pidDisplay   *actor.PID
	screen       int
}

func NewActor(disp *actor.Props) actor.Actor {

	a := &ActorUI{}
	a.propsDisplay = disp
	return a
}

func (a *ActorUI) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
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
			Label: display.WARNING_TEXT,
			Text:  msg.Text,
		}, time.Second*1))
		if ctx.Sender() != nil {
			ctx.Respond(AckMsg{Error: result})
		}
	case *TextConfirmationMsg:
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.WriteTextMsg{
			Label: display.CONFIRMATION_TEXT,
			Text:  msg.Text,
		}, 1*time.Second))
		if ctx.Sender() != nil {
			ctx.Respond(AckMsg{Error: result})
		}
	case *TextConfirmationPopupMsg:
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.PopupMsg{
			Label: display.POPUP_TEXT,
			Text:  msg.Text,
		}, 1*time.Second))
		if ctx.Sender() != nil {
			ctx.Respond(AckMsg{Error: result})
		}
	case *TextWarningPopupMsg:
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.PopupMsg{
			Label: display.POPUP_WARN_TEXT,
			Text:  msg.Text,
		}, 1*time.Second))
		if ctx.Sender() != nil {
			ctx.Respond(AckMsg{Error: result})
		}
	case *InputsMsg:
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.WriteNumberMsg{
			Label: display.INPUT_NUM,
			Num:   int64(msg.In),
		}, 1*time.Second))
		if ctx.Sender() != nil {
			ctx.Respond(AckMsg{Error: result})
		}
	case *OutputsMsg:
		ctx.Request(a.pidDisplay, &display.WriteNumberMsg{
			Label: display.OUTPUTS_TEXT,
			Num:   int64(msg.Out),
		})

	case *DeviationInputsMsg:
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.WriteNumberMsg{
			Label: display.DEVIATION_TEXT,
			Num:   int64(msg.Dev),
		}, 1*time.Second))
		if ctx.Sender() != nil {
			ctx.Respond(AckMsg{Error: result})
		}
	case *RouteMsg:
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.WriteTextMsg{
			Label: display.ROUTE_TEXT,
			Text:  msg.Route,
		}, 1*time.Second))
		if ctx.Sender() != nil {
			ctx.Respond(AckMsg{Error: result})
		}
	case *DriverMsg:
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.WriteTextMsg{
			Label: display.ROUTE_TEXT,
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
			Label: display.DATE_TEXT,
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
			Label: display.PROGRAMATION_DRIVER_TEXT,
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
			Label: display.PROGRAMATION_VEH_TEXT,
		}, time.Second*1))
		if ctx.Sender() != nil {
			ctx.Respond(AckMsg{Error: result})
		}

	case *ShowStatsMsg:

	case *BrightnessMsg:
	case *ServiceCurrentStateMsg:
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.WriteTextMsg{
			Label: display.SERVICE_CURRENT_STATE,
			Text:  []string{msg.Prompt},
		}, 1*time.Second))
		if ctx.Sender() != nil {
			ctx.Respond(AckMsg{Error: result})
		}

	default:
		ctx.Respond(fmt.Errorf("unhandled message type: %T", msg))
	}
}
