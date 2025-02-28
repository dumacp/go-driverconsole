package app

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/dumacp/go-logs/pkg/logs"
	"github.com/dumacp/go-schservices/api/services"
)

func (a *App) takeservice() error {
	if a.selectedService == nil && a.currentService == nil {
		return fmt.Errorf("no hay un servicio seleccionado")
	}

	selectedService := a.currentService
	if a.selectedService != nil {
		selectedService = a.selectedService
	}
	a.selectedService = nil
	if a.pidSvc == nil {
		return fmt.Errorf("service pid is nil")
	}
	switch {
	case len(a.deviceId) <= 0:
		return fmt.Errorf("device id is empty")
	case len(a.platformId) <= 0:
		return fmt.Errorf("platform id is empty")
	case len(a.companyId) <= 0:
		return fmt.Errorf("company id is empty")
	case a.driver == nil || len(a.driver.DocumentId) <= 0:
		return fmt.Errorf("conductor no seleccionado")
	}

	funcRequest := func(mss interface{}) error {
		res, err := a.ctx.RequestFuture(a.pidSvc, mss, 10*time.Second).Result()
		if err != nil {
			return fmt.Errorf("request service error: %s", err)
		}
		// if resSvc, ok := res.(*services.StartServiceResponseMsg); ok {
		// 	if len(resSvc.GetError()) > 0 {
		// 		return fmt.Errorf("error request service: %s", resSvc.GetError())
		// 	}
		// 	if resSvc.DataCode != 200 {
		// 		return fmt.Errorf("error request service: %d, %s", resSvc.DataCode, resSvc.DataMsg)
		// 	}
		// 	if err := a.uix.TextConfirmationPopup("servicio iniciado"); err != nil {
		// 		logs.LogWarn.Printf("textConfirmation error: %s", err)
		// 	}
		// } else
		if resSvc, ok := res.(*services.TakeServiceResponseMsg); ok {

			if len(resSvc.GetError()) > 0 {
				logs.LogWarn.Printf("error request service: %s", err)
				return fmt.Errorf(resSvc.GetError())
			}
			// if resSvc.DataCode != 200 {
			// 	logs.LogWarn.Printf("error request service: %s", err)
			// 	return fmt.Errorf(resSvc.DataMsg)
			// }
			// a.currentService = a.selectedService

		} else {
			return fmt.Errorf("response type error")
		}

		if err := a.uix.TextConfirmationPopup("servicio iniciado\n"); err != nil {
			logs.LogWarn.Printf("textConfirmation error: %s", err)
		}
		if a.cancelPop != nil {
			a.cancelPop()
		}
		contxt, cancel := context.WithCancel(context.Background())
		a.cancelPop = cancel
		go func() {
			defer cancel()
			select {
			case <-contxt.Done():
			case <-time.After(4 * time.Second):
			}
			if err := a.uix.TextConfirmationPopupclose(); err != nil {

				logs.LogWarn.Printf("textConfirmation error: %s", err)
			}
		}()
		return nil
	}

	switch {
	case selectedService == nil:
		return fmt.Errorf("no hay un servicio iniciado")
	case a.currentService != selectedService &&
		selectedService.State == services.State_SCHEDULED.String():
		mss := &services.TakeServiceMsg{
			DeviceId:   a.deviceId,
			PlatformId: a.platformId,
			CompanyId:  a.companyId,
			ServiceId:  selectedService.Id,
			DriverId:   a.driver.Id,
		}
		if err := funcRequest(mss); err != nil {
			return err
		}
		a.ctx.Send(a.ctx.Self(), &MsgSetRoute{
			Route:     int(selectedService.GetItinerary().GetId()),
			RouteName: selectedService.GetItinerary().GetName(),
		})
	case a.currentService == nil:
		return fmt.Errorf("no hay un servicio iniciado")
	case a.currentService.State == services.State_STARTED.String():
		return fmt.Errorf("servicio ya iniciado")
	case a.currentService.State == services.State_ENDED.String():
		return fmt.Errorf("servicio ya finalizado")
	case a.currentService.State == services.State_ABORTED.String():
		return fmt.Errorf("servicio ya abortado")
	case a.currentService.State == services.State_WAITING_TO_ARRIVE_TO_STARTING_POINT.String() ||
		a.currentService.State == services.State_READY_TO_START.String() ||
		a.currentService.State == services.State_SCHEDULED.String():
		mss := &services.TakeServiceMsg{
			DeviceId:   a.deviceId,
			PlatformId: a.platformId,
			CompanyId:  a.companyId,
			ServiceId:  a.currentService.Id,
			DriverId:   a.driver.Id,
		}
		if err := funcRequest(mss); err != nil {
			return err
		}
		a.ctx.Send(a.ctx.Self(), &MsgSetRoute{
			Route:     int(selectedService.GetItinerary().GetId()),
			RouteName: selectedService.GetItinerary().GetName(),
		})
	default:
		fmt.Printf("current service: %v\n", a.currentService)
		return fmt.Errorf("servicio no programado")
	}

	return nil
}

