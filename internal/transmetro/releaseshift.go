package app

import (
	"context"
	"fmt"
	"time"

	"github.com/dumacp/go-logs/pkg/logs"
	"github.com/dumacp/go-schservices/api/services"
	"github.com/google/uuid"
)

func (a *App) releaseshift() error {
	// 	if a.shift == nil {
	// 		return fmt.Errorf(`no hay un turno
	// para liberar`)
	// 	}
	// selectedShift := a.shift
	a.shift = nil

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
	}

	funcRequest := func(mss interface{}) error {
		res, err := a.ctx.RequestFuture(a.pidSvc, mss, 10*time.Second).Result()
		if err != nil {
			return fmt.Errorf("request release service error: %s", err)
		}

		if resSvc, ok := res.(*services.ReleaseShiftResponseMsg); ok {

			if len(resSvc.GetError()) > 0 {
				logs.LogWarn.Printf("error release take shift: %s", err)
				return fmt.Errorf(resSvc.GetError())
			}
		} else {
			return fmt.Errorf("response type error (%T)", res)
		}

		if err := a.uix.TextConfirmationPopup("turno liberado\n"); err != nil {
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
	// case selectedShift == nil:
	// 	return fmt.Errorf("no hay un turno seleccionado")
	// case len(selectedShift.GetShift()) <= 0:
	// 	return fmt.Errorf("el iD del turno es invalido")
	default:
		mss := &services.ReleaseShiftMsg{
			DeviceId:   a.deviceId,
			PlatformId: a.platformId,
			// ShiftId:             selectedShift.GetShift(),
			// ServiceSchedulingId: selectedShift.ServiceSchedulingID,
			MessageId: uuid.New().String(),
			Timestamp: time.Now().UnixMilli(),
		}
		if err := funcRequest(mss); err != nil {
			return err
		}
	}

	return nil
}
