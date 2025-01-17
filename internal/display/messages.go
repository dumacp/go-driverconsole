package display

import "time"

type DeviceMsg struct {
	Device interface{}
}
type Reset struct {
	// CountManual  int
	// CountParcial int
	// CountAppFare int
	// Route        string
	// Itininerary  string
}
type UpdateMainScreen struct {
	// timeLapse    int
	// CountManual  int
	// CountParcial int
	// CountAppFare int
	// Route        string
	// Itininerary  string
}
type UpdateDate struct{}
type DisplayCount struct {
	CountManual  int
	CountParcial int
	CountAppFare int
}
type DisplayError struct {
	Data string
}
type Route struct {
	Route       string
	Itininerary string
}
type Itininerary struct {
	Data string
}
type TimeLapse struct {
	Data int
}
type StopRecorrido struct{}
type InitRecorrido struct{}
type ResetCounter struct{}
type UpVoc struct{}
type EnterVoc struct{}
type DisableSelectPASO struct{}
type EnterPASO struct{}
type AddText struct {
	Text string
}
type DelText struct {
	Count int
}
type EnterText struct{}

type MsgRoutes struct {
	Data map[int]string
}
type DisplayDeviceError struct {
	Err string
}

///////////////////////

// AckMsg is a message to indicate response
type AckMsg struct {
	Error error
}

// InitMsg is a message to initialize the Display.
type InitMsg struct{}

// CloseMsg is a message to close the Display.
type CloseMsg struct{}

// SwitchScreenMsg is a message to switch to a specific screen on the Display.
type SwitchScreenMsg struct {
	Num int // Screen number to switch to.
}

// WriteTextMsg is a message to write text to a specific area on the Display.
type WriteTextMsg struct {
	Label int
	Text  []string
}

// ReadTextMsg is a message to write text to a specific area on the Display.
type ReadTextMsg struct {
	Label int
	Len   int
}
type ResponseTextMsg struct {
	Label int
	Value string
	Error error
}

// ReadNumMsg is a message to write text to a specific area on the Display.
type ReadNumMsg struct {
	Label int
}
type ResponseNumMsg struct {
	Label int
	Value int64
	Error error
}

// ReadBytesMsg is a message to write text to a specific area on the Display.
type ReadBytesMsg struct {
	Label int
}
type ResponseBytesMsg struct {
	Label int
	Value []byte
	Error error
}

// WriteNumberMsg is a message to write a number to a specific area on the Display.
type WriteNumberMsg struct {
	Label int
	Num   int64
}

// ArrayPict is a message to write a number to a specific area on the Display.
type ArrayPictMsg struct {
	Label int
	Num   int
}

// PopupMsg is a message to display a popup with text on the Display.
type PopupMsg struct {
	Label int
	Text  []string
	// RestoreData map[int]interface{}
	// Temout      time.Duration
}

// PopupCloseMsg is a message to close a popup on the Display.
type PopupCloseMsg struct {
	Label int
}

// BeepMsg is a message to play a beep sound on the Display.
type BeepMsg struct {
	Repeat int
	Duty   int
	Period time.Duration
}

// VerifyMsg is a message to perform a verification on the Display.
type VerifyMsg struct{}

// ScreenMsg is a message to get the current screen number on the Display.
type ScreenMsg struct{}

// ScreenResponseMsg is a message to return the current screen number and any error encountered.
type ScreenResponseMsg struct {
	Screen int
	Error  error
}

// ResetMsg is a message to reset the Display.
type ResetMsg struct{}

// LedMsg is a message to control an LED on the Display.
type LedMsg struct {
	Label int
	State int
}

// KeyNumMsg is a message to prompt the user to enter a number on the Display.
type KeyNumMsg struct {
	Prompt string
}

// KeyNumResponseMsg is a message to return the entered number and any error encountered.
type KeyNumResponseMsg struct {
	Num   int
	Error error
}

// KeyboardMsg is a message to prompt the user to enter text on the Display.
type KeyboardMsg struct {
	Prompt string
}

// KeyboardResponseMsg is a message to return the entered text and any error encountered.
type KeyboardResponseMsg struct {
	Text  string
	Error error
}

// BrightnessMsg is a message to set the brightness of the Display.
type BrightnessMsg struct {
	Percent int
}
