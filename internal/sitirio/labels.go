package app

import (
	"github.com/dumacp/go-driverconsole/internal/display"
	"github.com/dumacp/go-driverconsole/internal/ui"
)

func Label2DisplayRegister(label int) display.Register {

	switch label {
	case ui.ROUTE_TEXT:
		return display.Register{
			Type:   display.INPUT_TEXT,
			Addr:   AddrTextRoute,
			Len:    1,
			Size:   32,
			Gap:    0,
			Toogle: 0,
		}
	case ui.DRIVER_TEXT:
		return display.Register{
			Type:   display.INPUT_TEXT,
			Addr:   AddrTextDriver,
			Len:    1,
			Size:   20,
			Gap:    0,
			Toogle: 0,
		}

	case ui.ROUTE_TEXT_READ:
		return display.Register{
			Type:   display.INPUT_TEXT,
			Addr:   AddrNumRoute,
			Len:    1,
			Size:   8,
			Gap:    0,
			Toogle: 0,
		}
	case ui.DRIVER_TEXT_READ:
		return display.Register{
			Type:   display.INPUT_TEXT,
			Addr:   AddrNumDriver,
			Len:    1,
			Size:   8,
			Gap:    0,
			Toogle: 0,
		}
	case ui.DATE_TEXT:
		return display.Register{
			Type:   display.INPUT_TEXT,
			Addr:   AddrTextDate,
			Len:    1,
			Size:   32,
			Gap:    0,
			Toogle: 0,
		}
	case ui.CONFIRMATION_TEXT:
		return display.Register{
			Type:   display.INPUT_TEXT,
			Addr:   AddrTextConfirmation,
			Len:    1,
			Size:   60,
			Gap:    0,
			Toogle: 0,
		}
	case ui.WARNING_TEXT:
		return display.Register{
			Type:   display.INPUT_TEXT,
			Addr:   AddrTextWarning,
			Len:    1,
			Size:   60,
			Gap:    0,
			Toogle: 0,
		}
	case ui.POPUP_TEXT:
		return display.Register{
			Type:   display.INPUT_TEXT,
			Addr:   AddrTextPopup,
			Len:    1,
			Size:   60,
			Gap:    0,
			Toogle: AddrTogglePopup,
		}
	case ui.POPUP_WARN_TEXT:
		return display.Register{
			Type:   display.INPUT_TEXT,
			Addr:   AddrTextWarnPopup,
			Len:    1,
			Size:   60,
			Gap:    0,
			Toogle: AddrToggleWarnPopup,
		}
	case ui.CASH_INPUTS_TEXT:
		return display.Register{
			Type:   display.INPUT_NUM,
			Addr:   AddrNumCashInputs,
			Len:    1,
			Size:   4,
			Gap:    0,
			Toogle: 0,
		}
	case ui.ELECT_INPUTS_TEXT:
		return display.Register{
			Type:   display.INPUT_NUM,
			Addr:   AddrNumElectonicInputs,
			Len:    1,
			Size:   4,
			Gap:    0,
			Toogle: 0,
		}
	case ui.DEVIATION_TEXT:
		return display.Register{
			Type:   display.INPUT_NUM,
			Addr:   AddrNumDeviation,
			Len:    1,
			Size:   4,
			Gap:    0,
			Toogle: 0,
		}
	case ui.SERVICE_CURRENT_STATE_TEXT:
		return display.Register{
			Type:   display.INPUT_TEXT,
			Addr:   AddrTextCurrentService,
			Len:    1,
			Size:   30,
			Gap:    0,
			Toogle: 0,
		}
	case ui.SERVICE_CURRENT_STATE:
		return display.Register{
			Type:   display.ARRAY_PICT,
			Addr:   AddrLedCurrentService,
			Len:    1,
			Size:   2,
			Gap:    0,
			Toogle: 0,
		}
	case ui.PROGRAMATION_VEH_TEXT:
		return display.Register{
			Type:   display.INPUT_TEXT,
			Addr:   AddrTextProgVeh,
			Len:    10,
			Size:   100,
			Gap:    100,
			Toogle: 0,
		}
	case ui.PROGRAMATION_DRIVER_TEXT:
		return display.Register{
			Type:   display.INPUT_TEXT,
			Addr:   AddrTextProgDriver,
			Len:    10,
			Size:   100,
			Gap:    100,
			Toogle: 0,
		}
	case ui.NOTIFICATIONS_ALARM_TEXT:
		return display.Register{
			Type:   display.INPUT_TEXT,
			Addr:   AddrTextNotiAlarm,
			Len:    10,
			Size:   100,
			Gap:    100,
			Toogle: 0,
		}
	case ui.RESET:
		return display.Register{
			Type:   display.LED,
			Addr:   AddrLedReset,
			Len:    0,
			Size:   0,
			Gap:    0,
			Toogle: 0,
		}
	case ui.GPS_LED:
		return display.Register{
			Type:   display.LED,
			Addr:   AddrLedGsp,
			Len:    0,
			Size:   0,
			Gap:    0,
			Toogle: 0,
		}
	case ui.NETWORK_LED:
		return display.Register{
			Type:   display.LED,
			Addr:   AddrLedNetwork,
			Len:    0,
			Size:   0,
			Gap:    0,
			Toogle: 0,
		}
	case ui.STEP_ENABLE:
		return display.Register{
			Type:   display.LED,
			Addr:   AddrSwitchStep,
			Len:    0,
			Size:   0,
			Gap:    0,
			Toogle: 0,
		}
	case AddrEnterRuta:
		return display.Register{
			Type:   display.LED,
			Addr:   AddrEnterRuta,
			Len:    0,
			Size:   0,
			Gap:    0,
			Toogle: 0,
		}
	case AddrEnterDriver:
		return display.Register{
			Type:   display.LED,
			Addr:   AddrEnterDriver,
			Len:    0,
			Size:   0,
			Gap:    0,
			Toogle: 0,
		}
	}
	return display.Register{}
}

