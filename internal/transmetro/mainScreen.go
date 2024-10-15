package app

import "fmt"

func (a *App) mainScreen() error {
	if err := a.uix.MainScreen(); err != nil {
		return fmt.Errorf("main screen error: %s", err)
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
	return nil
}
