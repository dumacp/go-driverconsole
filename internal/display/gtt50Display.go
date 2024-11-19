package display

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/dumacp/matrixorbital/gtt43a"
)

type gtt50Display struct {
	// Campos de tu estructura gtt50Display
	dev          gtt43a.Display
	screenActual int
	label2addr   func(label int) Register
	enable       bool
}

func NewGtt50Display(label2addr func(label int) Register) Display {
	display := &gtt50Display{}
	display.label2addr = label2addr

	return display
}

func (m *gtt50Display) Init(dev interface{}) error {
	pi, ok := dev.(gtt43a.Display)
	if !ok {
		var ii gtt43a.Display
		return fmt.Errorf("device is not %T", ii)
	}
	m.dev = pi
	m.screenActual = 1
	return nil
}

func (m *gtt50Display) Close() error {
	if m.dev == nil {
		return nil
	}
	m.dev.Close()
	return nil
}

func (m *gtt50Display) SwitchScreen(num int) error {
	// Implementación de SwitchScreen
	reg := m.label2addr(num)
	fmt.Printf("reg: %+v\n", reg)
	if m.screenActual == reg.Addr {
		return nil
	}
	fmt.Printf("switch screen: %d\n", reg.Addr)

	if err := m.dev.RunScript(fmt.Sprintf("GTTProject1\\Screen%d\\Screen%d.bin", reg.Addr, reg.Addr)); err != nil {
		return err
	}
	time.Sleep(1000 * time.Millisecond)
	m.screenActual = reg.Addr

	return nil
}

func (m *gtt50Display) WriteText(label int, text ...string) error {
	// Implementación de WriteTtext
	reg := m.label2addr(label)
	fmt.Printf("reg: %+v\n", reg)
	if reg.Type != INPUT_TEXT {
		return fmt.Errorf("invalid data input")
	}

	textFinal := ""
	switch {
	case reg.Len <= 1 && len(text) > 0:
		textFinal = strings.Join(text, "\n")
	case reg.Len > 0:
		maxSize := 0
		for _, v := range text {
			if len(v) > maxSize {
				maxSize = len(v)
			}
		}
		if maxSize > reg.Size {
			return fmt.Errorf("len text in greather that register (%d > %d)", maxSize, reg.Size)
		}
		if len(text) > reg.Len {
			textFinal = strings.Join(text[0:reg.Len], "\n")
		} else {
			textFinal = strings.Join(text[:], "\n")
		}
	}

	result := m.dev.SetPropertyText(reg.Addr,
		gtt43a.LabelText)

	return result(textFinal)
}

func (m *gtt50Display) WriteNumber(label int, num int64) error {
	// Implementación de WriteNumber
	reg := m.label2addr(label)
	fmt.Printf("reg: %+v\n", reg)
	if reg.Type != INPUT_NUM {
		return fmt.Errorf("invalid data input")
	}
	result := m.dev.SetPropertyText(reg.Addr,
		gtt43a.LabelText)

	return result(fmt.Sprintf("%d", num))
}

func (m *gtt50Display) Popup(label int, text ...string) error {
	// Implementación de Popup
	reg := m.label2addr(label)
	fmt.Printf("reg: %+v\n", reg)
	if reg.Type != INPUT_TEXT {
		return fmt.Errorf("invalid data input")
	}
	fmt.Printf("current screen: %d\n", m.screenActual)
	if m.screenActual != 1 {
		return nil
	}
	if err := m.dev.SetPropertyValueU8(reg.Toogle, gtt43a.ButtonState)(1); err != nil {
		return err
	}

	result := m.dev.SetPropertyText(reg.Addr, gtt43a.ButtonText)
	if err := result(strings.Join(text, "\n")); err != nil {
		return err
	}
	// go func() {
	// 	time.Sleep(3 * time.Second)
	// 	if err := m.dev.SetPropertyValueU8(reg.Toogle, gtt43a.ButtonState)(2); err != nil {
	// 		fmt.Printf("off popup (reg: %v) error: %s\n", reg, err)
	// 	}
	// }()
	return nil
}

