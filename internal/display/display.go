package display

import (
	"time"
)

type Display interface {
	Init() error
	Close() error
	SwitchScreen(num int) error
	WriteText(label Label, text ...string) error
	WriteNumber(label Label, num int64) error
	Popup(label Label, text ...string) error
	PopupClose(label Label) error
	Beep(repeat int, timeout time.Duration) error
	Verify() error
	Screen() (int, error)
	Reset() error
	Led(label Label, state int) error
	KeyNum(prompt string) (int, error)
	Keyboard(prompt string) (string, error)
	Brightness(percent int) error
}
