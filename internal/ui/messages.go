package ui

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

// InitUIMsg is a message for initializing the UI.
type InitUIMsg struct{}

// MainScreenMsg is a message for displaying the main screen.
type MainScreenMsg struct{}

// TextWarningMsg is a message for displaying a warning text.
type TextWarningMsg struct {
	Text []string
}

// TextConfirmationMsg is a message for displaying a confirmation text.
type TextConfirmationMsg struct {
	Text []string
}

// TextConfirmationPopupMsg is a message for displaying a confirmation popup with a timeout.
type TextConfirmationPopupMsg struct {
	Timeout time.Duration
	Text    []string
}

// TextWarningPopupMsg is a message for displaying a warning popup with a timeout.
type TextWarningPopupMsg struct {
	Timeout time.Duration
	Text    []string
}

// InputsMsg is a message for setting the number of inputs.
type InputsMsg struct {
	In int
}

// OutputsMsg is a message for setting the number of outputs.
type OutputsMsg struct {
	Out int
}

// DeviationInputsMsg is a message for setting the number of deviation inputs.
type DeviationInputsMsg struct {
	Dev int
}

// RouteMsg is a message for displaying the route.
type RouteMsg struct {
	Route []string
}

// DriverMsg is a message for displaying the driver's data.
type DriverMsg struct {
	Data string
}

// BeepMsg is a message for making a beep sound.
type BeepMsg struct {
	Repeat  int
	Timeout time.Duration
}

// DateMsg is a message for displaying the date.
type DateMsg struct {
	Date time.Time
}

// ScreenMsg is a message for switching to a specific screen.
type ScreenMsg struct {
	Num   int
	Force bool
}

// GetScreenMsg is a message for getting the current screen number.
type GetScreenMsg struct{}
type ScreenResponseMsg struct {
	Num   int
	Error error
}

// KeyNumMsg is a message for getting a number from a keypad.
type KeyNumMsg struct {
	Prompt string
}
type KeyNumResponseMsg struct {
	Num   int
	Error error
}

// KeyboardMsg is a message for getting a string from a keyboard.
type KeyboardMsg struct {
	Prompt string
}
type KeyboarResponsedMsg struct {
	Text  string
	Error error
}

// DoorsMsg is a message for displaying the doors' state.
type DoorsMsg struct {
	State []bool
}

// GpsMsg is a message for displaying the GPS state.
type GpsMsg struct {
	State bool
}

// NetworkMsg is a message for displaying the network state.
type NetworkMsg struct {
	State bool
}

// AddNotificationsMsg is a message for adding notifications.
type AddNotificationsMsg struct {
	Add string
}

// ShowNotificationsMsg is a message for showing notifications.
type ShowNotificationsMsg struct{}

// ShowProgDriverMsg is a message for showing the driver's progress.
type ShowProgDriverMsg struct {
	Text []string
}

// ShowProgVehMsg is a message for showing the vehicle's progress.
type ShowProgVehMsg struct {
	Text []string
}

// ShowStatsMsg is a message for showing the statistics.
type ShowStatsMsg struct {
	Text []string
}

// BrightnessMsg is a message for setting the brightness.
type BrightnessMsg struct {
	Percent int
}

type ServiceCurrentStateMsg struct {
	State  int
	Prompt string
}

type AckMsg struct {
	Error error
}

type AddInputsHandlerMsg struct {
	handler actor.Actor
}
