package app

import (
	"context"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dumacp/go-driverconsole/internal/buttons"
	"github.com/dumacp/go-driverconsole/internal/ui"
	"github.com/dumacp/go-logs/pkg/logs"
)

func ButtonsPi(a *App) func(evt *buttons.InputEvent) {

	return func(evt *buttons.InputEvent) {

		if evt.Error == nil && !a.isDisplayEnable {
			a.isDisplayEnable = true
			a.ctx.Send(a.ctx.Self(), &MsgMainScreen{})
		}

		if err := func() error {
			switchScreen := false
			fmt.Printf("event: %v\n", evt.Value)
			switch evt.KeyCode {
			case AddrEnterService:
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				if err := a.uix.SetLed(AddrShowSelectProgVeh, false); err != nil {
					return fmt.Errorf("error setLed (SELECT_PROG_VEH): %s", err)
				}
				// a.ctx.Send(a.ctx.Self(), &RequestTakeService{})
				a.ctx.Send(a.ctx.Self(), &RequestTakeShift{})
			case AddrScreenSwitch:
				// release button
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				if err := a.uix.SetLed(AddrScreenSwitch, false); err != nil {
					return fmt.Errorf("error setLed (AddrScreenSwitch): %s", err)
				}
				if err := a.uix.Screen(int(ui.MAIN_SCREEN), switchScreen); err != nil {
					return fmt.Errorf("event SCREEN error: %s", err)
				}
			case AddrShowSelectProgVeh:
				// release button
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				if err := a.uix.SetLed(AddrShowSelectProgVeh, false); err != nil {
					return fmt.Errorf("error setLed (SELECT_PROG_VEH): %s", err)
				}
				if v, err := a.uix.ReadBytesRawDisplay(AddrCurrentSelectProgVeh); err != nil {
					return fmt.Errorf("error ReadBytesRawDisplay (SELECT_PROG_VEH): %s", err)
				} else {
					fmt.Printf("data SELECT_PROG_VEH: %s\n", v)
					num := binary.LittleEndian.Uint16(v)
					fmt.Printf("num SELECT_PROG_VEH: %d\n", num)
					if len(a.companyShiftsShow) > int(num) {
						fmt.Printf("companySchServices: %v\n", a.companyShiftsShow[num])
						a.selectedShift = a.companyShiftsShow[num].Shift
						a.uix.WriteTextRawDisplay(AddrResumeSelectProgVeh, []string{a.companyShiftsShow[num].String})
						prompt := fmt.Sprintf(`turno preseleccionado:
%s`, a.companyShiftsShow[num].ResumeString)
						if err := a.uix.WriteTextRawDisplay(AddrTextCurrentItinerary, []string{prompt}); err != nil {
							logs.LogWarn.Printf("error TextCurrentItinerary: %s", err)
						}
						// if err := a.uix.Route(fmt.Sprintf("%d", a.companySchServicesShow[num].Services.Itinerary.Id)); err != nil {
						// 	fmt.Printf("error Route: %s\n", err)
						// }
					}
				}
			case AddrSelectItinerary:
				// release button
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				if err := a.uix.SetLed(AddrSelectItinerary, false); err != nil {
					return fmt.Errorf("error setLed (SELECT_ITI_VEH): %s", err)
				}
				if v, err := a.uix.ReadBytesRawDisplay(AddrItineraryProgVeh); err != nil {
					return fmt.Errorf("error ReadBytesRawDisplay (SELECT_ITI_VEH): %s", err)
				} else {
					data := uint64(0)
					switch len(v) {
					case 0:
					case 2:
						data = uint64(binary.LittleEndian.Uint16(v))
					case 4:
						data = uint64(binary.LittleEndian.Uint32(v))
					case 8:
						data = binary.LittleEndian.Uint64(v)
					}
					if err := func() error {
						fmt.Printf("data SELECT_ITI_VEH: %d\n", data)
						codeShift := int(data)
						fmt.Printf("ShiftCode: %d\n", codeShift)

						a.ctx.Send(a.ctx.Self(), &RequestShitfsVeh{
							Shift: fmt.Sprintf("%d", codeShift),
						})

						return nil
					}(); err != nil {
						logs.LogWarn.Printf("error route: %s", err)
						if err := a.uix.TextWarningPopup(fmt.Sprintf("%s\n", err)); err != nil {
							logs.LogWarn.Printf("textWarningPopup error: %s", err)
						}
						if a.cancelPop != nil {
							a.cancelPop()
						}
						contxt, cancel := context.WithCancel(context.TODO())
						a.cancelPop = cancel
						go func() {
							defer cancel()
							select {
							case <-contxt.Done():
							case <-time.After(4 * time.Second):
							}
							if err := a.uix.TextWarningPopupClose(); err != nil {
								logs.LogWarn.Printf("textWarningPopupClose error: %s", err)
							}
						}()
						return err
					}
				}

			case AddrScreenProgDriver:
				// release button
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				if err := a.uix.SetLed(AddrScreenProgDriver, false); err != nil {
					return fmt.Errorf("error setLed (AddrScreenProgDriver): %s", err)
				}
				if err := a.uix.Screen(int(ui.PROGRAMATION_DRIVER_SCREEN), switchScreen); err != nil {
					return fmt.Errorf("event SCREEN error: %s", err)
				}
				a.ctx.Send(a.ctx.Self(), &ListProgDriver{})

			case AddrScreenProgVeh:
				// release button
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				if err := a.uix.SetLed(AddrScreenProgVeh, false); err != nil {
					return fmt.Errorf("error setLed (AddrScreenProgVeh): %s", err)
				}
				if err := a.uix.Screen(int(ui.PROGRAMATION_VEH_SCREEN), switchScreen); err != nil {
					return fmt.Errorf("event SCREEN error: %s", err)
				}

				dataSlice := make([]string, 0)
				for i := 0; i <= 9; i++ {
					size := Label2DisplayRegister(ui.PROGRAMATION_VEH_SCREEN).Size
					// un string de tamaño size de espacios
					spaces := strings.Repeat(" ", size)
					dataSlice = append(dataSlice, spaces)
				}
				if len(dataSlice) > 0 {
					if err := a.uix.ShowProgVeh(dataSlice...); err != nil {
						fmt.Printf("clean event ShowProgVeh error: %s", err)
					}
				}

				// a.ctx.Send(a.ctx.Self(), &RequestProgVeh{})
				a.ctx.Send(a.ctx.Self(), &RequestShitfsVeh{})

			case AddrScreenAlarms:
				// release button
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				if err := a.uix.SetLed(AddrScreenAlarms, false); err != nil {
					return fmt.Errorf("error setLed (AddrScreenAlarms): %s", err)
				}
				if err := a.uix.Screen(int(ui.ALARMS_SCREEN), switchScreen); err != nil {
					return fmt.Errorf("event SCREEN error: %s", err)
				}
				if len(a.notif) > 0 {
					if err := a.uix.ShowNotifications(a.notif...); err != nil {
						return fmt.Errorf("event ShowNotifications error: %s", err)
					}
				}
			case AddrScreenMore:
				// release button
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				if err := a.uix.SetLed(AddrScreenMore, false); err != nil {
					return fmt.Errorf("error setLed (AddrScreenMore): %s", err)
				}
				if err := a.uix.Screen(int(ui.ADDITIONALS_SCREEN), switchScreen); err != nil {
					return fmt.Errorf("event SCREEN error: %s", err)
				}
				if err := a.uix.ShowStats(); err != nil {
					return fmt.Errorf("event ShowStats error: %s", err)
				}
			case AddrExitSwitch:
				// release button
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				if err := a.uix.SetLed(AddrExitSwitch, false); err != nil {
					return fmt.Errorf("error setLed (AddrExitSwitch): %s", err)
				}
				a.ctx.Send(a.ctx.Self(), &ReleaseShitfsVeh{})

			case AddrEnterRuta:
				// release button
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				if err := a.uix.SetLed(AddrEnterRuta, false); err != nil {
					return fmt.Errorf("error setLed (ROUTE_TEXT_READ): %s", err)
				}
				// if v, err := a.uix.ReadBytesRawDisplay(ui.ROUTE_TEXT_READ); err != nil {
				// 	return fmt.Errorf("error ReadBytesRawDisplay (ROUTE_TEXT_READ): %s", err)
				// } else {
				// 	data := strings.ReplaceAll(string(v), "\x00", "")
				// 	if len(data) < 1 {
				// 		return fmt.Errorf("error ReadBytesRawDisplay (len < 1): %s", data)
				// 	}
				// 	rutaCodeInt, err := strconv.Atoi(strings.TrimSpace(data))
				// 	if err != nil {
				// 		return fmt.Errorf("error route: %s", err)
				// 	}
				// 	a.route = rutaCodeInt
				// 	if a.routes != nil {
				// 		if routeString, ok := a.routes[int32(rutaCodeInt)]; ok {
				// 			a.routeString = routeString
				// 			a.ctx.Send(a.ctx.Self(), &MsgSetRoute{
				// 				Route: rutaCodeInt,
				// 			})
				// 			routeS := func() string {
				// 				if len(a.routeString) > 32 {
				// 					return a.routeString[:32]
				// 				}
				// 				return a.routeString
				// 			}()
				// 			if err := a.uix.Route(routeS); err != nil {
				// 				return fmt.Errorf("error Route: %s", err)
				// 			}
				// 		}
				// 	}
				// }
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
					// fmt.Printf("driverCode: %d\n", data)
					if a.driver == nil || !strings.EqualFold(a.driver.DocumentId, data) {
						a.ctx.Send(a.ctx.Self(), &MsgSetDriver{
							Driver: driverCodeInt,
						})
					}
					// if err := a.uix.Driver(fmt.Sprintf(" %s", data)); err != nil {
					// 	return fmt.Errorf("error Driver: %s", err)
					// }
				}
			default:
				if evt.Error != nil {
					a.ctx.Send(a.ctx.Self(), &ErrorDisplay{
						Error: evt.Error,
					})
				}

			}
			return nil
		}(); err != nil {
			// fmt.Printf("error event: %s\n", err)
			logs.LogWarn.Println("error event: ", err)
		}
	}
}
