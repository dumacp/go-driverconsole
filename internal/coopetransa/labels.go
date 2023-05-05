package app

import (
	"github.com/dumacp/go-driverconsole/internal/display"
	"github.com/dumacp/go-driverconsole/internal/ui"
)

const (
	addrConfirmation    int = 600
	addrError           int = 500
	addrAlarms          int = 3000
	addrEfectivoCounter int = 90
	addrTagCounter      int = 88

	addrInputs  int = 80
	addrOutputs int = 84

	addrFrontalDoor int = 10
	addrBackDoor    int = 11

	addrNoRoute    = 120
	addrNameRoute  = 100
	addrNoDriver   = 160
	addrNameDriver = 140

	addrConfirmationTextMainScreen      int = 400
	addrConfirmationToggleMainScreen    int = 5
	addrConfirmationTextMainScreenErr   int = 300
	addrConfirmationToggleMainScreenErr int = 6

	addrIconGPS int = 12
	addrIconNET int = 11

	addrTimeDate int = 60
)

func label2addr(label any) display.Register {
	switch label {
	case ui.DRIVER_TEXT:
		return display.Register{}
	}
	return display.Register{}
}
