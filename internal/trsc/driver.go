package app

import (
	"context"
	"fmt"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/dumacp/go-schservices/api/services"
)

func (a *App) setDriver(ctx actor.Context, msg *MsgSetDriver) error {
	if msg.Driver <= 0 {
		return fmt.Errorf("driver is zero")
	}
	if err := func() error {
		if a.driver != nil && a.driver.GetDocumentId() == fmt.Sprintf("%d", msg.Driver) {
			return nil
		}
		if a.pidSvc == nil {
			return fmt.Errorf("service pid is nil")
		}
		res, err := ctx.RequestFuture(a.pidSvc, &services.GetCompanyDriverMsg{
			CompanyId: a.companyId,
			DriverDoc: fmt.Sprintf("%d", msg.Driver),
		}, 6*time.Second).Result()
		if err != nil {
			return fmt.Errorf("error request driver: %s", err)
		}

		switch res := res.(type) {
		case *services.CompanyDriverMsg:
			if res.Driver == nil || len(res.Driver.GetDocumentId()) <= 0 {
				return fmt.Errorf("driver not found")
			}
			a.driver = res.Driver
			if a.uix != nil {
				if err := a.uix.Driver(fmt.Sprintf("%d", msg.Driver)); err != nil {
					return fmt.Errorf("driver error: %s", err)
				}
			}
		default:
			return fmt.Errorf("error response: %s (%T)", res, res)
		}
		return nil
	}(); err != nil {
		fmt.Printf("error driver: %s\n", err)
		if a.uix != nil {
			a.uix.Beep(5, 90, 300*time.Millisecond)
			if a.cancelPop != nil {
				a.cancelPop()
			}
			if err := a.uix.TextWarningPopup("conductor no encontrado\n"); err != nil {
				return fmt.Errorf("driver error: %s", err)
			}
			contxt, cancel := context.WithCancel(context.Background())
			a.cancelPop = cancel
			go func() {
				defer cancel()
				select {
				case <-contxt.Done():
				case <-time.After(4 * time.Second):
				}
				if err := a.uix.TextWarningPopupClose(); err != nil {
					fmt.Printf("textWarningPopupClose error: %s", err)
				}
			}()
		}
		return err
	}
	fmt.Printf("driver: %v\n", a.driver)

	return nil
}
