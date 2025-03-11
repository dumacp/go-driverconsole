package app

import (
	"fmt"

	"github.com/dumacp/go-schservices/api/services"
)

func (a *App) mainScreen() error {
	if err := a.uix.MainScreen(); err != nil {
		return fmt.Errorf("main screen error: %s", err)
	}

	if a.isItineraryProgEnable {
		if err := a.uix.SetLed(AddrSwitchLang1, true); err != nil {
			return fmt.Errorf("setLed error: %s", err)
		}
	} else {
		if err := a.uix.SetLed(AddrSwitchLang1, false); err != nil {
			return fmt.Errorf("setLed error: %s", err)
		}
	}
	if a.hasCashInput {
		if err := a.uix.SetLed(AddrShowStep, true); err != nil {
			return fmt.Errorf("setLed error: %s", err)
		}
		if err := a.uix.WriteTextRawDisplay(AddrTextCashInputs, []string{"Pagos", "      Conductor"}); err != nil {
			return fmt.Errorf("writeText cash error: %s", err)
		}
	} else {
		if err := a.uix.SetLed(AddrShowStep, false); err != nil {
			return fmt.Errorf("setLed error: %s", err)
		}
		if err := a.uix.WriteTextRawDisplay(AddrTextCashInputs, []string{"Contador", "      Pasajeros"}); err != nil {
			return fmt.Errorf("writeText cash error: %s", err)
		}
	}

	// if err := a.uix.ElectronicInputs(int32(a.electInput)); err != nil {
	// 	return fmt.Errorf("electInput error: %s", err)
	// }
	if err := a.uix.CashInputs(int32(a.cashInput + a.electInput)); err != nil {
		return fmt.Errorf("cashInput error: %s", err)
	}
	if err := a.uix.DateWithFormat(a.updateTime, "2006/01/02 15:04"); err != nil {
		return fmt.Errorf("date error: %s", err)
	}
	if a.driver != nil {
		if err := a.uix.Driver(a.driver.DocumentId); err != nil {
			return fmt.Errorf("driver error: %s", err)
		}
	} else {
		fmt.Printf("driver is nil\n")
		if err := a.uix.Driver(" "); err != nil {
			return fmt.Errorf("driver error: %s", err)
		}
	}
	if err := a.uix.Gps(a.gps); err != nil {
		return fmt.Errorf("gps error: %s", err)
	}
	if err := a.uix.Network(a.network); err != nil {
		return fmt.Errorf("network error: %s", err)
	}

	if len(a.routeString) > 0 {
		// routeS := fmt.Sprintf("%d", a.route)
		if err := a.uix.Route(a.routeString); err != nil {
			return fmt.Errorf("route error: %s", err)
		}
	}
	if a.currentService != nil && a.currentService.GetState() != services.State_ENDED.String() {
		a.showCurrentService(a.currentService)
	}
	return nil
}
