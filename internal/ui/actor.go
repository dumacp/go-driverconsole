package ui

import (
	"fmt"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/dumacp/go-driverconsole/internal/buttons"
	"github.com/dumacp/go-driverconsole/internal/device"
	"github.com/dumacp/go-driverconsole/internal/display"
	"github.com/dumacp/go-logs/pkg/logs"
)

type ActorUI struct {
	actorDisplay actor.Actor
	actorDevice  actor.Actor
	pidDisplay   *actor.PID
	pidDevice    *actor.PID
	pidInputs    *actor.PID
	dev          interface{}
	// evt2Label    map[int]EventType
	evt2Func func(evt *buttons.InputEvent)
	screen   Screen
}

func NewActor(dev, disp actor.Actor) actor.Actor {

	a := &ActorUI{}
	a.actorDisplay = disp
	a.actorDevice = dev
	a.screen = MAIN
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
		pidDev, err := ctx.SpawnNamed(propsDev, "device-actor")
		if err != nil {
			time.Sleep(3 * time.Second)
			logs.LogError.Panicf("%q error: %s", ctx.Self().GetId(), err)
		}
		propsDisplay := actor.PropsFromFunc(a.actorDisplay.Receive, actor.WithMailbox(actor.Bounded(10)))
		pidDisplay, err := ctx.SpawnNamed(propsDisplay, "display-actor")
		if err != nil {
			time.Sleep(3 * time.Second)
			logs.LogError.Panicf("%q error: %s", ctx.Self().GetId(), err)
		}
		a.pidDevice = pidDev
		a.pidDisplay = pidDisplay
	case *device.MsgDevice:
		a.dev = msg.Device
		if a.pidDisplay != nil {
			ctx.Send(a.pidDisplay, msg)
		}
		if a.pidInputs != nil {
			ctx.Send(a.pidInputs, msg)
		}
		// TODO: replace this slepp (i need wait because init internal devices)
		time.Sleep(1 * time.Second)
	case *InitUIMsg:
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.InitMsg{}, 1*time.Second))
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: result})
		}
	case *MainScreenMsg:
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.SwitchScreenMsg{
			Num: 0,
		}, time.Second*1))
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: result})
		}
		if result == nil {
			a.screen = 0
		}
	case *TextWarningMsg:
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.WriteTextMsg{
			Label: WARNING_TEXT,
			Text:  msg.Text,
		}, time.Second*1))
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: result})
		}
	case *TextConfirmationMsg:
		fmt.Printf("screen: %d\n", a.screen)
		if a.screen != CONFIRMATION {
			result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.SwitchScreenMsg{
				Num: int(CONFIRMATION),
			}, 1*time.Second))
			if ctx.Sender() != nil {
				ctx.Respond(&AckMsg{Error: result})
			}
			a.screen = CONFIRMATION
		}
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.WriteTextMsg{
			Label: CONFIRMATION_TEXT,
			Text:  msg.Text,
		}, 1*time.Second))
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: result})
		}
	case *TextConfirmationPopupMsg:
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.PopupMsg{
			Label:  POPUP_TEXT,
			Text:   msg.Text,
			Temout: msg.Timeout,
		}, 1*time.Second))
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: result})
		}
	case *TextWarningPopupMsg:
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.PopupMsg{
			Label:  POPUP_WARN_TEXT,
			Text:   msg.Text,
			Temout: msg.Timeout,
		}, 1*time.Second))
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: result})
		}
	case *InputsMsg:
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.WriteNumberMsg{
			Label: INPUTS_TEXT,
			Num:   int64(msg.In),
		}, 1*time.Second))
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: result})
		}
	case *OutputsMsg:
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.WriteNumberMsg{
			Label: OUTPUTS_TEXT,
			Num:   int64(msg.Out),
		}, 1*time.Second))
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: result})
		}

	case *DeviationInputsMsg:
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.WriteNumberMsg{
			Label: DEVIATION_TEXT,
			Num:   int64(msg.Dev),
		}, 1*time.Second))
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: result})
		}
	case *RouteMsg:
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.WriteTextMsg{
			Label: ROUTE_TEXT,
			Text:  msg.Route,
		}, 1*time.Second))
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: result})
		}
	case *DriverMsg:
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.WriteTextMsg{
			Label: DRIVER_TEXT,
			Text:  []string{msg.Data},
		}, 1*time.Second))
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: result})
		}
	case *BeepMsg:
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.BeepMsg{
			Repeat:  3,
			Timeout: 1 * time.Second,
		}, 1*time.Second))
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: result})
		}
	case *DateMsg:
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.WriteTextMsg{
			Label: DATE_TEXT,
			Text:  []string{msg.Date.Format("2006/01/02 15:04:05")},
		}, 1*time.Second))
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: result})
		}
	case *ScreenMsg:
		res, err := ctx.RequestFuture(a.pidDisplay, &display.SwitchScreenMsg{
			Num: msg.Num,
		}, 1*time.Second).Result()
		if err != nil {
			if ctx.Sender() != nil {
				ctx.Respond(&AckMsg{Error: err})
			}
			logs.LogWarn.Printf("get ScreenMsg error: %s", err)
			break
		}
		if v, ok := res.(*display.AckMsg); ok {
			if v.Error != nil && ctx.Sender() != nil {
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
			if v.Error == nil {
				a.screen = Screen(v.Screen)
			}
			ctx.Respond(&ScreenResponseMsg{Error: v.Error, Num: v.Screen})
		} else {
			ctx.Respond(&ScreenResponseMsg{Error: fmt.Errorf("unkown error")})
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
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.LedMsg{
			State: func() int {
				if msg.State {
					return 0
				}
				return 1
			}(),
			Label: GPS_LED,
		}, 1*time.Second))
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: result})
		}

	case *NetworkMsg:
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.LedMsg{
			State: func() int {
				if msg.State {
					return 0
				}
				return 1
			}(),
			Label: NETWORK_LED,
		}, 1*time.Second))
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: result})
		}

	case *AddNotificationsMsg:

	case *ShowNotificationsMsg:
		if a.screen != ALARM {
			break
		}
		// result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.SwitchScreenMsg{
		// 	Num: 3,
		// }, time.Second*1))
		// if ctx.Sender() != nil {
		// 	ctx.Respond(&AckMsg{Error: result})
		// }
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.WriteTextMsg{
			Text:  msg.Text,
			Label: NOTIFICATIONS_ALARM_TEXT,
		}, time.Second*1))
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: result})
		}

	case *ShowProgDriverMsg:
		if a.screen != DRIVER_SCREEN {
			break
		}
		// result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.SwitchScreenMsg{
		// 	Num: 3,
		// }, time.Second*1))
		// if ctx.Sender() != nil {
		// 	ctx.Respond(&AckMsg{Error: result})
		// }
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.WriteTextMsg{
			Text:  msg.Text,
			Label: PROGRAMATION_DRIVER_TEXT,
		}, time.Second*1))
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: result})
		}
	case *ShowProgVehMsg:
		if a.screen != VEHICLE {
			if ctx.Sender() != nil {
				ctx.Respond(&AckMsg{Error: nil})
			}
			break
		}
		// result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.SwitchScreenMsg{
		// 	Num: 4,
		// }, time.Second*1))
		// if ctx.Sender() != nil {
		// 	ctx.Respond(&AckMsg{Error: result})
		// }
		if len(msg.Text) > 0 {
			result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.WriteTextMsg{
				Text:  msg.Text,
				Label: PROGRAMATION_VEH_TEXT,
			}, time.Second*3))
			if ctx.Sender() != nil {
				ctx.Respond(&AckMsg{Error: result})
			}
		}

	case *ShowStatsMsg:

	case *BrightnessMsg:
	case *ServiceCurrentStateMsg:
		// if a.screen != MAIN {
		// 	break
		// }
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.WriteTextMsg{
			Label: SERVICE_CURRENT_STATE_TEXT,
			Text:  []string{msg.Prompt},
		}, 1*time.Second))
		if ctx.Sender() != nil && result != nil {
			ctx.Respond(&AckMsg{Error: result})
			break
		}
		fmt.Println("/// led ///")
		result = AckResponse(ctx.RequestFuture(a.pidDisplay, &display.ArrayPictMsg{
			Label: SERVICE_CURRENT_STATE,
			Num:   msg.State,
		}, 1*time.Second))
		fmt.Println("/// led ///")
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: result})
			break
		}
		fmt.Println("/// led ///")
	case *AddInputsHandlerMsg:
		propsDev := actor.PropsFromFunc(msg.Handler.Receive)
		pidDev, err := ctx.SpawnNamed(propsDev, "inputs-actor")
		if err != nil {
			// time.Sleep(3 * time.Second)
			logs.LogError.Printf("%q error: %s", ctx.Self().GetId(), err)
			if ctx.Sender() != nil {
				ctx.Respond(&AckMsg{Error: err})
			}
			break
		}
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: err})
		}
		a.evt2Func = msg.Evt2Func
		a.pidInputs = pidDev
		if a.dev != nil {
			ctx.Send(a.pidInputs, &device.MsgDevice{
				Device: a.dev,
			})
			// TODO: replace this slepp (i need wait because init internal devices)
			time.Sleep(1 * time.Second)
		}

	case *ReadBytesRawMsg:
		if ctx.Sender() == nil {
			break
		}
		res, err := ctx.RequestFuture(a.pidDisplay, &display.ReadBytesMsg{Label: msg.Label}, 1*time.Second).Result()
		if err != nil {
			logs.LogWarn.Printf("Get ScreenMsg error: %s", err)
		}
		if v, ok := res.(*display.ResponseBytesMsg); ok {
			if v.Error != nil {
				ctx.Respond(&ReadBytesRawResponseMsg{
					Label: v.Label,
					Value: nil,
					Error: v.Error,
				})
				break
			}
			ctx.Respond(&ReadBytesRawResponseMsg{Label: v.Label, Value: v.Value})
		} else {
			ctx.Respond(&ReadBytesRawResponseMsg{Label: v.Label, Error: fmt.Errorf("unkown error")})
		}

	case *buttons.InputEvent:

		fmt.Printf("arrive event: %+v\n", msg)
		// TODO: change this innencesary dependency
		if err := func() error {
			switch msg.KeyCode {
			case buttons.AddrScreenSwitch:
				a.screen = MAIN
			case buttons.AddrScreenProgDriver:
				a.screen = DRIVER_SCREEN
				// ctx.Send(ctx.Self(), &ShowProgDriverMsg{})
			case buttons.AddrScreenProgVeh:
				a.screen = VEHICLE
				// ctx.Send(ctx.Self(), &ShowProgVehMsg{})
			case buttons.AddrScreenAlarms:
				a.screen = ALARM
				// ctx.Send(ctx.Self(), &ShowNotificationsMsg{})
			case buttons.AddrScreenMore:
				a.screen = SERVICE
				// ctx.Send(ctx.Self(), &ShowStatsMsg{})
			case buttons.AddrEnterRuta:
				result, err := ctx.RequestFuture(a.pidDisplay, &display.ReadBytesMsg{
					Label: ROUTE_TEXT_READ,
				}, 1*time.Second).Result()
				if err != nil {
					return fmt.Errorf("get route error: %s", err)
				}
				fmt.Printf("route: %s\n", result)
				// switch v := result.(type) {
				// case *display.ResponseBytesMsg:
				// 	if v.Error != nil {
				// 		return v.Error
				// 	}
				// 	select {
				// 	case a.ui.chEvents <- &Event{
				// 		Type:  label,
				// 		Value: v.Value,
				// 	}:
				// 	default:
				// 	}
				// }
			case buttons.AddrEnterDriver:
				// result, err := ctx.RequestFuture(a.pidDisplay, &display.ReadBytesMsg{
				// 	Label: DRIVER_TEXT_READ,
				// }, 1*time.Second).Result()
				// if err != nil {
				// 	return fmt.Errorf("get driver error: %s", err)
				// }
				// fmt.Printf("driver: %s\n", result)

			}
			fmt.Printf("SCREEN: %v\n", a.screen)
			go a.evt2Func(msg)
			return nil
		}(); err != nil {
			logs.LogWarn.Println(err)
		}

	default:
		ctx.Respond(fmt.Errorf("unhandled message type: %T", msg))
	}
}
