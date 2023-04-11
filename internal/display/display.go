package display

import (
	"time"
)

type Display interface {
	Init() error
	Close() error
	SwitchScreen(num int) error
	WriteText(label int, text ...string) error
	Popup(label int, text ...string) error
	PopupClose(label int) error
	Beep(repeat int, timeout time.Duration) error
	Verify() error
	Screen() (int, error)
	Reset() error
	Led(label int, state int) error
	KeyNum(prompt string) (int, error)
	Keyboard(prompt string) (string, error)
	Brightness(percent int) error
}
