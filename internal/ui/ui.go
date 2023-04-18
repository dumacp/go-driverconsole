package ui

import (
	"time"

	"github.com/dumacp/go-driverconsole/internal/display"
)

type ui struct {
	disp display.Display
}

func New(display display.Display) UI {
	u := &ui{}
	u.disp = display
	return nil
}

type UI interface {
	Init() error
	MainScreen() error
	TextWarning(text ...string) error
	TextConfirmation(text ...string) error
	TextConfirmationPopup(timeout time.Duration, text ...string) error
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
	ShowNotifications() error
	ShowProgDriver() error
	ShowProgVeh() error
	ShowStats() error
	Brightness(percent int) error
}

func (u *ui) Init() error {
	return u.disp.Init()
}

func (u *ui) MainScreen() error {
	return u.disp.SwitchScreen(0)
}

func (u *ui) TextWarning(text ...string) error {
	return u.disp.WriteText(display.WARNING_TEXT, text...)
}

func (u *ui) TextConfirmation(text ...string) error {
	return u.disp.WriteText(display.CONFIRMATION_TEXT, text...)
}

func (u *ui) TextConfirmationPopup(timeout time.Duration, sText ...string) error {
	return u.disp.Popup(display.POPUP_TEXT, sText...)
}

func (u *ui) TextWarningPopup(timeout time.Duration, sText ...string) error {
	return u.disp.Popup(display.POPUP_WARN_TEXT, sText...)
}

func (u *ui) Inputs(in int) error {
	return u.disp.WriteNumber(display.INPUT_NUM, int64(in))
}

func (u *ui) Outputs(out int) error {
	return u.disp.WriteNumber(display.INPUT_NUM, int64(out))
}

func (u *ui) DeviationInputs(dev int) error {
	return u.disp.WriteNumber(display.INPUT_NUM, int64(dev))
}

func (u *ui) Route(route ...string) error {
	return u.disp.WriteText(display.ROUTE_TEXT, route...)
}

func (u *ui) Driver(data string) error {
	return u.disp.WriteText(display.DRIVER_TEXT, data)
}

func (u *ui) Beep(repeat int, timeout time.Duration) error {
	return nil
}

func (u *ui) Date(date time.Time) error {
	panic("not implemented") // TODO: Implement
}

func (u *ui) Screen(num int, force bool) error {
	panic("not implemented") // TODO: Implement
}

func (u *ui) GetScreen() (int, error) {
	panic("not implemented") // TODO: Implement
}

func (u *ui) KeyNum(prompt string) (int, error) {
	panic("not implemented") // TODO: Implement
}

func (u *ui) Keyboard(prompt string) (string, error) {
	panic("not implemented") // TODO: Implement
}

func (u *ui) Doors(state ...bool) error {
	panic("not implemented") // TODO: Implement
}

func (u *ui) Gps(state bool) error {
	panic("not implemented") // TODO: Implement
}

func (u *ui) Network(state bool) error {
	panic("not implemented") // TODO: Implement
}

func (u *ui) AddNotifications(add string) error {
	panic("not implemented") // TODO: Implement
}

func (u *ui) Brightness(percent int) error {
	panic("not implemented") // TODO: Implement
}