func (m *gtt50Display) PopupClose(label int) error {
	// Implementación de PopupClose
	reg := m.label2addr(label)
	fmt.Printf("reg: %+v\n", reg)
	if err := m.dev.SetPropertyValueU8(reg.Toogle, gtt43a.ButtonState)(2); err != nil {
		return err
	}
	result := m.dev.SetPropertyText(reg.Addr, gtt43a.ButtonText)
	if err := result(""); err != nil {
		return err
	}
	// TODO: how make???

	return nil
}

func (m *gtt50Display) BeepWithContext(contxt context.Context, repeat, duty int, period time.Duration) error {
	// Implementación de Beep
	go func() {
		for range make([]int, repeat) {
			t0 := time.Now()
			tsleep := time.NewTimer(period)
			defer tsleep.Stop()
			tdown := time.NewTimer(200 * time.Millisecond)
			defer tdown.Stop()
			dutyDuration := int(time.Duration((float32(duty) / float32(100)) * float32(period)).Milliseconds())
			if err := m.dev.BuzzerActive(1000, dutyDuration); err != nil {
				fmt.Println(err)
			}
			select {
			case <-tsleep.C:
			case <-contxt.Done():
				return
			}
			fmt.Printf("(%d, %d, %s) millisecons = %s\n", repeat, duty, period, time.Since(t0))
		}
	}()
	return nil
}

func (m *gtt50Display) Beep(repeat, duty int, period time.Duration) error {
	return m.BeepWithContext(context.TODO(), repeat, duty, period)
}

var resetScratch = []byte{0, 1, 2, 3, 4, 5, 6, 0xA}

func (m *gtt50Display) Verify() error {
	// Implementación de Verify
	if !m.enable {
		if err := m.dev.WriteScratch(1, resetScratch); err != nil {
			return err
		}
	}
	scratchpad, err := m.dev.ReadScratch(1, len(resetScratch))
	if err != nil {
		m.enable = false
		return err
	}
	log.Printf("//////////// now scratchData: [% X]\n", scratchpad)
	if len(scratchpad) < len(resetScratch) {
		m.enable = false
		return fmt.Errorf("scartchpad is not the same")
	}
	for i, b := range scratchpad {
		if b != resetScratch[i] {
			m.enable = false
			return fmt.Errorf("scartchpad is not the same")
		}
	}
	m.enable = true
	return nil
}

func (m *gtt50Display) Screen() (int, error) {
	// Implementación de Screen
	return m.screenActual, nil
}

func (m *gtt50Display) Reset() error {
	// Implementación de Reset
	m.screenActual = 1
	if err := m.dev.Reset(); err != nil {
		return err
	}
	time.Sleep(1 * time.Second)
	return nil
}

func (m *gtt50Display) Led(label int, state int) error {
	// Implementación de Led
	reg := m.label2addr(label)
	fmt.Printf("reg: %+v\n", reg)
	if reg.Type == LED {
		if err := m.dev.SetPropertyValueU16(reg.Addr, gtt43a.VisualBitmap_SourceIndex)(state); err != nil {
			return err
		}
		return nil
	} else if reg.Type == BUTTON {
		if err := m.dev.SetPropertyValueU8(reg.Addr, gtt43a.ButtonState)(state); err != nil {
			return err
		}
		return nil
	}
	return nil
}

func (m *gtt50Display) ArrayPict(label int, state int) error {
	// Implementación de ArrayPict
	return nil
}

func (m *gtt50Display) KeyNum(prompt string) (int, error) {
	// Implementación de KeyNum
	return 0, nil
}

func (m *gtt50Display) Keyboard(prompt string) (string, error) {
	// Implementación de Keyboard
	return "", nil
}

func (m *gtt50Display) Brightness(percent int) error {
	// Implementación de Brightness
	data := percent * 255 / 100
	return m.dev.SetBacklightLegcay(data)
}

func (m *gtt50Display) ReadBytes(label int) ([]byte, error) {
	reg := m.label2addr(label)
	fmt.Printf("reg: %+v\n", reg)
	res, err := m.dev.ReadScratch(reg.Addr, reg.Size)
	if err != nil {
		// fmt.Println(err)
		return nil, err
	}

	// fmt.Printf("debug read bytes: %v\n", res)

	if len(res) <= 0 {
		return nil, fmt.Errorf("response is empty")
	}

	return res, nil
}

func (m *gtt50Display) DeviceRaw() (interface{}, error) {
	// Implementación de DeviceRaw
	return nil, nil
}
