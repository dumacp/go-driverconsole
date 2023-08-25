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

func ButtonsPi(a *App) func(evt *buttons.InputEvent) {

	return func(evt *buttons.InputEvent) {

		if err := func() error {
			switchScreen := false
			switch evt.KeyCode {
			case AddrScreenSwitch:
				if err := a.uix.Screen(int(ui.MAIN_SCREEN), switchScreen); err != nil {
					return fmt.Errorf("event SCREEN error: %s", err)
				}
			case AddrScreenProgDriver:
				if err := a.uix.Screen(int(ui.PROGRAMATION_DRIVER_SCREEN), switchScreen); err != nil {
					return fmt.Errorf("event SCREEN error: %s", err)
				}
				// release button
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				ss := make([]*services.ScheduleService, 0)
				if len(a.shcservices) > 0 {
					for _, v := range a.shcservices {
						ss = append(ss, v)
					}
					sort.SliceStable(ss, func(i, j int) bool {
						return ss[i].GetScheduleDateTime() > ss[j].GetScheduleDateTime()
					})
				}
				reverseSlice := make([]string, 0)
				// reverseSlice = append(reverseSlice, "")
				for _, v := range ss {
					if v.GetItinenary() == nil || v.GetDriver() == nil {
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
						v.GetItinenary().GetId(), capitalize(strings.ToLower(v.GetDriver().GetFullname())))
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
			case AddrScreenProgVeh:
				if err := a.uix.Screen(int(ui.PROGRAMATION_VEH_SCREEN), switchScreen); err != nil {
					return fmt.Errorf("event SCREEN error: %s", err)
				}
				// release button
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				ss := make([]*services.ScheduleService, 0)
				if len(a.shcservices) > 0 {
					for _, v := range a.shcservices {
						ss = append(ss, v)
					}
					sort.SliceStable(ss, func(i, j int) bool {
						return ss[i].GetScheduleDateTime() > ss[j].GetScheduleDateTime()
					})

				}
				reverseSlice := make([]string, 0)
				// reverseSlice = append(reverseSlice, "")
				for _, v := range ss {
					if v.GetItinenary() == nil || v.GetRoute() == nil {
						continue
					}
					ts := time.UnixMilli(v.GetScheduleDateTime())

					data := strings.ToLower(fmt.Sprintf(" %s: (%d) %s (%s)", ts.Format("01/02 15:04"),
						v.GetItinenary().GetId(), v.GetItinenary().GetName(), v.GetRoute().GetName()))
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
			case AddrScreenAlarms:
				if err := a.uix.Screen(int(ui.NOTIFICATIONS_SCREEN), switchScreen); err != nil {
					return fmt.Errorf("event SCREEN error: %s", err)
				}
				// release button
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				if len(a.notif) > 0 {
					if err := a.uix.ShowNotifications(a.notif...); err != nil {
						return fmt.Errorf("event ShowNotifications error: %s", err)
					}
				}
			case AddrScreenMore:
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
			case AddrEnterRuta:
				// release button
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				if err := a.uix.SetLed(AddrEnterRuta, false); err != nil {
					return fmt.Errorf("error setLed (ROUTE_TEXT_READ): %s", err)
				}
				if v, err := a.uix.ReadBytesRawDisplay(ui.ROUTE_TEXT_READ); err != nil {
					return fmt.Errorf("error ReadBytesRawDisplay (ROUTE_TEXT_READ): %s", err)
				} else {
					data := strings.ReplaceAll(string(v), "\x00", "")
					if len(data) < 1 {
						return fmt.Errorf("error ReadBytesRawDisplay (len < 1): %s", data)
					}
					rutaCodeInt, err := strconv.Atoi(strings.TrimSpace(data))
					if err != nil {
						return fmt.Errorf("error route: %s", err)
					}
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
					if err := a.uix.Route(fmt.Sprintf(" %s", a.routeString)); err != nil {
						return fmt.Errorf("error Route: %s", err)
					}
				}
			case AddrEnterDriver:
				// release button
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				if err := a.uix.SetLed(AddrEnterDriver, false); err != nil {
					return fmt.Errorf("error setLed (Driver_TEXT_READ): %s", err)
				}
				if v, err := a.uix.ReadBytesRawDisplay(ui.DRIVER_TEXT_READ); err != nil {
					return fmt.Errorf("error ReadBytesRawDisplay (DRIVER_TEXT_READ): %s", err)
				} else {
					data := strings.ReplaceAll(string(v), "\x00", "")
					if len(data) < 6 {
						return fmt.Errorf("error ReadBytesRawDisplay (len < 6): %s", data)
					}
					fmt.Printf("driverCode: %s\n", data)
					driverCodeInt, err := strconv.Atoi(strings.TrimSpace(data))
					if err != nil {
						return fmt.Errorf("error driver: %s", err)
					}
					fmt.Printf("driverCode: %d\n", driverCodeInt)
					a.driver = driverCodeInt
					if err := a.uix.Driver(fmt.Sprintf(" %s", data)); err != nil {
						return fmt.Errorf("error Driver: %s", err)
					}
				}
			case AddrSwitchStep:
				if v, ok := evt.Value.(bool); ok {
					if !v {
						fmt.Println("/////////// step enable: true ////////////")
						a.enableStep = true
						if a.cancelStep != nil {
							a.cancelStep()
						}

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
							a.uix.StepEnable(false)
						}()
					} else {
						a.enableStep = false
						if a.cancelStep != nil {
							a.cancelStep()
						}
					}
				}
			case AddrSendStep:
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

			}
			return nil
		}(); err != nil {
			logs.LogWarn.Println(err)
		}
	}
}
