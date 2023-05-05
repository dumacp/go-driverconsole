//go:build levis
// +build levis

package display

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/dumacp/go-levis"
)

var dayActual int = -1
var minActual int = -1
var statePuerta1 int = 0
var statePuerta2 int = 0
var timeoutRead time.Duration = 1 * time.Second

type TypeRegister int

const (
	LED = iota
	INPUT_TEXT
	INPUT_NUM
	ARRAY_PICT
)

type Register struct {
	Type   TypeRegister
	Addr   int
	Len    int
	Size   int
	Gap    int
	Toogle int
}

type display struct {
	dev            levis.Device
	screenActual   int
	scratchData    []byte
	notifications  []string
	lastUpdateDate time.Time
	inputs         int64
	outputs        int64
	label2addr     func(label Label) Register
}

const (
	SCREEN_INPUT_DRIVER = 3
	SCREEN_INPUT_ROUTE  = 2
)

func NewLevisDisplay(dev levis.Device) (Display, error) {
	display := &display{}
	display.dev = dev
	return display, nil
}

func (m *display) Init() error {
	return nil
}

func (m *display) Screen() (int, error) {
	return m.screenActual, nil
}

func (m *display) Close() error {
	return nil
}

func (m *display) Reset() error {
	return nil
}

func (m *display) SwitchScreen(num int) error {
	if err := m.dev.WriteRegister(0, []uint16{uint16(num)}); err != nil {
		return err
	}
	m.screenActual = num
	return nil
}

func (m *display) writeText(addr, length, size, gap int, text ...string) error {
	textBytes := make([]byte, 0)
	if length <= 1 {
		for i, v := range text {
			textBytes = append(textBytes, []byte(v)...)
			if i < len(text) {
				textBytes = append(textBytes, '\n')
			}
		}
		if size < len(textBytes) {
			return fmt.Errorf("len text in greather that register")
		}
		if err := m.dev.WriteRegister(addr, make([]uint16, size/2)); err != nil {
			return fmt.Errorf("error writeRegister: %s\n", err)
		}
		data := levis.EncodeFromChars(textBytes[:])
		if err := m.dev.WriteRegister(addr, data); err != nil {
			return err
		}
	} else {
		for i, v := range text {
			if i > length {
				break
			}
			if size < len(v) {
				return fmt.Errorf("len text in greather that register")
			}
			if err := m.dev.WriteRegister(addr+(i*gap), make([]uint16, size/2)); err != nil {
				return fmt.Errorf("error writeRegister: %s\n", err)
			}
		}
		for i, v := range text {
			if i > length {
				break
			}
			if err := m.dev.WriteRegister(addr+(i*gap),
				levis.EncodeFromChars([]byte(v))); err != nil {
				return fmt.Errorf("error writeRegister: %s\n", err)
			}
		}
	}
	return nil
}

func (m *display) WriteText(label Label, text ...string) error {
	reg := m.label2addr(label)
	if reg.Type != INPUT_TEXT {
		return fmt.Errorf("invalid data input")
	}
	return m.writeText(reg.Addr, reg.Len, reg.Size, reg.Gap, text...)
}

func (m *display) WriteNumber(label Label, num int64) error {
	reg := m.label2addr(label)
	if reg.Type != INPUT_NUM {
		return fmt.Errorf("invalid data input")
	}
	numBytes := make([]byte, reg.Size)
	switch reg.Size {
	case 2:
		binary.BigEndian.PutUint32(numBytes, uint32(num))
	case 4:
		binary.BigEndian.PutUint64(numBytes, uint64(num))
	default:
		return fmt.Errorf("invalis size (%d) to number input (%d)", reg.Size, num)
	}
	if err := m.dev.WriteRegister(reg.Addr, levis.EncodeFromBytes(numBytes)); err != nil {
		return fmt.Errorf("error writeRegister: %s\n", err)
	}
	return nil
}