func (a *App) showCurrentServiceWithAll(msg *services.ServiceAllMsg) {
	svcs := msg.GetUpdates()
	fmt.Printf("services len: %v\n", len(svcs))
	a.shcservices = make(map[string]*services.ScheduleService)
	startOfDay := time.Now().Truncate(24 * time.Hour)
	for _, svc := range svcs {
		if !time.UnixMilli(svc.GetScheduleDateTime()).Before(startOfDay) &&
			time.UnixMilli(svc.GetScheduleDateTime()).Before(startOfDay.Add(24*time.Hour)) {
			a.shcservices[svc.GetId()] = svc
			fmt.Printf("service in of range: %v (%s), %s\n", svc.Id, svc.GetState(), time.UnixMilli(svc.GetScheduleDateTime()))
		} else {
			fmt.Printf("service out of range: %v (%s), %s\n", svc.Id, svc.GetState(), time.UnixMilli(svc.GetScheduleDateTime()))
		}
	}

	arr := make([]string, 0)
	ss := make([]*services.ScheduleService, 0)
	for k, v := range a.shcservices {
		ss = append(ss, v)
		arr = append(arr, k)
	}
	sort.SliceStable(ss, func(i, j int) bool {
		return ss[i].GetScheduleDateTime() < ss[j].GetScheduleDateTime()
	})

	tn := time.Now().UnixMilli()
	idx := sort.Search(len(ss), func(i int) bool {
		return ss[i].GetScheduleDateTime() > tn // Busca el primer valor que sea mayor a `target`
	})

	if idx-1 > 0 &&
		(ss[idx-1].GetState() == services.State_STARTED.String()) {
		prompt := ""
		v := ss[idx]
		fmt.Printf("************ service: %v\n", v)
		a.currentService = v
		ts := time.UnixMilli(v.GetScheduleDateTime())
		prompt = strings.ToLower(fmt.Sprintf("servicio iniciado:\n%s: %s (%s)", ts.Format("01/02 15:04"),
			v.GetItinerary().GetName(), v.GetRoute().GetCode()))
		if err := a.uix.WriteTextRawDisplay(AddrTextCurrentItinerary, []string{prompt}); err != nil {
			fmt.Printf("error TextCurrentItinerary: %s\n", err)
		}
		fmt.Printf("//////////////// services (ori: %d): %v\n", len(msg.GetUpdates()), arr)
	} else if idx > 0 && idx < len(ss) {
		prompt := ""
		v := ss[idx]
		fmt.Printf("************ service: %v\n", v)
		if v.GetState() == services.State_STARTED.String() {
			a.currentService = v
			ts := time.UnixMilli(v.GetScheduleDateTime())
			prompt = strings.ToLower(fmt.Sprintf("servicio iniciado:\n%s: %s (%s)", ts.Format("01/02 15:04"),
				v.GetItinerary().GetName(), v.GetRoute().GetCode()))
		} else if v.GetState() == services.State_READY_TO_START.String() {
			a.currentService = v
			ts := time.UnixMilli(v.GetScheduleDateTime())
			prompt = strings.ToLower(fmt.Sprintf("próximo servicio (listo):\n%s: %s (%s)", ts.Format("01/02 15:04"),
				// prompt = strings.ToLower(fmt.Sprintf("servicio iniciado:\n%s: %s (%s)", ts.Format("01/02 15:04"),
				v.GetItinerary().GetName(), v.GetRoute().GetCode()))
		} else if v.GetState() == services.State_WAITING_TO_ARRIVE_TO_STARTING_POINT.String() {
			if a.currentService == nil || a.currentService.GetState() == services.State_ENDED.String() ||
				a.currentService.GetState() == services.State_ABORTED.String() ||
				a.currentService.GetState() == services.State_CANCELLED.String() {
				a.currentService = v
				ts := time.UnixMilli(v.GetScheduleDateTime())
				prompt = strings.ToLower(fmt.Sprintf("próximo servicio (esperando):\n%s: %s (%s)", ts.Format("01/02 15:04"),
					v.GetItinerary().GetName(), v.GetRoute().GetCode()))
			}
		} else if v.GetState() == services.State_SCHEDULED.String() {
			// a.currentService = v
			if a.currentService == nil || a.currentService.GetState() == services.State_ENDED.String() ||
				a.currentService.GetState() == services.State_ABORTED.String() {
				a.currentService = v
				ts := time.UnixMilli(v.GetScheduleDateTime())
				prompt = strings.ToLower(fmt.Sprintf("próximo servicio:\n%s: %s (%s)", ts.Format("01/02 15:04"),
					v.GetItinerary().GetName(), v.GetRoute().GetCode()))
			}
		}
		if err := a.uix.WriteTextRawDisplay(AddrTextCurrentItinerary, []string{prompt}); err != nil {
			fmt.Printf("error TextCurrentItinerary: %s\n", err)
		}
		fmt.Printf("//////////////// services (ori: %d): %v\n", len(msg.GetUpdates()), arr)
	}
}
