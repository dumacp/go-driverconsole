package ui

import "github.com/dumacp/go-driverconsole/internal/display"

const (
	ROUTE_TEXT int = iota
	ROUTE_TEXT_READ
	DRIVER_TEXT
	DRIVER_TEXT_READ
	INPUTS_TEXT
	OUTPUTS_TEXT
	DEVIATION_TEXT
	DATE_TEXT
	WARNING_TEXT
	CONFIRMATION_TEXT
	POPUP_TEXT
	POPUP_WARN_TEXT
	PROGRAMATION_VEH_TEXT
	PROGRAMATION_DRIVER_TEXT
	DOOR_0_LED
	DOOR_1_LED
	DOOR_2_LED
	GPS_LED
	NETWORK_LED
	NOTIFICATIONS_BUTTON
	PROGRAMATION_VEH_BUTTON
	PROGRAMATION_DRIVER_BUTTON
	ADDITIONALS_BUTTON
	SERVICE_CURRENT_STATE
	SERVICE_CURRENT_STATE_TEXT
)

func Label2DisplayRegister(label int) display.Register {

	switch label {
	case ROUTE_TEXT:
		return display.Register{
			Type:   display.INPUT_TEXT,
			Addr:   AddrTextRoute,
			Len:    1,
			Size:   32,
			Gap:    0,
			Toogle: 0,
		}
	case DRIVER_TEXT:
		return display.Register{
			Type:   display.INPUT_TEXT,
			Addr:   AddrTextRoute,
			Len:    1,
			Size:   20,
			Gap:    0,
			Toogle: 0,
		}
	case DATE_TEXT:
		return display.Register{
			Type:   display.INPUT_TEXT,
			Addr:   AddrTextDriver,
			Len:    1,
			Size:   32,
			Gap:    0,
			Toogle: 0,
		}
	case CONFIRMATION_TEXT:
		return display.Register{
			Type:   display.INPUT_TEXT,
			Addr:   AddrTextConfirmation,
			Len:    1,
			Size:   60,
			Gap:    0,
			Toogle: AddrToggleConfirmation,
		}
	case POPUP_TEXT:
		return display.Register{
			Type:   display.INPUT_NUM,
			Addr:   AddrTextPopup,
			Len:    1,
			Size:   60,
			Gap:    0,
			Toogle: AddrTogglePopup,
		}
	case INPUTS_TEXT:
		return display.Register{
			Type:   display.INPUT_NUM,
			Addr:   AddrNumInputs,
			Len:    1,
			Size:   2,
			Gap:    0,
			Toogle: 0,
		}
	case OUTPUTS_TEXT:
		return display.Register{
			Type:   display.INPUT_NUM,
			Addr:   AddrNumOutputs,
			Len:    1,
			Size:   2,
			Gap:    0,
			Toogle: 0,
		}
	case DEVIATION_TEXT:
		return display.Register{
			Type:   display.INPUT_NUM,
			Addr:   AddrNumDeviation,
			Len:    1,
			Size:   2,
			Gap:    0,
			Toogle: 0,
		}
	case SERVICE_CURRENT_STATE_TEXT:
		return display.Register{
			Type:   display.INPUT_TEXT,
			Addr:   AddrTextCurrentService,
			Len:    1,
			Size:   30,
			Gap:    0,
			Toogle: 0,
		}
	case SERVICE_CURRENT_STATE:
		return display.Register{
			Type:   display.ARRAY_PICT,
			Addr:   AddrLedCurrentService,
			Len:    1,
			Size:   2,
			Gap:    0,
			Toogle: 0,
		}

	}
	return display.Register{}
}
