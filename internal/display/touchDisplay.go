//+build gtt43 !levis

package display

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/dumacp/go-logs/pkg/logs"
	"github.com/dumacp/matrixorbital/gtt43a"
)

var statePuerta1 int = 0
var statePuerta2 int = 0
var timeoutRead time.Duration = 1 * time.Second

type touchDisplay struct {
	dev              gtt43a.Display
	screenActual     int
	scratchData      []byte
	inputsCash       int
	inputsApp        int
	inputsParcial    int
	doorss           [2]int
	percentRecorrido int
	timeLapse        time.Time

	routeName string
}

const (
	textCivica       int = 25
	textEfectivo     int = 27
	labelParcial     int = 5
	textParcial      int = 9
	usosSliceValue   int = 3
	labelError       int = 20
	labelTextInput   int = 29
	tittleRuta       int = 2
	textRuta         int = 1
	buttonEnter      int = 18
	buttonUp         int = 19
	buttonSelectPaso int = 15
	buttonEnterPaso  int = 16
	buttonRecorrido  int = 10
	buttonCounter    int = 6
	buttonGrid1      int = 0
	buttonGrid2      int = 3
	buttonGrid3      int = 6
	buttonGrid4      int = 1
	buttonGrid5      int = 4
	buttonGrid6      int = 7
	buttonGrid7      int = 2
	buttonGrid8      int = 5
	buttonGrid9      int = 8
	buttonGrid0      int = 10
	buttonGridEnter  int = 11
	buttonGridDel    int = 9
	timeRecorrido    int = 17
	timeHour         int = 24
	timeDate         int = 28
	textBoxError     int = 20

	addrNameRoute    int = 1
	addrPercent      int = 5
	addrConfirmation int = 32
	addrWarning      int = 34

	addrEfectivoCounter int = 3
	addrTotalCounter    int = 4

	addrServiceTime int = 8

	addrFrontalDoor int = 2
	addrBackDoor    int = 3

	addrTimeDate       int = 2
	addrResetRecorrido int = 9

	addrNoRoute int = 13
)

func NewDisplay(m interface{}) (Display, error) {

	dev, ok := m.(gtt43a.Display)
	if !ok {
		return nil, fmt.Errorf("device is not GTT43 device")
	}
	display := &touchDisplay{}
	display.dev = dev
	return display, nil
}

func (m *touchDisplay) screen() int {
	return m.screenActual
}

func (m *touchDisplay) close() {
	m.dev.Close()
}

func (m *touchDisplay) route(routes string) error {

	if len(routes) <= 0 {
		return nil
	}

	m.routeName = routes

	if m.screenActual != 1 {
		return nil
	}

	text := m.dev.SetPropertyText(addrNameRoute,
		gtt43a.ButtonText)

	if err := text(routes); err != nil {
		return err
	}

	return nil

}

func (m *touchDisplay) screenError(sError ...string) {

	if m.screenActual == 3 {
		return
	}
	m.switchScreen(3, true)
	// go func() {
	for i := 0; i < 3; i++ {
		m.dev.BuzzerActive(400, 500)
		time.Sleep(100 * time.Millisecond)
	}
	// }()
	text := m.dev.SetPropertyText(int(addrConfirmation), gtt43a.LabelText)
	s1 := ""
	for _, v := range sError {
		s1 = fmt.Sprintf("%s%s\n", s1, v)
	}

	if len(s1) > 0 {
		text(s1[:len(s1)-1])
	}
}

func (m *touchDisplay) textError(sError ...string) {
	if m.screenActual == 4 {
		return
	}
	m.dev.AnimationStopAll()
	m.dev.AnimationSetFrame(0, 0)
	m.dev.AnimationSetFrame(1, 0)

	if err := m.switchScreen(4, true); err != nil {
		log.Printf("switch screen err: %s", err)
		return
	}
	ch := make(chan int)
	go func() {
		for i := 0; i < 2; i++ {
			m.dev.BuzzerActive(250, 500)
			time.Sleep(600 * time.Millisecond)
		}
		ch <- 1
	}()
	text := m.dev.SetPropertyText(int(addrWarning), gtt43a.LabelText)
	s1 := ""
	for _, v := range sError {
		s1 = fmt.Sprintf("%s%s\n", s1, v)
	}

	if len(s1) > 0 {
		text(s1[:len(s1)-1])
	}
	<-ch
	//m.dev.Version()
}

func (m *touchDisplay) textConfirmation(sError ...string) {
	if m.screenActual == 3 {
		return
	}
	m.dev.AnimationStopAll()
	m.dev.AnimationSetFrame(0, 0)
	m.dev.AnimationSetFrame(1, 0)

	m.switchScreen(3, true)
	ch := make(chan int)
	go func() {
		for i := 0; i < 2; i++ {
			m.dev.BuzzerActive(400, 500)
			time.Sleep(600 * time.Millisecond)
		}
		ch <- 1
	}()
	text := m.dev.SetPropertyText(int(addrConfirmation), gtt43a.LabelText)
	s1 := ""
	for _, v := range sError {
		s1 = fmt.Sprintf("%s%s\n", s1, v)
	}

	if len(s1) > 0 {
		text(s1[:len(s1)-1])
	}
	<-ch
	//m.dev.Version()
}

