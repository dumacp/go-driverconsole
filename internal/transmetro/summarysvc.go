package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/dumacp/go-driverconsole/internal/ui"
	"github.com/dumacp/go-logs/pkg/logs"
	"github.com/dumacp/go-schservices/api/services"
)

func (a *App) summaryservice() error {
	if a.currentService == nil {
		return fmt.Errorf(`no hay un servicio iniciado`)
	}

	if time.Since(a.lastReqSummarySvc) < 1*time.Minute && a.summaryService != nil &&
		a.summaryService.GetXId() == a.currentService.Id {
		return nil
	}

	if a.pidSvc == nil {
		return fmt.Errorf("service pid is nil")
	}

	funcRequest := func(mss interface{}) error {
		res, err := a.ctx.RequestFuture(a.pidSvc, mss, 10*time.Second).Result()
		if err != nil {
			return fmt.Errorf("request retake service error: %s", err)
		}

		if resSvc, ok := res.(*services.ServiceSummaryMsg); ok {

			if len(resSvc.GetError()) > 0 {
				logs.LogWarn.Printf("error request retake service: %s", err)
				return fmt.Errorf(resSvc.GetError())
			}
			a.summaryService = resSvc.GetSummary()
		} else {
			a.summaryService = nil
			return fmt.Errorf("response type error (%T)", res)
		}

		a.lastReqSummarySvc = time.Now()
		return nil
	}

	switch {
	case a.currentService == nil:
		return fmt.Errorf("no hay un servicio iniciado")
	case a.currentService.State == services.State_READY_TO_START.String() ||
		a.currentService.State == services.State_WAITING_TO_ARRIVE_TO_STARTING_POINT.String() ||
		a.currentService.State == services.State_STARTED.String() ||
		a.currentService.State == services.State_SCHEDULED.String() ||
		a.currentService.State == services.State_ENDED.String():
		mss := &services.GetServiceSummaryMsg{
			ServiceId: a.currentService.Id,
		}
		if err := funcRequest(mss); err != nil {
			return err
		}
	default:
		fmt.Printf("current service: %v\n", a.currentService)
		return fmt.Errorf("servicio en estado inválido")
	}
	if a.summaryService != nil {
		if a.summaryService.GetPrevVehicle() != nil && len(a.summaryService.GetPrevVehicle().GetPlate()) > 0 {
			if err := a.uix.WriteTextRawDisplay(AddrPrevVehHeaderText, []string{
				fmt.Sprintf("%s (%s)", a.summaryService.GetPrevVehicle().GetPlate(), a.summaryService.GetPrevVehicle().GetNumber())}); err != nil {
				return fmt.Errorf("error prev header: %s", err)
			}
			if len(a.summaryService.GetPrevVehicle().GetCheckpoint()) > 0 {
				if err := a.uix.WriteTextRawDisplay(AddrPrevVehFooterText, []string{
					fmt.Sprintf("%s\nRegulación: %d", a.summaryService.GetPrevVehicle().GetCheckpoint(), a.summaryService.GetPrevVehicle().GetTimeDiff()),
				}); err != nil {
					return fmt.Errorf("error prev footer: %s", err)
				}
			}
		} else {
			sizeHeader := Label2DisplayRegister(AddrPrevVehHeaderText).Size
			spacesHeader := strings.Repeat(" ", sizeHeader)
			if err := a.uix.WriteTextRawDisplay(AddrPrevVehHeaderText, []string{
				spacesHeader,
			}); err != nil {
				return fmt.Errorf("error prev header: %s", err)
			}
			sizeFooter := Label2DisplayRegister(AddrPrevVehFooterText).Size
			spacesFooter := strings.Repeat(" ", sizeFooter)
			if err := a.uix.WriteTextRawDisplay(AddrPrevVehFooterText, []string{
				spacesFooter,
			}); err != nil {
				return fmt.Errorf("error prev footer: %s", err)
			}
		}
		if a.summaryService.GetNextVehicle() != nil && len(a.summaryService.GetNextVehicle().GetPlate()) > 0 {
			if err := a.uix.WriteTextRawDisplay(AddrNextVehHeaderText, []string{
				fmt.Sprintf("%s (%s)", a.summaryService.GetNextVehicle().GetPlate(), a.summaryService.GetNextVehicle().GetNumber())}); err != nil {

				return fmt.Errorf("error next header: %s", err)
			}
			if len(a.summaryService.GetNextVehicle().GetCheckpoint()) > 0 {
				if err := a.uix.WriteTextRawDisplay(AddrNextVehFooterText, []string{
					fmt.Sprintf("%s\nRegulación: %d", a.summaryService.GetNextVehicle().GetCheckpoint(), a.summaryService.GetNextVehicle().GetTimeDiff()),
				}); err != nil {
					return fmt.Errorf("error next footer: %s", err)
				}
			}
		} else {
			sizeHeader := Label2DisplayRegister(AddrNextVehHeaderText).Size
			spacesHeader := strings.Repeat(" ", sizeHeader)
			if err := a.uix.WriteTextRawDisplay(AddrNextVehHeaderText, []string{
				spacesHeader,
			}); err != nil {
				return fmt.Errorf("error next header: %s", err)
			}
			sizeFooter := Label2DisplayRegister(AddrNextVehFooterText).Size
			spacesFooter := strings.Repeat(" ", sizeFooter)
			if err := a.uix.WriteTextRawDisplay(AddrNextVehFooterText, []string{
				spacesFooter,
			}); err != nil {
				return fmt.Errorf("error next footer: %s", err)
			}
		}
		if a.summaryService.GetVehicle() != nil && len(a.summaryService.GetVehicle().GetPlate()) > 0 {
			state := int(services.TimingState_value[a.summaryService.GetVehicle().GetTimingState()])
			if err := a.uix.ArrayPict(ui.SERVICE_SUMMARY_STATE, state); err != nil {
				return fmt.Errorf("error curr pict: %s", err)
			}
			if err := a.uix.WriteTextRawDisplay(AddrCurrVehHeaderText, []string{
				fmt.Sprintf("%s (%s)", a.summaryService.GetVehicle().GetPlate(), a.summaryService.GetVehicle().GetNumber())}); err != nil {

				return fmt.Errorf("error curr header: %s", err)
			}
			if len(a.summaryService.GetVehicle().GetCheckpoint()) > 0 {
				if err := a.uix.WriteTextRawDisplay(AddrCurrVehFooterText, []string{
					fmt.Sprintf("%s\nRegulación: %d", a.summaryService.GetVehicle().GetCheckpoint(), a.summaryService.GetVehicle().GetTimeDiff()),
				}); err != nil {
					return fmt.Errorf("error curr footer: %s", err)
				}
			}
			if len(a.summaryService.GetVehicle().GetCheckpoint()) > 0 {
				if err := a.uix.WriteTextRawDisplay(AddrCurrCheckpointText, []string{
					a.summaryService.GetVehicle().GetCheckpoint(),
				}); err != nil {
					return fmt.Errorf("error curr footer: %s", err)
				}
			}
		} else {
			state := int(services.TimingState_NOT_APPLICABLE)
			if err := a.uix.ArrayPict(ui.SERVICE_SUMMARY_STATE, state); err != nil {
				return fmt.Errorf("error curr pict: %s", err)
			}
			if err := a.uix.WriteTextRawDisplay(AddrCurrVehHeaderText, []string{
				"Sin vehículo",
			}); err != nil {
				return fmt.Errorf("error curr header: %s", err)
			}
			size := Label2DisplayRegister(AddrCurrVehHeaderText).Size
			spaces := strings.Repeat(" ", size)
			if err := a.uix.WriteTextRawDisplay(AddrCurrVehFooterText, []string{
				spaces,
			}); err != nil {
				return fmt.Errorf("error curr footer: %s", err)
			}
			sizeCheckpoint := Label2DisplayRegister(AddrCurrCheckpointText).Size
			spacesCheckpoint := strings.Repeat(" ", sizeCheckpoint)
			if err := a.uix.WriteTextRawDisplay(AddrCurrCheckpointText, []string{
				spacesCheckpoint,
			}); err != nil {
				return fmt.Errorf("error curr checkpoint: %s", err)
			}
		}
		if err := a.uix.WriteTextRawDisplay(AddrCurrItineraryText, []string{
			fmt.Sprintf("%s (%d)", a.summaryService.GetItinerary().GetItineraryName(), a.summaryService.GetItinerary().GetItineraryPmc()),
		}); err != nil {
			return fmt.Errorf("error curr itinerary: %s", err)
		}
	}

	return nil
}
