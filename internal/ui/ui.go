package ui

import (
	"time"
)

type UI interface {
	Init() error
	MainScreen() error
	TextWarning(sError ...string) error
	TextConfirmation(sError ...string) error
	TextConfirmationPopup(timeout time.Duration, sError ...string) error
	TextWarningPopup(timeout time.Duration, sText ...string) error
	Inputs(in int) error
	Outputs(out int) error
	DeviationInputs(dev int) error
	Route(route ...string) error
	Driver(data string) error
	Beep(repeat int, timeout time.Duration) error
	Date(date time.Time) error
	Screen(num int, force bool) error
	GetScreen() (int, error)
	KeyNum(prompt string) (int, error)
	Keyboard(prompt string) (string, error)
	Doors(state ...bool) error
	Gps(state bool) error
	Network(state bool) error
	AddNotifications(add string) error
	Brightness(percent int) error
}
