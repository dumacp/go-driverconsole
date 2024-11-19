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
	screen   int
}

func NewActor(dev, disp actor.Actor) actor.Actor {

	a := &ActorUI{}
	a.actorDisplay = disp
	a.actorDevice = dev
	a.screen = MAIN_SCREEN
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
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.InitMsg{}, 5*time.Second))
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: result})
		}
	case *MainScreenMsg:
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.SwitchScreenMsg{
			Num: int(MAIN_SCREEN),
		}, time.Second*5))
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: result})
		}
		if result == nil {
			a.screen = MAIN_SCREEN
		}
	case *VerifyDisplayMsg:
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.VerifyMsg{}, time.Second*1))
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: result})
		}
	case *TextWarningMsg:
		fmt.Printf("screen: %d\n", a.screen)
		if a.screen != WARN_SCREEN {
			break
		}
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.WriteTextMsg{
			Label: WARNING_TEXT,
			Text:  msg.Text,
		}, time.Second*1))
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: result})
		}
	case *TextConfirmationMsg:
		fmt.Printf("screen: %d\n", a.screen)
		if a.screen != INFO_SCREEN {
			break
		}
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.WriteTextMsg{
			Label: CONFIRMATION_TEXT,
			Text:  msg.Text,
		}, 1*time.Second))
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: result})
		}
	case *TextConfirmationPopupMsg:
		fmt.Printf("screen: %d\n", a.screen)
		// if a.screen != MAIN_SCREEN {
		// 	if ctx.Sender() != nil {
		// 		ctx.Respond(&AckMsg{Error: nil})
		// 	}
		// 	break
		// }
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.PopupMsg{
			Label: POPUP_TEXT,
			Text:  msg.Text,
		}, 1*time.Second))
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: result})
		}
	case *TextConfirmationPopupCloseMsg:
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.PopupCloseMsg{
			Label: POPUP_TEXT,
		}, 1*time.Second))
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: result})
		}
	case *TextWarningPopupMsg:
		fmt.Printf("screen: %d\n", a.screen)
		// if a.screen != MAIN_SCREEN {
		// 	if ctx.Sender() != nil {
		// 		ctx.Respond(&AckMsg{Error: nil})
		// 	}
		// 	break
		// }
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.PopupMsg{
			Label: POPUP_WARN_TEXT,
			Text:  msg.Text,
		}, 1*time.Second))
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: result})
		}
	case *TextWarningPopupCloseMsg:
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.PopupCloseMsg{
			Label: POPUP_WARN_TEXT,
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
	case *CashInputsMsg:
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.WriteNumberMsg{
			Label: CASH_INPUTS_TEXT,
			Num:   int64(msg.In),
		}, 1*time.Second))
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: result})
		}
	case *ElectronicInputsMsg:
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.WriteNumberMsg{
			Label: ELECT_INPUTS_TEXT,
			Num:   int64(msg.In),
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
			Repeat: msg.Repeat,
			Period: msg.Period,
			Duty:   msg.Duty,
		}, 1*time.Second))
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: result})
		}
	case *DateMsg:
		fmt.Printf("screen: %d\n", a.screen)
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.WriteTextMsg{
			Label: DATE_TEXT,
			Text: func() []string {
				if len(msg.Format) <= 0 {
					return []string{msg.Date.Format("2006/01/02 15:04:05")}
				}
				return []string{msg.Date.Format(msg.Format)}
			}(),
		}, 1*time.Second))
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: result})
		}
	case *ScreenMsg:
		if msg.Force {
			res, err := ctx.RequestFuture(a.pidDisplay, &display.SwitchScreenMsg{
				Num: msg.Num,
			}, 5*time.Second).Result()
			if err != nil {
				if ctx.Sender() != nil {
					ctx.Respond(&AckMsg{Error: err})
				}
				logs.LogWarn.Printf("screenMsg error: %s", err)
				break
			}
			if v, ok := res.(*display.AckMsg); ok {
				if ctx.Sender() != nil {
					ctx.Respond(&AckMsg{Error: v.Error})
				}
				if v.Error == nil {
					a.screen = msg.Num
				}
			} else {
				if ctx.Sender() != nil {
					ctx.Respond(&AckMsg{Error: fmt.Errorf("unkown error")})
				}
			}
		} else {
			a.screen = msg.Num
			if ctx.Sender() != nil {
				ctx.Respond(&AckMsg{Error: nil})
			}
		}
		fmt.Printf("***** SCREEN SWITCH: %d\n", a.screen)
	case *GetScreenMsg:
		if ctx.Sender() == nil {
			break
		}
		ctx.Respond(&ScreenResponseMsg{Error: nil, Num: a.screen})
		// res, err := ctx.RequestFuture(a.pidDisplay, &display.ScreenMsg{}, 1*time.Second).Result()
		// if err != nil {
		// 	logs.LogWarn.Printf("getScreenMsg error: %s", err)
		// 	break
		// }
		// if v, ok := res.(*display.ScreenResponseMsg); ok {
		// 	if v.Error == nil {
		// 		a.screen = v.Screen
		// 	}
		// 	ctx.Respond(&ScreenResponseMsg{Error: v.Error, Num: v.Screen})
		// } else {
		// 	ctx.Respond(&ScreenResponseMsg{Error: fmt.Errorf("unkown error")})
		// }
	case *KeyNumMsg:
		if ctx.Sender() == nil {
			break
		}
		res, err := ctx.RequestFuture(a.pidDisplay, &display.KeyNumMsg{Prompt: msg.Prompt}, 1*time.Second).Result()
		if err != nil {
			logs.LogWarn.Printf("keyNumMsg error: %s", err)
			break
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
			break
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
	case *LedMsg:
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.LedMsg{
			State: func() int {
				if msg.State {
					return 0
				}
				return 1
			}(),
			Label: msg.Label,
		}, 1*time.Second))
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: result})
		}

	case *AddNotificationsMsg:

	case *ShowNotificationsMsg:
		if a.screen != ALARMS_SCREEN {
			break
		}
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.WriteTextMsg{
			Text:  msg.Text,
			Label: NOTIFICATIONS_ALARM_TEXT,
		}, time.Second*1))
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: result})
		}

	case *ShowProgDriverMsg:
		if a.screen != PROGRAMATION_DRIVER_SCREEN {
			break
		}
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.WriteTextMsg{
			Text:  msg.Text,
			Label: PROGRAMATION_DRIVER_TEXT,
		}, time.Second*1))
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: result})
		}
	case *ShowProgVehMsg:
		fmt.Printf("screen: %d\n", a.screen)
		if a.screen != PROGRAMATION_VEH_SCREEN {
			if ctx.Sender() != nil {
				ctx.Respond(&AckMsg{Error: nil})
			}
			break
		}
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
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.BrightnessMsg{Percent: msg.Percent}, 2*time.Second))
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: result})
		}
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
			logs.LogWarn.Printf("readBytesRawMsg error: %s", err)
			break
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
			ctx.Respond(&ReadBytesRawResponseMsg{Label: msg.Label, Error: fmt.Errorf("unkown error")})
		}
	case *WriteTextRawMsg:
		if ctx.Sender() == nil {
			break
		}
		res, err := ctx.RequestFuture(a.pidDisplay, &display.WriteTextMsg{Label: msg.Label, Text: msg.Text}, 1*time.Second).Result()
		if err != nil {
			logs.LogWarn.Printf("readBytesRawMsg error: %s", err)
			break
		}
		if v, ok := res.(*display.AckMsg); ok {
			if v.Error != nil {
				ctx.Respond(&AckMsg{
					Error: v.Error,
				})
				break
			}
			ctx.Respond(&AckMsg{Error: nil})
		} else {
			ctx.Respond(&AckMsg{Error: fmt.Errorf("unkown error (%T)", res)})
		}
	case *StepEnableMsg:
		result := AckResponse(ctx.RequestFuture(a.pidDisplay, &display.LedMsg{
			State: func() int {
				if msg.State {
					return 0
				}
				return 1
			}(),
			Label: STEP_ENABLE,
		}, 1*time.Second))
		if ctx.Sender() != nil {
			ctx.Respond(&AckMsg{Error: result})
		}
	case *buttons.StartButtons:

	case *buttons.InputEvent:

		fmt.Printf("arrive event: %+v\n", msg)
		// TODO: change this innencesary dependency
		if err := func() error {

			fmt.Printf("SCREEN: %v\n", a.screen)
			go a.evt2Func(msg)
			return nil
		}(); err != nil {
			logs.LogWarn.Println(err)
		}
	case error:
		fmt.Printf("error message: %s (%s)\n", msg, ctx.Self().GetId())
	default:
		fmt.Printf("unhandled message type: %T (%s)\n", msg, ctx.Self().GetId())
	}
}
