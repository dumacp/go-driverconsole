package app

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/dumacp/go-driverconsole/internal/buttons"
	"github.com/dumacp/go-driverconsole/internal/ui"
	"github.com/dumacp/go-logs/pkg/logs"
	"github.com/dumacp/go-schservices/api/services"
)

func ButtonsGtt50(a *App) func(evt *buttons.InputEvent) {

	return func(evt *buttons.InputEvent) {

		// label, ok := evt2EvtLabel[int(evt.KeyCode)]
		// if !ok {
		// 	return
		// }
		if err := func() error {
			switchScreen := true
			switch evt.KeyCode {
			case AddrGttReturnNotiAlarm, AddrGttReturnProgDriver,
				AddrGttReturnProgVeh, AddrGttReturnWarning, AddrGttReturnConfimration,
				AddrGttReturnStats:
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				a.ctx.Send(a.ctx.Self(), &MsgMainScreen{})
			case AddrGttEnterScreenProgDriver:
				// release button
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				if err := a.uix.Screen(int(ui.PROGRAMATION_DRIVER_SCREEN), switchScreen); err != nil {
					return fmt.Errorf("event SCREEN error: %s", err)
				}

				ss := make([]*services.ScheduleService, 0)
				if len(a.shcservices) > 0 {
					for _, v := range a.shcservices {
						ss = append(ss, v)
					}
					sort.SliceStable(ss, func(i, j int) bool {
						return ss[i].GetScheduleDateTime() < ss[j].GetScheduleDateTime()
					})
				}
				reverseSlice := make([]string, 0)
				// reverseSlice = append(reverseSlice, "")
				for i := range ss {
					v := ss[len(ss)-1-i]
					if v.GetItinerary() == nil || v.GetDriver() == nil {
						continue
					}
					ts := time.UnixMilli(v.GetScheduleDateTime())

					capitalize := func(sstr string) string {
						result := make([]string, 0)
						for _, str := range strings.Fields(sstr) {
							runes := []rune(str)
							runes[0] = unicode.ToUpper(runes[0])
							result = append(result, string(runes))
						}
						return strings.Join(result, " ")
					}

					data := fmt.Sprintf(" %s: (%d) %s", ts.Format("01/02 15:04"),
						v.GetItinerary().GetId(), capitalize(strings.ToLower(v.GetDriver().GetFullName())))
					fmt.Printf("servicio: %v\n", v)
					fmt.Printf("data: %s\n", data)
					reverseSlice = append(reverseSlice, data)
					if time.Until(ts) < 0 {
						break
					}
				}
				dataSlice := make([]string, 0)
				for i := range reverseSlice {
					dataSlice = append(dataSlice, reverseSlice[len(reverseSlice)-i-1])
					if i >= 9 {
						break
					}
				}
				fmt.Printf("dataslice: %v\n", dataSlice)
				if err := a.uix.ShowProgDriver(dataSlice...); err != nil {
					return fmt.Errorf("event ShowProgDriver error: %s", err)
				}
			case AddrGttEnterScreenProgVeh:
				// release button
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				if err := a.uix.Screen(int(ui.PROGRAMATION_VEH_SCREEN), switchScreen); err != nil {
					return fmt.Errorf("event SCREEN error: %s", err)
				}

				ss := make([]*services.ScheduleService, 0)
				if len(a.shcservices) > 0 {
					for _, v := range a.shcservices {
						ss = append(ss, v)
					}
					sort.SliceStable(ss, func(i, j int) bool {
						return ss[i].GetScheduleDateTime() < ss[j].GetScheduleDateTime()
					})
				}
				reverseSlice := make([]string, 0)
				// reverseSlice = append(reverseSlice, "")
				for i := range ss {
					v := ss[len(ss)-1-i]
					if v.GetItinerary() == nil || v.GetRoute() == nil {
						continue
					}
					ts := time.UnixMilli(v.GetScheduleDateTime())

					data := strings.ToLower(fmt.Sprintf(" %s: (%d) %s (%s)", ts.Format("01/02 15:04"),
						v.GetItinerary().GetId(), v.GetItinerary().GetName(), v.GetRoute().GetName()))
					fmt.Printf("servicio: %v\n", v)
					fmt.Printf("data: %s\n", data)
					reverseSlice = append(reverseSlice, data)
					if time.Until(ts) < 0 {
						break
					}
				}
				dataSlice := make([]string, 0)
				for i := range reverseSlice {
					dataSlice = append(dataSlice, reverseSlice[len(reverseSlice)-i-1])
					if i >= 9 {
						break
					}
				}
				fmt.Printf("dataslice: %v\n", dataSlice)
				if err := a.uix.ShowProgVeh(dataSlice...); err != nil {
					return fmt.Errorf("event ShowProgVeh error: %s", err)
				}
			case AddrGttEnterScreenAlarms:
				// release button
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				if err := a.uix.Screen(int(ui.ALARMS_SCREEN), switchScreen); err != nil {
					return fmt.Errorf("event SCREEN error: %s", err)
				}

				if len(a.notif) > 0 {
					if err := a.uix.ShowNotifications(a.notif...); err != nil {
						return fmt.Errorf("event ShowNotifications error: %s", err)
					}
				}
			case AddrGttEnterScreenMore:
				if err := a.uix.Screen(int(ui.ADDITIONALS_SCREEN), switchScreen); err != nil {
					return fmt.Errorf("event SCREEN error: %s", err)
				}
				// release button
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				if err := a.uix.ShowStats(); err != nil {
					return fmt.Errorf("event ShowStats error: %s", err)
				}
			case AddrGttButtonRoute:
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				if err := a.uix.Screen(int(ui.KEY_ROUTE_SCREEN), switchScreen); err != nil {
					return fmt.Errorf("event SCREEN error: %s", err)
				}
			case AddrGttButtonRoute_SEND:
				// release button
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				fmt.Printf("routeCode: %s\n", a.routeCode)
				rutaCodeInt, _ := strconv.Atoi(a.routeCode)
				// a.ctx.Send(a.ctx.Self(), &MsgSetRoute{Route: rutaCodeInt})
				a.routeCode = ""
				a.route = rutaCodeInt
				if a.routes != nil {
					if routeString, ok := a.routes[int32(rutaCodeInt)]; ok {
						a.routeString = routeString
					} else {
						a.routeString = fmt.Sprintf("%d", rutaCodeInt)
					}
				}
				a.ctx.Send(a.ctx.Self(), &MsgSetRoute{
					Route: rutaCodeInt,
				})
				a.ctx.Send(a.ctx.Self(), &MsgMainScreen{})
			case AddrGttButtonDriver:
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				if err := a.uix.Screen(int(ui.KEY_DRIVER_SCREEN), switchScreen); err != nil {
					return fmt.Errorf("event SCREEN error: %s", err)
				}
			case AddrGttButtonDriver_SEND:
				// release button
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				fmt.Printf("driverCode: %s\n", a.driverCode)
				driverCodeInt, _ := strconv.Atoi(a.driverCode)
				fmt.Printf("driverCode: %d\n", driverCodeInt)
				a.driver = driverCodeInt
				// a.ctx.Send(a.ctx.Self(), &MsgSetDriver{Driver: driverCodeInt})
				a.driverCode = ""
				a.ctx.Send(a.ctx.Self(), &MsgSetDriver{
					Driver: driverCodeInt,
				})
				a.ctx.Send(a.ctx.Self(), &MsgMainScreen{})
			case AddrGttSwitchStep:
				if v, ok := evt.Value.(bool); ok {
					if v {
						fmt.Println("/////////// step enable: true ////////////")
						if a.enableStep {
							a.uix.StepEnable(false)
							a.enableStep = false
							if a.cancelStep != nil {
								a.cancelStep()
							}
							break
						}
						a.enableStep = true
						if a.cancelStep != nil {
							a.cancelStep()
						}
						a.uix.StepEnable(false)

						go func() {

							contxt, cancel := context.WithCancel(context.Background())
							defer cancel()
							a.cancelStep = cancel

							func() {
								for {
									timer := time.NewTimer(10 * time.Second)
									defer timer.Stop()

									renewcontxt, renewcancel := context.WithCancel(context.TODO())
									defer renewcancel()
									a.renewStep = renewcancel

									select {
									case <-contxt.Done():
										return
									case <-renewcontxt.Done():
										if !timer.Stop() {
											select {
											case <-timer.C:
											case <-time.After(30 * time.Millisecond):
											}
										}
										timer.Reset(10 * time.Second)
									case <-timer.C:
										return
									}

								}
							}()
							a.enableStep = false
							fmt.Println("/////////// step enable: false ////////////")
							a.uix.StepEnable(true)
						}()
						// } else {
						// 	a.uix.StepEnable(false)
						// 	a.enableStep = false
						// 	if a.cancelStep != nil {
						// 		a.cancelStep()
						// 	}
					}
				}
			case AddrGttSendStep:
				// release button
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				if a.enableStep {
					a.ctx.Send(a.ctx.Self(), &StepMsg{})
					if a.renewStep != nil {
						a.renewStep()
					}
				}
			case AddrGttAddBright:
				// release button
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				if a.brightness >= 90 {
					a.brightness = 100
				} else {
					a.brightness += 10
				}
				if err := a.uix.Brightness(a.brightness); err != nil {
					return fmt.Errorf("event SCREEN error: %s", err)
				}

			case AddrGttSubBright:
				// release button
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				if a.brightness <= 10 {
					a.brightness = 10
				} else {
					a.brightness -= 10
				}
				if err := a.uix.Brightness(a.brightness); err != nil {
					return fmt.Errorf("event SCREEN error: %s", err)
				}

			case AddrGttButtonRoute_0:
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				if len(a.routeCode) > 0 {
					a.routeCode = fmt.Sprintf("%s%d", a.routeCode, 0)
				}
			case AddrGttButtonRoute_1:
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				a.routeCode = fmt.Sprintf("%s%d", a.routeCode, 1)
			case AddrGttButtonRoute_2:
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				a.routeCode = fmt.Sprintf("%s%d", a.routeCode, 2)
			case AddrGttButtonRoute_3:
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				a.routeCode = fmt.Sprintf("%s%d", a.routeCode, 3)
			case AddrGttButtonRoute_4:
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				a.routeCode = fmt.Sprintf("%s%d", a.routeCode, 4)
			case AddrGttButtonRoute_5:
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				a.routeCode = fmt.Sprintf("%s%d", a.routeCode, 5)
			case AddrGttButtonRoute_6:
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				a.routeCode = fmt.Sprintf("%s%d", a.routeCode, 6)
			case AddrGttButtonRoute_7:
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				a.routeCode = fmt.Sprintf("%s%d", a.routeCode, 7)
			case AddrGttButtonRoute_8:
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				a.routeCode = fmt.Sprintf("%s%d", a.routeCode, 8)
			case AddrGttButtonRoute_9:
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				a.routeCode = fmt.Sprintf("%s%d", a.routeCode, 9)
			case AddrGttButtonRoute_Clear:
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				a.routeCode = ""
			case AddrGttButtonRoute_Cancel:
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				a.routeCode = ""
				a.ctx.Send(a.ctx.Self(), &MsgMainScreen{})
			case AddrGttButtonDriver_0:
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				if len(a.driverCode) > 0 {
					a.driverCode = fmt.Sprintf("%s%d", a.driverCode, 0)
				}
			case AddrGttButtonDriver_1:
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				a.driverCode = fmt.Sprintf("%s%d", a.driverCode, 1)
			case AddrGttButtonDriver_2:
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				a.driverCode = fmt.Sprintf("%s%d", a.driverCode, 2)
			case AddrGttButtonDriver_3:
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				a.driverCode = fmt.Sprintf("%s%d", a.driverCode, 3)
			case AddrGttButtonDriver_4:
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				a.driverCode = fmt.Sprintf("%s%d", a.driverCode, 4)
			case AddrGttButtonDriver_5:
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				a.driverCode = fmt.Sprintf("%s%d", a.driverCode, 5)
			case AddrGttButtonDriver_6:
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				a.driverCode = fmt.Sprintf("%s%d", a.driverCode, 6)
			case AddrGttButtonDriver_7:
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				a.driverCode = fmt.Sprintf("%s%d", a.driverCode, 7)
			case AddrGttButtonDriver_8:
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				a.driverCode = fmt.Sprintf("%s%d", a.driverCode, 8)
			case AddrGttButtonDriver_9:
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				a.driverCode = fmt.Sprintf("%s%d", a.driverCode, 9)
			case AddrGttButtonDriver_Clear:
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				a.driverCode = ""
			case AddrGttButtonDriver_Cancel:
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				a.driverCode = ""
				a.ctx.Send(a.ctx.Self(), &MsgMainScreen{})
			}

			return nil
		}(); err != nil {
			logs.LogWarn.Println(err)
		}
	}
}
