package app

import (
	"fmt"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/dumacp/go-driverconsole/internal/buttons"
	"github.com/dumacp/go-driverconsole/internal/ui"
	"github.com/dumacp/go-logs/pkg/logs"
	"github.com/dumacp/go-schservices/api/services"
)

type Event struct {
	Label EventLabel
	Value interface{}
}

type EventLabel int

const (
	PROGRAMATION_DRIVER EventLabel = iota
	PROGRAMATION_VEH
	STATS
	SHOW_NOTIF
	ROUTE
	DRIVER
	SCREEN_SWITCH
)

func (a *App) Buttons() func(evt *buttons.InputEvent) {

	return func(evt *buttons.InputEvent) {

		evt2EvtLabel := map[int]EventLabel{
			AddrEnterRuta:        ROUTE,
			AddrEnterDriver:      DRIVER,
			AddrScreenSwitch:     SCREEN_SWITCH,
			AddrScreenProgVeh:    PROGRAMATION_VEH,
			AddrScreenProgDriver: PROGRAMATION_DRIVER,
			AddrScreenAlarms:     SHOW_NOTIF,
		}

		label, ok := evt2EvtLabel[int(evt.KeyCode)]
		if !ok {
			return
		}
		if err := func() error {
			switch label {
			case PROGRAMATION_DRIVER:
				// release button
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				if err := a.uix.Screen(int(ui.PROGRAMATION_DRIVER_SCREEN), false); err != nil {
					return fmt.Errorf("event SCREEN error: %s", err)
				}
				// if err := a.uix.ShowProgDriver(); err != nil {
				// 	return fmt.Errorf("event ShowProgDriver error: %s", err)
				// }
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
			case PROGRAMATION_VEH:
				if err := a.uix.Screen(int(ui.PROGRAMATION_VEH_SCREEN), false); err != nil {
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
			case SHOW_NOTIF:
				if err := a.uix.Screen(int(ui.ALARMS_SCREEN), false); err != nil {
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
			case STATS:
				if err := a.uix.Screen(int(ui.ADDITIONALS_SCREEN), false); err != nil {
					return fmt.Errorf("event SCREEN error: %s", err)
				}
				// release button
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				if err := a.uix.ShowStats(); err != nil {
					return fmt.Errorf("event ShowStats error: %s", err)
				}
			case ROUTE:
				// release button
				// if v, ok := evt.Value.(bool); !ok || v {
				// 	break
				// }
				// contxt, cancel := context.WithCancel(a.contxt)
				// a.buttonCancel = cancel
				// num, err := a.uix.KeyNum(contxt, "ingrese el número de ruta:")
				// if err != nil {
				// 	return fmt.Errorf("route keyNum error: %s", err)
				// }
				// go func() {
				// 	defer cancel()
				// 	self := a.ctx.Self()
				// 	rootctx := a.ctx.ActorSystem().Root

				// 	select {
				// 	case <-contxt.Done():
				// 	case v := <-num:
				// 		rootctx.Send(self, &MsgChangeRoute{
				// 			ID: v,
				// 		})
				// 	}
				// }()
			case DRIVER:
				// release button
				if v, ok := evt.Value.(bool); !ok || v {
					break
				}
				if v, err := a.uix.ReadBytesRawDisplay(ui.DRIVER_TEXT_READ); err != nil {
					return fmt.Errorf("error ReadBytesRawDisplay (DRIVER_TEXT_READ): %s", err)
				} else {
					data := strings.ReplaceAll(string(v), "\x00", "")
					if len(data) < 6 {
						return fmt.Errorf("error ReadBytesRawDisplay (len < 6): %s", data)
					}
					if err := a.uix.Driver(fmt.Sprintf(" %s", data)); err != nil {
						return fmt.Errorf("error Driver: %s", err)
					}
				}

			}
			return nil
		}(); err != nil {
			logs.LogWarn.Println(err)
		}
	}
}