func (m *display) Popup(label Label, text ...string) error {
	reg := m.label2addr(label)
	if err := m.writeText(reg.Addr, reg.Len, reg.Size, reg.Gap, text...); err != nil {
		return err
	}
	if err := m.dev.SetIndicator(reg.Toogle, true); err != nil {
		return err
	}
	return nil
}

func (m *display) PopupClose(label Label) error {
	reg := m.label2addr(label)
	if err := m.dev.SetIndicator(reg.Toogle, false); err != nil {
		return err
	}
	return nil
}

func (m *display) Beep(repeat int, timeout time.Duration) error {
	go func() {
		for range make([]int, repeat) {
			if err := m.dev.SetIndicator(23, true); err != nil {
				fmt.Println(err)
			}
			time.Sleep(timeout)
			if err := m.dev.SetIndicator(23, false); err != nil {
				fmt.Println(err)
			}
		}
	}()
	return nil
}

func (m *display) Verify() error {
	return nil
}

func (m *display) Led(label Label, state int) error {
	reg := m.label2addr(label)
	if reg.Type == LED {
		if err := m.dev.SetIndicator(reg.Addr, state == 0); err != nil {
			return err
		}
		return nil
	}

	return m.dev.WriteRegister(reg.Addr, []uint16{uint16(state)})
}

func (m *display) KeyNum(prompt string) (int, error) {
	return 0, nil
}

func (m *display) Keyboard(prompt string) (string, error) {
	return "", nil
}

func (m *display) Brightness(percent int) error {
	return nil
}

///////////////////////////////////////
//////////////////////////////////////////////////////////////

