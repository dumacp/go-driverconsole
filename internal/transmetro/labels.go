package app

import (
	"github.com/dumacp/go-driverconsole/internal/display"
	"github.com/dumacp/go-driverconsole/internal/ui"
)

// const (
// 	addrConfirmation    int = 600
// 	addrError           int = 500
// 	addrAlarms          int = 3000
// 	addrEfectivoCounter int = 90
// 	addrTagCounter      int = 88

// 	addrInputs  int = 80
// 	addrOutputs int = 84

// 	addrFrontalDoor int = 10
// 	addrBackDoor    int = 11

// 	addrNoRoute    = 120
// 	addrNameRoute  = 100
// 	addrNoDriver   = 160
// 	addrNameDriver = 140

// 	addrConfirmationTextMainScreen      int = 400
// 	addrConfirmationToggleMainScreen    int = 5
// 	addrConfirmationTextMainScreenErr   int = 300
// 	addrConfirmationToggleMainScreenErr int = 6

// 	addrIconGPS int = 12
// 	addrIconNET int = 11

// 	addrTimeDate int = 60
// )

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
	case ui.INPUTS_TEXT:
		return display.Register{
			Type:   display.INPUT_NUM,
			Addr:   AddrNumElectonicInputs,
			Len:    1,
			Size:   4,
			Gap:    0,
			Toogle: 0,
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
	case AddrTextCurrentItinerary:
		return display.Register{
			Type:   display.INPUT_TEXT,
			Addr:   AddrTextCurrentItinerary,
			Len:    1,
			Size:   90,
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
			Size:   50,
			Gap:    50,
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
	// case ui.STEP_ENABLE:
	// 	return display.Register{
	// 		Type:   display.LED,
	// 		Addr:   AddrSwitchStep,
	// 		Len:    0,
	// 		Size:   0,
	// 		Gap:    0,
	// 		Toogle: 0,
	// 	}
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
	case AddrScreenSwitch:
		return display.Register{
			Type:   display.LED,
			Addr:   AddrScreenSwitch,
			Len:    0,
			Size:   0,
			Gap:    0,
			Toogle: 0,
		}
	case AddrScreenProgDriver:
		return display.Register{
			Type:   display.LED,
			Addr:   AddrScreenProgDriver,
			Len:    0,
			Size:   0,
			Gap:    0,
			Toogle: 0,
		}
	case AddrScreenProgVeh:
		return display.Register{
			Type:   display.LED,
			Addr:   AddrScreenProgVeh,
			Len:    0,
			Size:   0,
			Gap:    0,
			Toogle: 0,
		}
	case AddrScreenAlarms:
		return display.Register{
			Type:   display.LED,
			Addr:   AddrScreenAlarms,
			Len:    0,
			Size:   0,
			Gap:    0,
			Toogle: 0,
		}
	case AddrScreenMore:
		return display.Register{
			Type:   display.LED,
			Addr:   AddrScreenMore,
			Len:    0,
			Size:   0,
			Gap:    0,
			Toogle: 0,
		}
	case AddrUpdateDropListProgVeh:
		return display.Register{
			Type:   display.LED,
			Addr:   AddrUpdateDropListProgVeh,
			Len:    0,
			Size:   0,
			Gap:    0,
			Toogle: 0,
		}
	case AddrCurrentSelectProgVeh:
		return display.Register{
			Type:   display.INPUT_NUM,
			Addr:   AddrCurrentSelectProgVeh,
			Len:    1,
			Size:   1,
			Gap:    0,
			Toogle: 0,
		}
	case AddrResumeSelectProgVeh:
		return display.Register{
			Type:   display.INPUT_TEXT,
			Addr:   AddrResumeSelectProgVeh,
			Len:    1,
			Size:   254,
			Gap:    0,
			Toogle: 0,
		}
	case AddrSelectItinerary:
		return display.Register{
			Type:   display.LED,
			Addr:   AddrSelectItinerary,
			Len:    0,
			Size:   0,
			Gap:    0,
			Toogle: 0,
		}
	case AddrItineraryProgVeh:
		return display.Register{
			Type:   display.INPUT_NUM,
			Addr:   AddrItineraryProgVeh,
			Len:    1,
			Size:   2,
			Gap:    0,
			Toogle: 0,
		}
	case AddrItineraryProgVehVerify:
		return display.Register{
			Type:   display.INPUT_NUM,
			Addr:   AddrItineraryProgVehVerify,
			Len:    1,
			Size:   2,
			Gap:    0,
			Toogle: 0,
		}
	}
	return display.Register{}
}
