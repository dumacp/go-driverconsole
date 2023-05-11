package display

import (
	"time"
)

type Display interface {
	Init(dev interface{}) error
	Close() error
	SwitchScreen(num int) error
	WriteText(label int, text ...string) error
	WriteNumber(label int, num int64) error
	Popup(label int, text ...string) error
	PopupClose(label int) error
	Beep(repeat int, timeout time.Duration) error
	Verify() error
	Screen() (int, error)
	Reset() error
	Led(label int, state int) error
	ArrayPict(label int, state int) error
	KeyNum(prompt string) (int, error)
	Keyboard(prompt string) (string, error)
	Brightness(percent int) error
	ReadBytes(label int) ([]byte, error)
	DeviceRaw() (interface{}, error)
}

// func New(devi interface{}) (Display, error) {

// 	switch dev := devi.(type) {
// 	case levis.Device:
// 		return NewPiDisplay(dev)
// 	}

// 	return nil, fmt.Errorf("Display device not foud")
// }
