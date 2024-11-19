package app

import (
	"time"

	"github.com/dumacp/go-driverconsole/internal/itinerary"
	"github.com/dumacp/matrixorbital/gtt43a"
)

type DisplayDevice struct {
	Device gtt43a.Display
}
type InputText struct {
	Text string
}

type SelectPaso struct{}
type StepMsg struct{}
type TestStepMsg struct{}
type ResetCounter struct{}
type ErrorDisplay struct {
	Error error
}

type MsgAppPaso struct {
	Value int
}
type MsgAppPercentRecorrido struct {
	Data int
}
type MsgMainScreen struct{}
type MsgScreen struct {
	ID     int
	Switch bool
}
type MsgInitRecorrido struct{}
type MsgStopRecorrido struct{}
type ResetRecorrido struct{}
type MsgSubscribe struct{}
type MsgChangeRoute struct {
	ID int
}
type MsgChangeDriver struct {
	ID int
}

type MsgCounters struct {
	Parcial  int
	Efectivo int
	App      int
}

type MsgRoute struct {
	Route string
}
type MsgSetRoute struct {
	Route     int
	RouteName string
}
type MsgDriver struct {
	Driver int
}
type MsgSetDriver struct {
	Driver int
}

type MsgSetRoutes struct {
	Routes map[int32]string
}

type MsgDoors struct {
	Value [2]int
}

type MsgConfirmationText struct {
	Text []byte
}
type MsgConfirmationTextMainScreen struct {
	Text    []byte
	Timeout time.Duration
}
type MsgWarningText struct {
	Text []byte
}
type MsgWarningTextInMainScreen struct {
	Text    []byte
	Timeout time.Duration
}
type MsgGpsDown struct{}
type MsgGpsUP struct{}
type MsgNetDown struct{}
type MsgNetUP struct{}

type MsgShowAlarms struct{}
type MsgAddAlarm struct {
	Data string
}

type MsgBrightnessPlus struct{}
type MsgBrightnessMinus struct{}

type MsgGetItinieary struct {
	ID             string
	OrganizationID string
	PaymentID      int
}
type MsgItinirary struct {
	Data itinerary.Itinerary
}

type MsgUpdateTime struct{}

type MsgShowCounters struct{}

type TestTextProgDriver struct {
	Text []string
}
type ListProgDriver struct {
	Itinerary      int
	DriverDocument string
}
type RequestProgVeh struct {
	Itinerary int
}
type ListProgVeh struct {
	Itinerary int
}
type RequestDriver struct {
	Driver string
}
type RequestTakeService struct {
}
