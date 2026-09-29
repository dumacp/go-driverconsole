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

const rowsPerPage = 6

// fetchDriverDailyServices sends a request to the service actor, stores the data and displays page 0.
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
		// Store full list and reset page
		a.dailyServices = result.Services
		a.dailyServicesPage = 0
		return result, nil
	default:
		return nil, fmt.Errorf("unexpected response type: %T", res)
	}
}

// showDailyServices displays page 0 of stored daily services.
func (a *App) showDailyServices(data *platform.DriverDailyServicesResponse) {
	if a.uix == nil {
		logs.LogWarn.Println("showDailyServices: ui is nil")
		return
	}
	a.dailyServices = data.Services
	a.dailyServicesPage = 0
	a.renderDailyPage()
}

// renderDailyPage writes the current page of services to the HMI.
func (a *App) renderDailyPage() {
	addrs := []int{
		ui.SERVICE_DAILY_ROW_0,
		ui.SERVICE_DAILY_ROW_1,
		ui.SERVICE_DAILY_ROW_2,
		ui.SERVICE_DAILY_ROW_3,
		ui.SERVICE_DAILY_ROW_4,
		ui.SERVICE_DAILY_ROW_5,
	}

	// Clear all rows
	for _, addr := range addrs {
		if err := a.uix.WriteTextRawDisplay(addr, []string{""}); err != nil {
			logs.LogWarn.Printf("daily clear row %d error: %s", addr, err)
		}
	}

	// Calculate slice for current page
	start := a.dailyServicesPage * rowsPerPage
	end := start + rowsPerPage
	if end > len(a.dailyServices) {
		end = len(a.dailyServices)
	}
	page := a.dailyServices[start:end]

	// Write rows
	for i, svc := range page {
		line := formatServiceRow(start+i+1, &svc)
		if err := a.uix.WriteTextRawDisplay(addrs[i], []string{line}); err != nil {
			logs.LogWarn.Printf("daily write row %d error: %s", addrs[i], err)
		}
	}

	logs.LogInfo.Printf("daily services page %d: %d services (total %d)",
		a.dailyServicesPage, len(page), len(a.dailyServices))
}

// nextDailyPage advances to the next page of services and re-renders.
func (a *App) nextDailyPage() {
	if len(a.dailyServices) == 0 {
		logs.LogWarn.Println("nextDailyPage: no services stored")
		return
	}
	totalPages := (len(a.dailyServices) + rowsPerPage - 1) / rowsPerPage
	a.dailyServicesPage = (a.dailyServicesPage + 1) % totalPages
	a.renderDailyPage()
}

// formatServiceRow formats a single service into a fixed-width display string.
func formatServiceRow(num int, svc *platform.DriverDailyService) string {
	t := time.UnixMilli(svc.StartTime)
	timeStr := t.Format("15:04")

	routeName := svc.RouteName
	maxRouteLen := 32
	if len(routeName) > maxRouteLen {
		routeName = routeName[:maxRouteLen]
	}

	return fmt.Sprintf(" %d %-5s %-32s D:%d/%d T:%d/%d",
		num, timeStr, routeName,
		svc.FrontDoor.Inputs, svc.FrontDoor.Outputs,
		svc.BackDoor.Inputs, svc.BackDoor.Outputs,
	)
}
