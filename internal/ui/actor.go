package ui

import (
	"fmt"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/dumacp/go-driverconsole/internal/display"
	"github.com/dumacp/go-logs/pkg/logs"
)

type ActorUI struct {
	ui         UI
	pidDisplay *actor.PID
	screen     int
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
		res, err := ctx.RequestFuture(a.pidDisplay, &display.ScreenMsg{}, 1*time.Second).Result()
		if err != nil {
			logs.LogWarn.Printf("Get ScreenMsg error: %s", err)
		}
		if v, ok := res.(*display.ScreenResponseMsg); ok {
			a.screen = v.Screen
		}

	case *GetScreenMsg:
		screen, err := a.ui.GetScreen()
		if err != nil {
			ctx.Respond(err)
		} else {
			ctx.Respond(screen)
		}
	case *KeyNumMsg:
		num, err := a.ui.KeyNum(msg.Prompt)
		if err != nil {
			ctx.Respond(err)
		} else {
			ctx.Respond(num)
		}
	case *KeyboardMsg:
		text, err := a.ui.Keyboard(msg.Prompt)
		if err != nil {
			ctx.Respond(err)
		} else {
			ctx.Respond(text)
		}
	case *DoorsMsg:
		err := a.ui.Doors(msg.State...)
		if err != nil {
			ctx.Respond(err)
		}
	case *GpsMsg:
		err := a.ui.Gps(msg.State)
		if err != nil {
			ctx.Respond(err)
		}
	case *NetworkMsg:
		err := a.ui.Network(msg.State)
		if err != nil {
			ctx.Respond(err)
		}
	case *AddNotificationsMsg:
		err := a.ui.AddNotifications(msg.Add)
		if err != nil {
			ctx.Respond(err)
		}
	case *ShowNotificationsMsg:
		err := a.ui.ShowNotifications()
		if err != nil {
			ctx.Respond(err)
		}
	case *ShowProgDriverMsg:
		err := a.ui.ShowProgDriver()
		if err != nil {
			ctx.Respond(err)
		}
	case *ShowProgVehMsg:
		err := a.ui.ShowProgVeh()
		if err != nil {
			ctx.Respond(err)
		}

	case *ShowStatsMsg:
		err := a.ui.ShowStats()
		if err != nil {
			ctx.Respond(err)
		}
	case *BrightnessMsg:
		err := a.ui.Brightness(msg.Percent)
		if err != nil {
			ctx.Respond(err)
		}
	default:
		ctx.Respond(fmt.Errorf("unhandled message type: %T", msg))
	}
}
