package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/dumacp/go-logs/pkg/logs"
	"github.com/dumacp/go-schservices/api/services"
)

func (a *App) showCurrentService(svc *services.ScheduleService) {

	// changeCurrentService := false
	state := svc.GetState()
	if len(state) > 0 {
		if v, ok := a.shcservices[svc.GetId()]; ok {
			UpdateService(v, svc)
			fmt.Printf("////// update: %v\n", v)
		} else {
			a.shcservices[svc.GetId()] = svc
		}
		svc = a.shcservices[svc.GetId()]
		if svc.GetScheduleDateTime() <= 0 {
			svc.ScheduleDateTime = time.Now().UnixMilli()
		}
		startOfDay := time.Now().Truncate(24 * time.Hour)
		if time.UnixMilli(svc.GetScheduleDateTime()).Before(startOfDay) ||
			!time.UnixMilli(svc.GetScheduleDateTime()).Before(startOfDay.Add(24*time.Hour)) {
			fmt.Printf("service mod out of range: %v (%s), %s\n", svc.Id, svc.GetState(), time.UnixMilli(svc.GetScheduleDateTime()))
			return
		}
		a.companySchServices[svc.GetId()] = svc

		if !strings.EqualFold(state, svc.GetState()) {
			data := strings.ToLower(fmt.Sprintf(" %s: (%d) %s (%s)", time.Now().Format("01/02 15:04"),
				svc.GetItinerary().GetId(), svc.GetItinerary().GetName(), svc.GetState()))
			a.notif = append(a.notif, data)
			if len(a.notif) > 10 {
				copy(a.notif, a.notif[1:])
				a.notif = a.notif[:len(a.notif)-1]
			}
			fmt.Printf("notif len: %d, %v\n", len(a.notif), a.notif)
		}
	}

	prompt := ""
	promptNext := ""

	if svc.GetCheckpointTimingState() != nil && len(svc.GetCheckpointTimingState().GetState()) > 0 {
		if a.currentService == nil {
			a.currentService = svc
			if len(a.shcservices) > 0 && a.shcservices[svc.GetId()] != nil {
				a.currentService = a.shcservices[svc.GetId()]
			}
			ts := time.Now()
			if a.currentService.GetScheduleDateTime() > 0 {
				ts = time.UnixMilli(a.currentService.GetScheduleDateTime())
			}
			if len(a.currentService.GetState()) <= 0 {
				a.currentService.State = services.State_UNKNOWN.String()
				prompt = strings.ToLower(fmt.Sprintf("servicio iniciado:\n%s: %s (%s)", ts.Format("01/02 15:04"),
					svc.GetItinerary().GetName(), svc.GetRoute().GetCode()))
			}
		}
		state := int(services.TimingState_value[svc.GetCheckpointTimingState().GetState()])
		promtp := ""
		if len(svc.GetCheckpointTimingState().GetName()) > 30 {
			promtp = fmt.Sprintf("%s (%d)", svc.GetCheckpointTimingState().GetName()[0:30], svc.GetCheckpointTimingState().GetTimeDiff())
		} else {
			promtp = fmt.Sprintf("%s (%d)", svc.GetCheckpointTimingState().GetName(), svc.GetCheckpointTimingState().GetTimeDiff())
		}
		if len(svc.CheckpointTimingState.GetNextCheckPointName()) > 20 {
			promtp = fmt.Sprintf(`%s
		prox: %s`, promtp, svc.CheckpointTimingState.GetNextCheckPointName()[0:20])
		} else if len(svc.CheckpointTimingState.GetNextCheckPointName()) > 0 {
			promtp = fmt.Sprintf(`%s
			prox: %s`, promtp, svc.CheckpointTimingState.GetNextCheckPointName())
		}
		fmt.Printf("///// state: %d\n", state)
		if err := a.uix.ServiceCurrentState(state, promtp); err != nil {
			logs.LogWarn.Printf("textConfirmation error: %s", err)
		}
	}

	fmt.Printf("************ service: %v\n", svc)
	if state == services.State_STARTED.String() {
		a.currentService = svc
		a.lastService = nil
		a.nextService = nil
		ts := time.UnixMilli(svc.GetScheduleDateTime())
		prompt = strings.ToLower(fmt.Sprintf("servicio iniciado:\n%s: %s (%s)", ts.Format("01/02 15:04"),
			svc.GetItinerary().GetName(), svc.GetRoute().GetCode()))
	} else if state == services.State_UNKNOWN.String() {
		ts := time.Now()
		if a.currentService.GetScheduleDateTime() > 0 {
			ts = time.UnixMilli(a.currentService.GetScheduleDateTime())
		}
		if len(a.currentService.GetState()) <= 0 {
			a.currentService.State = services.State_UNKNOWN.String()
			prompt = strings.ToLower(fmt.Sprintf("servicio iniciado:\n%s: %s (%s)", ts.Format("01/02 15:04"),
				svc.GetItinerary().GetName(), svc.GetRoute().GetCode()))
		}
	} else if state == services.State_READY_TO_START.String() {
		a.currentService = svc
		ts := time.UnixMilli(svc.GetScheduleDateTime())
		prompt = strings.ToLower(fmt.Sprintf("próximo servicio (listo):\n%s: %s (%s)", ts.Format("01/02 15:04"),
			// prompt = strings.ToLower(fmt.Sprintf("servicio iniciado:\n%s: %s (%s)", ts.Format("01/02 15:04"),
			svc.GetItinerary().GetName(), svc.GetRoute().GetCode()))
	} else if state == services.State_WAITING_TO_ARRIVE_TO_STARTING_POINT.String() {
		if a.currentService == nil || a.currentService.GetState() == services.State_ENDED.String() ||
			a.currentService.GetState() == services.State_ABORTED.String() ||
			a.currentService.GetState() == services.State_CANCELLED.String() ||
			a.currentService.GetState() == services.State_UNKNOWN.String() {
			a.currentService = svc
			ts := time.UnixMilli(svc.GetScheduleDateTime())
			prompt = strings.ToLower(fmt.Sprintf("próximo servicio (esperando):\n%s: %s (%s)", ts.Format("01/02 15:04"),
				svc.GetItinerary().GetName(), svc.GetRoute().GetCode()))
		} else if a.currentService.GetState() == services.State_STARTED.String() {
			a.nextService = svc
		} else {
			a.nextService = nil
		}
	} else if state == services.State_SCHEDULED.String() {
		// a.currentService = v
		if a.currentService == nil || a.currentService.GetState() == services.State_ENDED.String() ||
			a.currentService.GetState() == services.State_ABORTED.String() ||
			a.currentService.GetState() == services.State_CANCELLED.String() ||
			a.currentService.GetState() == services.State_UNKNOWN.String() {
			a.currentService = svc
			ts := time.UnixMilli(svc.GetScheduleDateTime())
			prompt = strings.ToLower(fmt.Sprintf("próximo servicio:\n%s: %s (%s)", ts.Format("01/02 15:04"),
				svc.GetItinerary().GetName(), svc.GetRoute().GetCode()))
		}
	} else if state == services.State_ENDED.String() {
		// a.currentService = v
		ts := time.UnixMilli(svc.GetScheduleDateTime())
		prompt = strings.ToLower(fmt.Sprintf("servicio finalizado:\n%s: %s (%s)", ts.Format("01/02 15:04"),
			svc.GetItinerary().GetName(), svc.GetRoute().GetCode()))
		if a.nextService != nil && a.nextService.GetId() != svc.GetId() {
			ts := time.UnixMilli(a.nextService.GetScheduleDateTime())
			promptNext = strings.ToLower(fmt.Sprintf("próximo servicio:\n%s: %s (%s)", ts.Format("01/02 15:04"),
				a.nextService.GetItinerary().GetName(), a.nextService.GetRoute().GetCode()))
			a.currentService = a.nextService
			a.nextService = nil
		}
	} else if state == services.State_ABORTED.String() {
		// a.currentService = v
		ts := time.UnixMilli(svc.GetScheduleDateTime())
		if a.currentService != nil && a.currentService.GetId() == svc.GetId() {
			prompt = strings.ToLower(fmt.Sprintf("servicio abortado:\n%s: %s (%s)", ts.Format("01/02 15:04"),
				svc.GetItinerary().GetName(), svc.GetRoute().GetCode()))
			if a.nextService != nil && a.nextService.GetId() != svc.GetId() {
				ts := time.UnixMilli(a.nextService.GetScheduleDateTime())
				promptNext = strings.ToLower(fmt.Sprintf("próximo servicio:\n%s: %s (%s)", ts.Format("01/02 15:04"),
					a.nextService.GetItinerary().GetName(), a.nextService.GetRoute().GetCode()))
				a.currentService = a.nextService
				a.nextService = nil
			}
		}
	} else if state == services.State_CANCELLED.String() {
		// a.currentService = v
		ts := time.UnixMilli(svc.GetScheduleDateTime())
		if a.currentService != nil && a.currentService.GetId() == svc.GetId() {
			prompt = strings.ToLower(fmt.Sprintf("servicio cancelado:\n%s: %s (%s)", ts.Format("01/02 15:04"),
				svc.GetItinerary().GetName(), svc.GetRoute().GetCode()))
			if a.nextService != nil && a.nextService.GetId() != svc.GetId() {
				ts := time.UnixMilli(a.nextService.GetScheduleDateTime())
				promptNext = strings.ToLower(fmt.Sprintf("próximo servicio:\n%s: %s (%s)", ts.Format("01/02 15:04"),
					a.nextService.GetItinerary().GetName(), a.nextService.GetRoute().GetCode()))
				a.currentService = a.nextService
				a.nextService = nil
			}
		}
	} else if len(state) == 0 {
		if svc.GetCheckpointTimingState().GetState() == services.TimingState_ON_TIME.String() {
			if a.currentService != nil && svc.GetId() == a.currentService.GetId() {
				if a.currentService.GetState() == services.State_SCHEDULED.String() ||
					a.currentService.GetState() == services.State_UNKNOWN.String() ||
					a.currentService.GetState() == services.State_WAITING_TO_ARRIVE_TO_STARTING_POINT.String() {
					a.currentService.State = services.State_STARTED.String()
					ts := time.UnixMilli(svc.GetScheduleDateTime())
					prompt = strings.ToLower(fmt.Sprintf("servicio iniciado:\n%s: %s (%s)", ts.Format("01/02 15:04"),
						svc.GetItinerary().GetName(), svc.GetRoute().GetCode()))
				}
			}
		}
	}

	if len(prompt) > 0 {
		if err := a.uix.WriteTextRawDisplay(AddrTextCurrentItinerary, []string{prompt}); err != nil {
			logs.LogWarn.Printf("error TextCurrentItinerary: %s", err)
		}
	}

	if len(promptNext) > 0 {
		time.Sleep(2 * time.Second)
		if err := a.uix.WriteTextRawDisplay(AddrTextCurrentItinerary, []string{prompt}); err != nil {
			logs.LogWarn.Printf("error TextCurrentItinerary: %s", err)
		}
	}
}