func (m *touchDisplay) ingresos(usosEfectivo, usosCivica, usosParcial int) {
	/**/
	m.inputsCash = usosEfectivo
	m.inputsApp = usosCivica
	m.inputsParcial = usosParcial
	if m.screenActual != 1 {
		return
	}
	textUsosEfectivo := m.dev.SetPropertyText(int(addrEfectivoCounter), gtt43a.LabelText)
	textUsosEfectivo(fmt.Sprintf("%d", usosEfectivo))

	textUsosTotal := m.dev.SetPropertyText(int(addrTotalCounter), gtt43a.LabelText)
	textUsosTotal(fmt.Sprintf("%d", usosCivica+usosEfectivo))
}

func (m *touchDisplay) ingresosPartial(usosParcial int) {
	m.inputsParcial = usosParcial
	if m.screenActual != 1 {
		return
	}
	textUsosParcial := m.dev.SetPropertyText(int(textParcial), gtt43a.LabelText)
	textUsosParcial(fmt.Sprintf("%d", usosParcial))
}

func (m *touchDisplay) timeRecorrido(timeLapse int) {

	var timeRef time.Time
	if timeLapse == 0 {
		m.timeLapse = time.Now()
		timeRef = time.Now()
	} else if timeLapse < 0 {
		m.timeLapse = time.Time{}
	} else {
		timeRef = time.Now()
	}

	if m.screenActual != 1 {
		return
	}

	since := timeRef.Sub(m.timeLapse)

	hours := uint((since).Hours())
	minutes := uint((since).Minutes()) % 60
	text := fmt.Sprintf("%02d:%02d", hours, minutes)
	fmt.Printf("debug time: %s\n", text)

	textRecorrido := m.dev.SetPropertyText(int(addrServiceTime), gtt43a.LabelText)
	textRecorrido(text)
}

func (m *touchDisplay) selectionRuta() {
	fmt.Println("salir de SELECCION ruta GTT")
}

func (m *touchDisplay) updateRuta(ruta, subruta string) {

}

func (m *touchDisplay) alertBeep(repeat int) {
	for i := 0; i < repeat; i++ {
		m.dev.BuzzerActive(1000, 150)
		time.Sleep(time.Millisecond * 400)
	}
}

func (m *touchDisplay) init() {

	id := m.screen()
	log.Printf("LISTEN init: %d\n", id)
	if id != 1 {
		if err := m.switchScreen(1, true); err != nil {
			logs.LogWarn.Println(err)
		}
	}
	m.screenActual = 1
	// m.dev.Listen()
	log.Printf("LISTEN 1\n")
	time.Sleep(1 * time.Second)
	// sizeFont1 := m.dev.SetPropertyValueU16(int(textParcial), gtt43a.LabelFontSize)
	// sizeFont1(12)
	// m.scratchData = []byte{0x01, 0x10, 0x0A, 0xA0, 0x0B, 0xB0}
	// m.dev.WriteScratch(1, m.scratchData)

}

func (m *touchDisplay) verifyReset() chan int {
	ch := make(chan int)
	/**/
	go func() {
		defer close(ch)
		tick := time.NewTicker(time.Second * 10)
		defer tick.Stop()
		for {
			//log.Printf("memory scratchData: [% X]\n", m.scratchData)
			select {
			case <-tick.C:
				scratchpad, err := m.dev.ReadScratch(1, len(m.scratchData))
				if err != nil {
					log.Println(err)
					continue
				}
				log.Printf("now scratchData: [% X]\n", scratchpad)
				if len(scratchpad) < len(m.scratchData) {
					m.screenActual = 0
					ch <- 1
					continue
				}
				for i, b := range scratchpad {
					if b != m.scratchData[i] {
						m.screenActual = 0
						ch <- 1
						break
					}
				}
			}
		}
		// log.Println("verifyReset stop")
	}()
	/**/
	return ch
}

func (m *touchDisplay) switchScreen(screen int, active bool) error {
	/**/
	if m.screenActual == screen {
		return nil
	}
	log.Printf("switch screen: %d, active: %v", screen, active)

	// m.dev.AnimationSetFrame(0, 0)
	// m.dev.AnimationStartStop(0, 0)
	// m.dev.AnimationSetFrame(1, 0)
	// m.dev.AnimationStartStop(0, 0)

	/**/

	if active {
		if err := m.dev.RunScript(fmt.Sprintf("GTTPdirverconsole\\Screen%d\\Screen%d.bin", screen, screen)); err != nil {
			return err
		}
		time.Sleep(1000 * time.Millisecond)
	}
	m.screenActual = screen
	return nil
}

