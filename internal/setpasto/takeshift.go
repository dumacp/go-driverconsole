package app

import (
	"context"
	"fmt"
	"time"

	"github.com/dumacp/go-logs/pkg/logs"
	"github.com/dumacp/go-schservices/api/services"
	"github.com/google/uuid"
)

func (a *App) takeshift() error {
	if a.selectedShift == nil {
		return fmt.Errorf(`no hay un turno
sobre el cual iniciar`)
	}
	selectedShift := a.selectedShift
	// a.selectedShift = nil

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

		if resSvc, ok := res.(*services.TakeShiftResponseMsg); ok {

			if len(resSvc.GetError()) > 0 {
				logs.LogWarn.Printf("error request take shift: %s", err)
				return fmt.Errorf(resSvc.GetError())
			}
		} else {
			return fmt.Errorf("response type error (%T)", res)
		}

		if err := a.uix.TextConfirmationPopup("turno iniciado\n"); err != nil {
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
		a.selectedShift = nil
		return nil
	}

	switch {
	case selectedShift == nil:
		return fmt.Errorf("no hay un turno seleccionado")
	case len(selectedShift.GetShift()) <= 0:
		return fmt.Errorf("el iD del turno es invalido")
	case selectedShift.Itinerary == nil:
		return fmt.Errorf(`no hay un itinerario seleccionado
dentro del turno`)
	case len(selectedShift.ServiceSchedulingID) <= 0:
		return fmt.Errorf(`no hay un servicio programado
dentro del turno`)
	case a.currentService != nil && selectedShift.ServiceSchedulingID == a.currentService.Id:
		return fmt.Errorf(`el turno ya esta iniciado`)
	case selectedShift.GetServiceAmount() <= 0:
		return fmt.Errorf(`el turno no tiene servicios programados`)
	// case a.currentService != nil && a.selectedShift.ServiceSchedulingID != a.currentService.Id:
	default:
		mss := &services.TakeShiftMsg{
			DeviceId:            a.deviceId,
			PlatformId:          a.platformId,
			DriverId:            a.driver.Id,
			ShiftId:             selectedShift.GetShift(),
			ServiceSchedulingId: selectedShift.ServiceSchedulingID,
			MessageId:           uuid.New().String(),
			Timestamp:           time.Now().UnixMilli(),
		}
		if err := funcRequest(mss); err != nil {
			return err
		}
		a.ctx.Send(a.ctx.Self(), &MsgSetRoute{
			Route:     int(selectedShift.GetItinerary().GetId()),
			RouteName: selectedShift.GetItinerary().GetName(),
		})
	}

	a.shift = selectedShift

	return nil
}
