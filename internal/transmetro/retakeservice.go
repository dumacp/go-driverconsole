package app

import (
	"context"
	"fmt"
	"time"

	"github.com/dumacp/go-logs/pkg/logs"
	"github.com/dumacp/go-schservices/api/services"
)

func (a *App) retakeservice() error {
	if a.currentService == nil {
		return fmt.Errorf(`no hay un servicio iniciado
sobre el cual retomar`)
	}

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
			return fmt.Errorf("request retake service error: %s", err)
		}

		if resSvc, ok := res.(*services.TakeServiceResponseMsg); ok {

			if len(resSvc.GetError()) > 0 {
				logs.LogWarn.Printf("error request retake service: %s", err)
				return fmt.Errorf(resSvc.GetError())
			}
		} else {
			return fmt.Errorf("response type error (%T)", res)
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
	case a.currentService == nil:
		return fmt.Errorf("no hay un servicio iniciado")
	case a.currentService.State == services.State_WAITING_TO_ARRIVE_TO_STARTING_POINT.String():
		return fmt.Errorf("servicio actual esperando llegar al punto de inicio")
	case a.currentService.State == services.State_SCHEDULED.String():
		return fmt.Errorf("servicio actual programado")
	case a.currentService.State == services.State_READY_TO_START.String():
		return fmt.Errorf("servicio actual listo para iniciar")
	case a.currentService.State == services.State_ABORTED.String() ||
		a.currentService.State == services.State_STARTED.String() ||
		a.currentService.State == services.State_CANCELLED.String() ||
		a.currentService.State == services.State_ENDED.String():
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
			Route:     int(a.currentService.GetItinerary().GetId()),
			RouteName: a.currentService.GetItinerary().GetName(),
		})
	default:
		fmt.Printf("current service: %v\n", a.currentService)
		return fmt.Errorf("servicio anterior no finalizado")
	}

	return nil
}