func Label2DisplayRegisterGtt(label int) display.Register {

	switch label {
	case ui.ROUTE_TEXT:
		return display.Register{
			Type:   display.INPUT_TEXT,
			Addr:   AddrGttTextRoute,
			Len:    1,
			Size:   32,
			Gap:    0,
			Toogle: 0,
		}
	case ui.DRIVER_TEXT:
		return display.Register{
			Type:   display.INPUT_TEXT,
			Addr:   AddrGttTextDriver,
			Len:    1,
			Size:   32,
			Gap:    0,
			Toogle: 0,
		}

	case ui.ROUTE_TEXT_READ:
		return display.Register{
			Type:   display.INPUT_TEXT,
			Addr:   AddrGttNumRoute,
			Len:    1,
			Size:   8,
			Gap:    0,
			Toogle: 0,
		}
	case ui.DRIVER_TEXT_READ:
		return display.Register{
			Type:   display.INPUT_TEXT,
			Addr:   AddrGttNumDriver,
			Len:    1,
			Size:   8,
			Gap:    0,
			Toogle: 0,
		}
	case ui.DATE_TEXT:
		return display.Register{
			Type:   display.INPUT_TEXT,
			Addr:   AddrGttTextDate,
			Len:    1,
			Size:   32,
			Gap:    0,
			Toogle: 0,
		}
	case ui.CONFIRMATION_TEXT:
		return display.Register{
			Type:   display.INPUT_TEXT,
			Addr:   AddrGttTextConfirmation,
			Len:    1,
			Size:   180,
			Gap:    0,
			Toogle: 0,
		}
	case ui.WARNING_TEXT:
		return display.Register{
			Type:   display.INPUT_TEXT,
			Addr:   AddrGttTextWarning,
			Len:    1,
			Size:   180,
			Gap:    0,
			Toogle: 0,
		}
	case ui.POPUP_TEXT:
		return display.Register{
			Type:   display.INPUT_TEXT,
			Addr:   AddrGttTextPopup,
			Len:    1,
			Size:   100,
			Gap:    0,
			Toogle: AddrGttTogglePopup,
		}
	case ui.POPUP_WARN_TEXT:
		return display.Register{
			Type:   display.INPUT_TEXT,
			Addr:   AddrGttTextWarnPopup,
			Len:    1,
			Size:   100,
			Gap:    0,
			Toogle: AddrGttToggleWarnPopup,
		}
	case ui.CASH_INPUTS_TEXT:
		return display.Register{
			Type:   display.INPUT_NUM,
			Addr:   AddrGttNumCashInputs,
			Len:    1,
			Size:   4,
			Gap:    0,
			Toogle: 0,
		}
	case ui.ELECT_INPUTS_TEXT:
		return display.Register{
			Type:   display.INPUT_NUM,
			Addr:   AddrGttNumElectonicInputs,
			Len:    1,
			Size:   4,
			Gap:    0,
			Toogle: 0,
		}
	case ui.INPUTS_TEXT:
		return display.Register{
			Type:   display.INPUT_NUM,
			Addr:   AddrGttCounterInputs,
			Len:    1,
			Size:   4,
			Gap:    0,
			Toogle: 0,
		}
	case ui.OUTPUTS_TEXT:
		return display.Register{
			Type:   display.INPUT_NUM,
			Addr:   AddrGttCounterOutputs,
			Len:    1,
			Size:   4,
			Gap:    0,
			Toogle: 0,
		}
	case ui.DEVIATION_TEXT:
		return display.Register{
			Type:   display.INPUT_NUM,
			Addr:   AddrGttNumDeviation,
			Len:    1,
			Size:   4,
			Gap:    0,
			Toogle: 0,
		}
	case ui.SERVICE_CURRENT_STATE_TEXT:
		return display.Register{
			Type:   display.INPUT_TEXT,
			Addr:   AddrGttTextCurrentService,
			Len:    1,
			Size:   30,
			Gap:    0,
			Toogle: 0,
		}
	case ui.SERVICE_CURRENT_STATE:
		return display.Register{
			Type:   display.ARRAY_PICT,
			Addr:   AddrGttLedCurrentService,
			Len:    1,
			Size:   2,
			Gap:    0,
			Toogle: 0,
		}
	case ui.PROGRAMATION_VEH_TEXT:
		return display.Register{
			Type:   display.INPUT_TEXT,
			Addr:   AddrGttTextProgVeh,
			Len:    10,
			Size:   100,
			Gap:    0,
			Toogle: 0,
		}
	case ui.PROGRAMATION_DRIVER_TEXT:
		return display.Register{
			Type:   display.INPUT_TEXT,
			Addr:   AddrGttTextProgDriver,
			Len:    10,
			Size:   100,
			Gap:    0,
			Toogle: 0,
		}
	case ui.NOTIFICATIONS_ALARM_TEXT:
		return display.Register{
			Type:   display.INPUT_TEXT,
			Addr:   AddrGttTextNotiAlarm,
			Len:    10,
			Size:   100,
			Gap:    0,
			Toogle: 0,
		}
	case ui.RESET:
		return display.Register{
			Type:   display.LED,
			Addr:   AddrGttLedReset,
			Len:    0,
			Size:   0,
			Gap:    0,
			Toogle: 0,
		}
	case ui.GPS_LED:
		return display.Register{
			Type:   display.LED,
			Addr:   AddrGttLedGsp,
			Len:    0,
			Size:   0,
			Gap:    0,
			Toogle: 0,
		}
	case ui.NETWORK_LED:
		return display.Register{
			Type:   display.LED,
			Addr:   AddrGttLedNetwork,
			Len:    0,
			Size:   0,
			Gap:    0,
			Toogle: 0,
		}
	case ui.STEP_ENABLE:
		return display.Register{
			Type:   display.BUTTON,
			Addr:   AddrGttSwitchStep,
			Len:    0,
			Size:   0,
			Gap:    0,
			Toogle: 0,
		}
	case ui.MAIN_SCREEN:
		return display.Register{
			Type:   display.SCREEN_NUM,
			Addr:   AddrGttScreenMain,
			Len:    0,
			Size:   0,
			Gap:    0,
			Toogle: 0,
		}
	case ui.KEY_ROUTE_SCREEN:
		return display.Register{
			Type:   display.SCREEN_NUM,
			Addr:   AddrGttScreenKeyRoute,
			Len:    0,
			Size:   0,
			Gap:    0,
			Toogle: 0,
		}
	case int(ui.KEY_DRIVER_SCREEN):
		return display.Register{
			Type:   display.SCREEN_NUM,
			Addr:   AddrGttScreenKeyDriver,
			Len:    0,
			Size:   0,
			Gap:    0,
			Toogle: 0,
		}
	case ui.INFO_SCREEN:
		return display.Register{
			Type:   display.SCREEN_NUM,
			Addr:   AddrGttScreenInfo,
			Len:    0,
			Size:   0,
			Gap:    0,
			Toogle: 0,
		}
	case ui.WARN_SCREEN:
		return display.Register{
			Type:   display.SCREEN_NUM,
			Addr:   AddrGttScreenWarn,
			Len:    0,
			Size:   0,
			Gap:    0,
			Toogle: 0,
		}
	case ui.PROGRAMATION_DRIVER_SCREEN:
		return display.Register{
			Type:   display.SCREEN_NUM,
			Addr:   AddrGttScreenDriver,
			Len:    0,
			Size:   0,
			Gap:    0,
			Toogle: 0,
		}
	case ui.PROGRAMATION_VEH_SCREEN:
		return display.Register{
			Type:   display.SCREEN_NUM,
			Addr:   AddrGttScreenVeh,
			Len:    0,
			Size:   0,
			Gap:    0,
			Toogle: 0,
		}
	case ui.ALARMS_SCREEN:
		return display.Register{
			Type:   display.SCREEN_NUM,
			Addr:   AddrGttScreenAlarm,
			Len:    0,
			Size:   0,
			Gap:    0,
			Toogle: 0,
		}

	}
	return display.Register{}
}
