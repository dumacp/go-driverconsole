//+build levis

package display

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/dumacp/go-levis"
	"github.com/dumacp/go-logs/pkg/logs"
)

var dayActual int = -1
var minActual int = -1
var statePuerta1 int = 0
var statePuerta2 int = 0
var timeoutRead time.Duration = 1 * time.Second

type display struct {
	dev          levis.Device
	screenActual int
	scratchData  []byte
}

const (
	addrConfirmation    int = 200
	addrEfectivoCounter int = 3
	addrTotalCounter    int = 2
	addrServiceTime     int = 10
	addrNoRoute         int = 110
	addrNameRoute       int = 20
	addrPercent         int = 4

	addrFrontalDoor int = 10
	addrBackDoor    int = 11
)

func NewDisplay(m interface{}) (Display, error) {

	dev, ok := m.(levis.Device)
	if !ok {
		return nil, fmt.Errorf("device is not LEVIS device")
	}
	display := &display{}
	display.dev = dev
	return display, nil
}

func (m *display) screen() int {
	return m.screenActual
}

func (m *display) close() {
}

func (m *display) route(routes string) error {

	if err := m.dev.WriteRegister(addrNameRoute,
		levis.EncodeFromChars([]byte(routes))); err != nil {
		return err
	}

	return nil

}

func (m *display) screenError(sError ...string) {

	textBytes := make([]byte, 0)

	for _, v := range sError {
		textBytes = append(textBytes, []byte(v)...)
		textBytes = append(textBytes, '\n')
	}

	text := levis.EncodeFromChars(textBytes[:len(textBytes)-1])
	m.dev.WriteRegister(addrConfirmation, text)
	m.switchScreen(1, false)
}

func (m *display) ingresos(usosEfectivo, usosTarjeta, usosParcial int) {

	if err := m.dev.WriteRegister(addrTotalCounter,
		[]uint16{uint16(usosEfectivo) + uint16(usosTarjeta)}); err != nil {
		logs.LogWarn.Println(err)
	}
	if err := m.dev.WriteRegister(addrEfectivoCounter,
		[]uint16{uint16(usosEfectivo)}); err != nil {
		logs.LogWarn.Println(err)
	}

}

func (m *display) ingresosPartial(usosParcial int) {

}

func (m *display) timeRecorrido(timeLapse int) {
	hours := uint((time.Duration(timeLapse) * time.Minute).Hours())
	minutes := uint((time.Duration(timeLapse) * time.Minute).Minutes()) % 60
	text := fmt.Sprintf("%02d:%02d", hours, minutes)
	fmt.Printf("debug time: %s\n", text)
	if err := m.dev.WriteRegister(addrServiceTime,
		levis.EncodeFromChars([]byte(text))); err != nil {
		logs.LogWarn.Println(err)
	}
}

func (m *display) selectionRuta() {
	fmt.Println("salir de SELECCION ruta GTT")
}

func (m *display) updateRuta(ruta, subruta string) {

}

func (m *display) alertBeep(repeat int) {
}

func (m *display) init() {

}

func (m *display) verifyReset() chan int {
	ch := make(chan int)
	/**/
	go func() {
		defer close(ch)

		log.Println("verifyReset stop")
	}()
	/**/
	return ch
}

func (m *display) switchScreen(screen int, active bool) error {
	/**/
	if m.screenActual == screen {
		return nil
	}
	/**/
	m.dev.WriteRegister(0, []uint16{uint16(screen)})
	m.screenActual = screen

	return nil
}

func (m *display) mainScreen() {
	m.dev.WriteRegister(0, []uint16{0})
}

func (m *display) disableSelectButton() {
}

func (m *display) updateDate(period int) {

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

func (m *display) keyNum(text string) {}

func (m *display) textConfirmation(sError ...string) {}

func (m *display) textError(sError ...string) {}
