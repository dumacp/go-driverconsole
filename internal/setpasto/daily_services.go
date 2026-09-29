package app

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dumacp/go-driverconsole/internal/platform"
	"github.com/dumacp/go-driverconsole/internal/ui"
	"github.com/dumacp/go-logs/pkg/logs"
	"github.com/dumacp/go-schservices/api/services"
)

// fetchDriverDailyServices sends a request to the service actor and returns parsed data.
func (a *App) fetchDriverDailyServices() (*platform.DriverDailyServicesResponse, error) {
	if a.driver == nil || len(a.driver.GetDocumentId()) <= 0 {
		return nil, fmt.Errorf("driver document not available")
	}
	if a.pidSvc == nil {
		return nil, fmt.Errorf("service pid is nil")
	}

	mss := &services.GetDriverDailyServicesMsg{
		DeviceId: a.deviceId,
		DriverId: a.driver.GetDocumentId(),
	}

	res, err := a.ctx.RequestFuture(a.pidSvc, mss, 15*time.Second).Result()
	if err != nil {
		return nil, fmt.Errorf("request daily services error: %w", err)
	}

	switch r := res.(type) {
	case *services.DriverDailyServicesResponseMsg:
		if len(r.GetError()) > 0 {
			return nil, fmt.Errorf(r.GetError())
		}
		result := &platform.DriverDailyServicesResponse{}
		if err := json.Unmarshal(r.GetData(), result); err != nil {
			return nil, fmt.Errorf("unmarshal error: %w", err)
		}
		return result, nil
	default:
		return nil, fmt.Errorf("unexpected response type: %T", res)
	}
}

// showDailyServices formats and writes daily services to the HMI display (screen 8).
func (a *App) showDailyServices(data *platform.DriverDailyServicesResponse) {
	if a.uix == nil {
		logs.LogWarn.Println("showDailyServices: ui is nil")
		return
	}

	addrs := []int{
		ui.SERVICE_DAILY_ROW_0,
		ui.SERVICE_DAILY_ROW_1,
		ui.SERVICE_DAILY_ROW_2,
		ui.SERVICE_DAILY_ROW_3,
		ui.SERVICE_DAILY_ROW_4,
		ui.SERVICE_DAILY_ROW_5,
	}

	// Clear all rows first
	for _, addr := range addrs {
		if err := a.uix.WriteTextRawDisplay(addr, []string{""}); err != nil {
			logs.LogWarn.Printf("showDailyServices clear row %d error: %s", addr, err)
		}
	}

	// Write each service row
	for i, svc := range data.Services {
		if i >= len(addrs) {
			break
		}
		line := formatServiceRow(i+1, &svc)
		if err := a.uix.WriteTextRawDisplay(addrs[i], []string{line}); err != nil {
			logs.LogWarn.Printf("showDailyServices write row %d error: %s", addrs[i], err)
		}
	}

	logs.LogInfo.Printf("daily services displayed: %d services", len(data.Services))
}

// formatServiceRow formats a single service into a fixed-width display string.
// Layout: "#  HH:MM  RouteName           FDin FDout BDin BDout"
func formatServiceRow(num int, svc *platform.DriverDailyService) string {
	t := time.UnixMilli(svc.StartTime)
	timeStr := t.Format("15:04")

	routeName := svc.RouteName
	maxRouteLen := 22
	if len(routeName) > maxRouteLen {
		routeName = routeName[:maxRouteLen]
	}

	return fmt.Sprintf("%-2d %-5s %-20s D:%d/%d T:%d/%d",
		num, timeStr, routeName,
		svc.FrontDoor.Inputs, svc.FrontDoor.Outputs,
		svc.BackDoor.Inputs, svc.BackDoor.Outputs,
	)
}