func (m *touchDisplay) mainScreen() {
	if m.screenActual != 1 {
		if err := m.switchScreen(1, true); err != nil {
			logs.LogWarn.Println(err)
		}
	}
	m.screenActual = 1
	m.updateDate(0)
	if !m.timeLapse.IsZero() {
		m.timeRecorrido(100)
		if err := m.dev.SetPropertyValueU8(addrResetRecorrido, gtt43a.ButtonState)(1); err != nil {
			logs.LogWarn.Println(err)
		}
	}
	m.ingresos(m.inputsCash, m.inputsApp, 0)
	m.doors(m.doorss)
	m.recorridoPercent(m.percentRecorrido)
	m.route(m.routeName)
	// m.setPuerta(LEDPuerta1, statePuerta1)
	// m.setPuerta(LEDPuerta2, statePuerta2)
	// time.Sleep(time.Millisecond * 800)
	// m.dev.Version()
}

func (m *touchDisplay) disableSelectButton() {
	setState := m.dev.SetPropertyValueU8(int(buttonSelectPaso), gtt43a.ButtonState)
	setState(0x00)
}

func (m *touchDisplay) updateDate(period int) {

	if m.screenActual != 1 {
		return
	}

	tNow := time.Now()

	textTime := m.dev.SetPropertyText(addrTimeDate, gtt43a.LabelText)
	textTime(tNow.Format("15:04"))
}

func (m *touchDisplay) reset() {
	m.dev.Reset()
}

func (m *touchDisplay) setPuerta(id, state int) {
	if state == 0 {
		m.dev.AnimationSetFrame(id, 0)
	}
	// if id == LEDPuerta1 {
	// 	statePuerta1 = state
	// }
	// if id == LEDPuerta2 {
	// 	statePuerta2 = state
	// }

	m.dev.AnimationStartStop(id, state)
}

func (m *touchDisplay) textInput(stext ...string) {
	// if m.screenActual < 3 {
	// 	log.Println("STOP animation")
	// 	m.dev.AnimationStopAll()
	// 	// m.dev.AnimationSetFrame(LEDPuerta1, 0)
	// 	// m.dev.AnimationSetFrame(LEDPuerta2, 0)
	// }
	// m.switchScreen(4)
	// text := m.dev.SetPropertyText(int(labelTextInput), gtt43a.LabelText)
	// s1 := ""
	// for _, v := range stext {
	// 	s1 = fmt.Sprintf("%s%s\n", s1, v)
	// }

	// if len(s1) > 0 {
	// 	text(s1[:len(s1)-1])
	// }
	//m.dev.Version()
}

var cacheTextInput string = ""

func (m *touchDisplay) addTextInput(schar string) {
	cacheTextInput = strings.Join([]string{cacheTextInput, schar}, "")
	log.Printf("cacheTextInput: %s", cacheTextInput)
	m.textInput(cacheTextInput)
}

func (m *touchDisplay) delTextInput(count int) {
	if len(cacheTextInput) > count {
		cacheTextInput = cacheTextInput[0 : len(cacheTextInput)-count]
	} else {
		cacheTextInput = ""
	}
	m.textInput(cacheTextInput)
}

func (m *touchDisplay) clearTextInput() {
	m.textInput("")
}

func (m *touchDisplay) doors(value [2]int) error {

	if len(value) < 2 {
		return fmt.Errorf("value is wrong")
	}

	m.doorss = value

	if m.screenActual != 1 {
		return nil
	}

	if value[0] == 0 {
		m.dev.AnimationSetFrame(addrFrontalDoor, 1)
	} else {
		m.dev.AnimationSetFrame(addrFrontalDoor, 2)
	}

	if value[1] == 0 {
		m.dev.AnimationSetFrame(addrBackDoor, 1)
	} else {
		m.dev.AnimationSetFrame(addrBackDoor, 2)
	}

	m.dev.AnimationStartStop(addrFrontalDoor, 1)
	m.dev.AnimationStartStop(addrFrontalDoor, 0)
	m.dev.AnimationStartStop(addrBackDoor, 1)
	m.dev.AnimationStartStop(addrBackDoor, 0)

	return nil
}

func (m *touchDisplay) recorridoPercent(value int) error {
	m.percentRecorrido = value
	if m.screenActual != 1 {
		return nil
	}
	valuefunc := m.dev.SetPropertyValueS16(addrPercent, gtt43a.SliderValue)

	if err := valuefunc(value); err != nil {
		return err
	}

	return nil
}

func (m *touchDisplay) keyNum(text string) {
	if m.screenActual == 2 {
		return
	}
	m.switchScreen(2, true)
	if len(text) > 0 {
		if err := m.dev.SetPropertyText(addrNoRoute, gtt43a.LabelText)(text); err != nil {
			logs.LogWarn.Println(err)
		}
	}
}