/**

func (m *display) screenError(sError ...string) error {

	textBytes := make([]byte, 0)

	for _, v := range sError {
		textBytes = append(textBytes, []byte(v)...)
		textBytes = append(textBytes, '\n')
	}

	text := levis.EncodeFromChars(textBytes[:len(textBytes)-1])
	if err := m.dev.WriteRegister(addrConfirmation, text); err != nil {
		return err
	}
	return m.switchScreen(1, false)
}

func (m *display) ingresos(usosEfectivo, usosTarjeta, usoParcial int) error {

	if err := m.dev.WriteRegister(addrEfectivoCounter,
		[]uint16{uint16(usosEfectivo)}); err != nil {
		return err
	}
	if err := m.dev.WriteRegister(addrTagCounter,
		[]uint16{uint16(usosTarjeta)}); err != nil {
		return err
	}
	return nil
}

func (m *display) counters(inputs1, outputs1, inputs2, outputs2 int64) error {

	if m.inputs != inputs1+inputs2 {
		m.inputs = inputs1 + inputs2
		if err := m.dev.WriteRegister(addrInputs,
			[]uint16{uint16((inputs1 + inputs2) & 0xFFFF), uint16((inputs1 + inputs2) & 0xFFFF0000)}); err != nil {
			return err
		}
	}
	if m.outputs != outputs1+outputs2 {
		m.outputs = outputs1 + outputs2
		if err := m.dev.WriteRegister(addrOutputs,
			[]uint16{uint16((outputs1 + outputs2) & 0xFFFF), uint16((outputs1 + outputs2) & 0xFFFF0000)}); err != nil {
			return err
		}
	}
	return nil
}

func (m *display) eventCount(input, output int) error {

	if m.inputs > 0 && input > 0 {
		m.inputs += int64(input)
		if err := m.dev.WriteRegister(addrInputs,
			[]uint16{uint16((m.inputs) & 0xFFFF), uint16((m.inputs) & 0xFFFF0000)}); err != nil {
			return err
		}

	}
	if m.outputs > 0 && output > 0 {
		m.outputs += int64(output)
		if err := m.dev.WriteRegister(addrOutputs,
			[]uint16{uint16((m.outputs) & 0xFFFF), uint16((m.outputs) & 0xFFFF0000)}); err != nil {
			return err
		}
	}
	return nil
}

func (m *display) ingresosPartial(usosParcial int) {

}

func (m *display) timeRecorrido(timeLapse int) {

}

func (m *display) selectionRuta() {
	fmt.Println("salir de SELECCION ruta GTT")
}

func (m *display) updateRuta(ruta, subruta string) {

}

func (m *display) alertBeep(repeat int) {
	go func() {
		for range make([]int, repeat) {
			if err := m.dev.SetIndicator(23, true); err != nil {
				fmt.Println(err)
			}
			time.Sleep(1 * time.Second)
			if err := m.dev.SetIndicator(23, false); err != nil {
				fmt.Println(err)
			}
		}
	}()
}

func (m *display) Init() error {
	return nil
}

func (m *display) verifyReset(quit chan int, ctx actor.Context) {
}

func (m *display) switchScreen(screen int, active bool) error {
	if err := m.dev.WriteRegister(0, []uint16{uint16(screen)}); err != nil {
		return err
	}
	m.screenActual = screen
	return nil
}

func (m *display) mainScreen() error {
	if err := m.dev.WriteRegister(0, []uint16{0}); err != nil {
		return err
	}
	m.screenActual = 0
	return nil
}

func (m *display) disableSelectButton() {
}

func (m *display) updateDate(period int) error {
	tNow := time.Now()
	if m.lastUpdateDate.Minute() == tNow.Minute() && m.lastUpdateDate.Hour() == tNow.Hour() {
		return nil
	}
	m.lastUpdateDate = tNow
	text := levis.EncodeFromChars([]byte(tNow.Format("2006/01/02 15:04")))
	if err := m.dev.WriteRegister(addrTimeDate, text); err != nil {
		return err
	}
	return nil

}

func (m *display) reset() {
}

func (m *display) setPuerta(id, state int) {

}

func (m *display) textInput(stext ...string) {

}

var cacheTextInput string = ""

func (m *display) addTextInput(schar string) {
	cacheTextInput = strings.Join([]string{cacheTextInput, schar}, "")
	log.Printf("cacheTextInput: %s", cacheTextInput)
	m.textInput(cacheTextInput)
}

func (m *display) delTextInput(count int) {
	if len(cacheTextInput) > count {
		cacheTextInput = cacheTextInput[0 : len(cacheTextInput)-count]
	} else {
		cacheTextInput = ""
	}
	m.textInput(cacheTextInput)
}

func (m *display) clearTextInput() {
	m.textInput("")
}

func (m *display) recorridoPercent(value int) error {
	return m.dev.WriteRegister(4, []uint16{uint16(value)})
}

func (m *display) doors(value [2]int) error {

	if err := m.dev.SetIndicator(addrFrontalDoor, value[0] == 0); err != nil {
		return err
	}
	if err := m.dev.SetIndicator(addrBackDoor, value[1] == 0); err != nil {
		return err
	}

	return nil
}

func (m *display) textConfirmationMainScreen(timeout time.Duration, sError ...string) error {
	return m.messageInMainScreen(addrConfirmationToggleMainScreen, addrConfirmationTextMainScreen, 400, timeout, sError...)
}

func (m *display) warningInMainScreen(timeout time.Duration, sError ...string) error {
	return m.messageInMainScreen(addrConfirmationToggleMainScreenErr, addrConfirmationTextMainScreenErr, 300, timeout, sError...)
}

func (m *display) messageInMainScreen(addrToogle, addrText, freq int, timeout time.Duration, sError ...string) error {

	m.alertBeep(1)

	if err := m.dev.WriteRegister(addrText, make([]uint16, 64)); err != nil {
		return fmt.Errorf("error writRegister: %w", err)
	}
	textBytes := make([]byte, 0)

	for _, v := range sError {
		if len(v) > 26 {
			for _, vv := range SplitHeader(v, 26) {
				textBytes = append(textBytes, []byte(vv)...)
				textBytes = append(textBytes, '\n')
			}
		} else {
			textBytes = append(textBytes, []byte(v)...)
		}
		textBytes = append(textBytes, '\n')
	}

	text := levis.EncodeFromChars(textBytes[:len(textBytes)-1])
	if err := m.dev.WriteRegister(addrText, text); err != nil {
		return err
	}
	if err := m.dev.SetIndicator(addrToogle, true); err != nil {
		return err
	}
	go func() {
		time.Sleep(2 * time.Second)
		if err := m.dev.SetIndicator(addrToogle, false); err != nil {
			fmt.Println(err)
		}
	}()
	return nil
}

func (m *display) textConfirmation(sError ...string) error {

	m.alertBeep(1)

	if err := m.dev.WriteRegister(addrConfirmation, make([]uint16, 120)); err != nil {
		return fmt.Errorf("error writRegister: %w", err)
	}
	textBytes := make([]byte, 0)

	for _, v := range sError {
		if len(v) > 26 {
			for _, vv := range SplitHeader(v, 26) {
				textBytes = append(textBytes, []byte(vv)...)
				textBytes = append(textBytes, '\n')
			}
		} else {
			textBytes = append(textBytes, []byte(v)...)
		}
		textBytes = append(textBytes, '\n')
	}

	text := levis.EncodeFromChars(textBytes[:len(textBytes)-1])
	if err := m.dev.WriteRegister(addrConfirmation, text); err != nil {
		return err
	}
	return m.switchScreen(5, false)
}

func (m *display) textError(sError ...string) error {

	if err := m.dev.WriteRegister(addrError, make([]uint16, 120)); err != nil {
		return fmt.Errorf("error writRegister: %w", err)
	}
	textBytes := make([]byte, 0)

	for _, v := range sError {
		if len(v) > 26 {
			for _, vv := range SplitHeader(v, 26) {
				textBytes = append(textBytes, []byte(vv)...)
				textBytes = append(textBytes, '\n')
			}
		} else {
			textBytes = append(textBytes, []byte(v)...)
		}
		textBytes = append(textBytes, '\n')
	}

	text := levis.EncodeFromChars(textBytes[:len(textBytes)-1])
	if err := m.dev.WriteRegister(addrError, text); err != nil {
		return fmt.Errorf("error writRegister: %w", err)
	}
	return m.switchScreen(6, false)
}

func (m *display) gpsstate(state int) error {
	return m.dev.SetIndicator(addrIconGPS, state == 0)
}
func (m *display) netstate(state int) error {
	return m.dev.SetIndicator(addrIconNET, state == 0)
}

func (m *display) addnotification(msg string) error {
	MAX_LEN := 10
	if len(m.notifications) <= 0 {
		m.notifications = make([]string, 0)
		m.notifications = append(m.notifications, msg)

	} else if len(m.notifications) > MAX_LEN {
		for i, v := range m.notifications[1:] {
			m.notifications[i] = v
		}
		m.notifications[MAX_LEN] = msg
	} else {
		m.notifications = append(m.notifications, msg)
	}
	fmt.Printf("notifs: %v\n", m.notifications)
	return nil
}
func (m *display) shownotifications() error {
	for i := range make([]int, 10) {
		fmt.Printf("%d\n", i)
		if err := m.dev.WriteRegister(addrAlarms+(i*100), make([]uint16, 32)); err != nil {
			return fmt.Errorf("error writRegister: %w", err)
		}
	}
	arrayText := make([]string, 0)
	for _, v := range m.notifications {
		arrayText = append(arrayText, fmt.Sprintf(" ==> %s", v))
	}
	for i, data := range arrayText {
		text := levis.EncodeFromChars([]byte(data))
		if err := m.dev.WriteRegister(addrAlarms+(i*100), text); err != nil {
			return fmt.Errorf("error writRegister (%s): %w", data, err)
		}
	}
	return m.switchScreen(3, true)
}
func (m *display) setBrightness(percent int) error {
	return nil
}

func (m *display) inputValue(initialText string, screen int) {
}
**/
