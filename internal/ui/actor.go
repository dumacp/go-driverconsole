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
		ctx.Request(a.pidDisplay, &display.SwitchScreenMsg{
			Num: 0,
		})

	case *MainScreenMsg:
		ctx.Request(a.pidDisplay, &display.SwitchScreenMsg{
			Num: 0,
		})
	case *TextWarningMsg:
		ctx.Request(a.pidDisplay, &display.WriteTextMsg{
			Label: display.WARNING_TEXT,
			Text:  msg.Text,
		})

	case *TextConfirmationMsg:
		ctx.Request(a.pidDisplay, &display.WriteTextMsg{
			Label: display.CONFIRMATION_TEXT,
			Text:  msg.Text,
		})

	case *TextConfirmationPopupMsg:
		ctx.Request(a.pidDisplay, &display.PopupMsg{
			Label: display.POPUP_TEXT,
			Text:  msg.Text,
		})

	case *TextWarningPopupMsg:
		ctx.Request(a.pidDisplay, &display.PopupMsg{
			Label: display.POPUP_WARN_TEXT,
			Text:  msg.Text,
		})

	case *InputsMsg:
		ctx.Request(a.pidDisplay, &display.WriteNumberMsg{
			Label: display.INPUT_NUM,
			Num:   int64(msg.In),
		})

	case *OutputsMsg:
		ctx.Request(a.pidDisplay, &display.WriteNumberMsg{
			Label: display.OUTPUTS_TEXT,
			Num:   int64(msg.Out),
		})

	case *DeviationInputsMsg:
		ctx.Request(a.pidDisplay, &display.WriteNumberMsg{
			Label: display.DEVIATION_TEXT,
			Num:   int64(msg.Dev),
		})

	case *RouteMsg:
		ctx.Request(a.pidDisplay, &display.WriteTextMsg{
			Label: display.ROUTE_TEXT,
			Text:  msg.Route,
		})
	case *DriverMsg:
		ctx.Request(a.pidDisplay, &display.WriteTextMsg{
			Label: display.ROUTE_TEXT,
			Text:  []string{msg.Data},
		})
	case *BeepMsg:
		ctx.Request(a.pidDisplay, &display.BeepMsg{
			Repeat:  3,
			Timeout: 1 * time.Second,
		})

	case *DateMsg:
		ctx.Request(a.pidDisplay, &display.WriteTextMsg{
			Label: display.DATE_TEXT,
			Text:  []string{msg.Date.Format("2006/01/02 15:04:05")},
		})

	case *ScreenMsg:
		res, err := ctx.RequestFuture(a.pidDisplay, &display.SwitchScreenMsg{
			Num: msg.Num,
		}, 1*time.Second).Result()
		if err != nil {
			logs.LogWarn.Printf("Get ScreenMsg error: %s", err)
		}
		if v, ok := res.(*display.ScreenResponseMsg); ok {
			a.screen = v.Screen
		}

	case *GetScreenMsg:
		res, err := ctx.RequestFuture(a.pidDisplay, &display.ScreenMsg{}, 1*time.Second).Result()
		if err != nil {
			logs.LogWarn.Printf("Get ScreenMsg error: %s", err)
		}
		if v, ok := res.(*display.ScreenResponseMsg); ok {
			a.screen = v.Screen
		}
	case *KeyNumMsg:
		ctx.Request(a.pidDisplay, &display.KeyNumMsg{
			Prompt: msg.Prompt,
		})
	case *KeyboardMsg:
		ctx.Request(a.pidDisplay, &display.KeyboardMsg{
			Prompt: msg.Prompt,
		})
	case *DoorsMsg:

	case *GpsMsg:

	case *NetworkMsg:

	case *AddNotificationsMsg:

	case *ShowNotificationsMsg:

	case *ShowProgDriverMsg:

	case *ShowProgVehMsg:

	case *ShowStatsMsg:

	case *BrightnessMsg:

	default:
		ctx.Respond(fmt.Errorf("unhandled message type: %T", msg))
	}
}
