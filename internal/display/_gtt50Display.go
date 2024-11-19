package display

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/dumacp/go-logs/pkg/logs"
	"github.com/dumacp/matrixorbital/gtt43a"
)

type display struct {
	dev          gtt43a.Display
	screenActual int
	inputsCash   int
	inputsApp    int
	counterInput int64
	// counterInput2    int64
	counterOutput int64
	// counterOutput2   int64
	inputsParcial    int
	percentRecorrido int
	timeLapse        time.Time

	routeName      string
	driverID       string
	netState       int
	gpsState       int
	notifications  []string
	lastRunScript  time.Time
	lastUpdateDate time.Time
}

const (
	SCREEN_INPUT_DRIVER = 3
	SCREEN_INPUT_ROUTE  = 2
)

const (
	buttonSelectPaso   int = 15
	addrCounterInputs  int = 3
	addrCounterOutputs int = 14

	addrNameRoute    int = 18
	addrIDDriver     int = 21
	addrPercent      int = 5
	addrConfirmation int = 70
	addrWarning      int = 29

	addrEfectivoCounter int = 16
	addrTagCounter      int = 15

	addrServiceTime int = 8

	addrTimeDate       int = 90
	addrResetRecorrido int = 9

	addrNoRoute int = 13

	addrConfirmationTextMainScreen      int = 17
	addrConfirmationToggleMainScreen    int = 17
	addrConfirmationTextMainScreenErr   int = 33
	addrConfirmationToggleMainScreenErr int = 33

	addrIconGPS int = 34
	addrIconNET int = 88

	// addrVehicleAgend  int = 78
	// addrDisplayAlarms int = 75

	// addrCenterSpeed     int = 116
	// addrLeftSpeed       int = 108
	// addrRightSpeed      int = 127
	// addrCenterSpeedText int = 120
	// addrLeftSpeedText   int = 117
	// addrRightSpeedText  int = 121

	// addrItineraryVehicle int = 114

	// addrCurrentVehicle int = 142
	// addrNextStop       int = 140

	addrTextAlarms int = 87
)

var resetScratch = []byte{0, 1, 2, 3, 4, 5, 6, 7}

func NewDisplay(m interface{}) (Display, error) {
	dev, ok := m.(gtt43a.Display)
	if !ok {
		return nil, fmt.Errorf("device is not GTT43 device")
	}
	display := &display{}
	display.dev = dev
	display.lastRunScript = time.Now().Add(-10 * time.Second)
	return display, nil
}

func (m *display) screen() int {
	return m.screenActual
}

func (m *display) close() {
	m.dev.Close()
}

func (m *display) route(routes string) error {

	if len(routes) <= 0 {
		return nil
	}
	m.routeName = routes
	if m.screenActual != 1 {
		return nil
	}
	text := m.dev.SetPropertyText(addrNameRoute,
		gtt43a.LabelText)

	if err := text(routes); err != nil {
		return err
	}
	return nil

}

func (m *display) screenError(sError ...string) error {

	if m.screenActual == 5 {
		return nil
	}
	if err := m.switchScreen(5, true); err != nil {
		return err
	}
	ch := make(chan int)
	go func() {
		for i := 0; i < 2; i++ {
			m.dev.BuzzerActive(250, 500)
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
		if err := text(s1[:len(s1)-1]); err != nil {
			return err
		}
	}
	<-ch
	return nil
}

func (m *display) textError(sError ...string) error {
	if m.screenActual == 5 {
		return nil
	}

	if err := m.switchScreen(5, true); err != nil {
		log.Printf("switch screen err: %s", err)
		return err
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
		if err := text(s1[:len(s1)-1]); err != nil {
			return err
		}
	}
	<-ch
	return nil
}

func (m *display) textConfirmation(sError ...string) error {
	fmt.Printf("current screen: %d\n", m.screenActual)
	if m.screenActual == 4 {
		return nil
	}
	if err := m.switchScreen(4, true); err != nil {
		return err
	}
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
		if err := text(s1[:len(s1)-1]); err != nil {
			return err
		}
	}
	<-ch
	return nil
}

func (m *display) textConfirmationMainScreen(timeout time.Duration, sError ...string) error {
	return m.messageInMainScreen(addrConfirmationToggleMainScreen, addrConfirmationTextMainScreen, 400, timeout, sError...)
}

func (m *display) warningInMainScreen(timeout time.Duration, sError ...string) error {
	return m.messageInMainScreen(addrConfirmationToggleMainScreenErr, addrConfirmationTextMainScreenErr, 300, timeout, sError...)
}

func (m *display) messageInMainScreen(addrButton, addrText, freq int, timeout time.Duration, sError ...string) error {
	fmt.Printf("current screen: %d\n", m.screenActual)
	if m.screen() != 1 {
		return nil
	}

	ch := make(chan int)
	m.dev.SetPropertyValueU8(addrButton, gtt43a.ButtonState)(1)
	text := m.dev.SetPropertyText(int(addrText), gtt43a.ButtonText)
	s1 := ""
	for _, v := range sError {
		s1 = fmt.Sprintf("%s%s\n", s1, v)
	}
	if len(s1) > 0 {
		if err := text(s1[:len(s1)-1]); err != nil {
			return err
		}
	}
	go func() {
		for i := 0; i < 2; i++ {
			m.dev.BuzzerActive(freq, 500)
			time.Sleep(600 * time.Millisecond)
		}
		time.Sleep(200 * time.Millisecond)
		ch <- 1
	}()
	<-ch
	time.Sleep(timeout)
	m.dev.SetPropertyValueU8(addrButton, gtt43a.ButtonState)(2)
	text("")
	m.ingresos(m.inputsCash, m.inputsApp, 0)
	m.counters(m.counterInput, 0, m.counterOutput, 0)
	return nil
}

func (m *display) ingresos(usosEfectivo, usosCivica, usosParcial int) error {
	/**/
	m.inputsCash = usosEfectivo
	m.inputsApp = usosCivica
	m.inputsParcial = usosParcial
	if m.screenActual != 1 {
		return nil
	}
	textUsosEfectivo := m.dev.SetPropertyText(int(addrEfectivoCounter), gtt43a.LabelText)
	if err := textUsosEfectivo(fmt.Sprintf("%d", usosEfectivo)); err != nil {
		return err
	}
	textUsosTotal := m.dev.SetPropertyText(int(addrTagCounter), gtt43a.LabelText)
	if err := textUsosTotal(fmt.Sprintf("%d", usosCivica)); err != nil {
		return err
	}
	return nil
}

func (m *display) ingresosPartial(usosParcial int) {

}

func (m *display) timeRecorrido(timeLapse int) {

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

func (m *display) selectionRuta() {
	fmt.Println("salir de SELECCION ruta GTT")
}

func (m *display) updateRuta(ruta, subruta string) {

}

func (m *display) alertBeep(repeat, duty int, period time.Duration) {
	for i := 0; i < repeat; i++ {
		m.dev.BuzzerActive(1000, 150)
		time.Sleep(time.Millisecond * 400)
	}
}

func (m *display) init() error {

	id := 0
	log.Printf("LISTEN init: %d\n", id)
	if err := m.switchScreen(1, true); err != nil {
		logs.LogWarn.Println(err)
		return err
	}
	// m.screenActual = 1
	time.Sleep(1 * time.Second)
	if err := m.dev.WriteScratch(1, resetScratch); err != nil {
		return err
	}
	return nil
}

func (m *display) verifyReset(quit chan int, ctx actor.Context) {

	rootctx := ctx.ActorSystem().Root
	self := ctx.Self()

	fnc := func(err error) {
		rootctx.Send(self, &DisplayDeviceError{
			Err: err.Error(),
		})
	}
	log.Println("*********** VERIFY RESET ***********")
	go func() {
		tick := time.NewTicker(time.Second * 10)
		defer tick.Stop()
		for {

			select {
			case <-tick.C:
				if err := func() error {
					var scratchpad []byte
					var err error
					for range []int{1, 2, 3} {
						scratchpad, err = m.dev.ReadScratch(1, len(resetScratch))
						if err != nil {
							continue
						}
						log.Printf("//////////// now scratchData: [% X]\n", scratchpad)
						if len(scratchpad) < len(resetScratch) {
							return fmt.Errorf("scartchpad is not same")
						}
						for i, b := range scratchpad {
							if b != resetScratch[i] {
								if err := m.dev.WriteScratch(1, resetScratch); err != nil {
									return err
								}
								return fmt.Errorf("scartchpad is not same")
							}
						}
						break
					}
					if err != nil {
						return err
					}
					return nil
				}(); err != nil {
					log.Println(err)
					fnc(err)
				}
			case <-quit:
				return
			}
		}
	}()
}

func (m *display) switchScreen(screen int, active bool) error {
	/**/
	if m.screenActual == screen {
		return nil
	}
	if time.Since(m.lastRunScript) <= 3*time.Second {
		return fmt.Errorf("lastScreen")
	}
	m.lastRunScript = time.Now()
	log.Printf("switch screen: %d, active: %v", screen, active)

	if active {
		if err := m.dev.RunScript(fmt.Sprintf("GTTProject1\\Screen%d\\Screen%d.bin", screen, screen)); err != nil {
			return err
		}
		time.Sleep(1000 * time.Millisecond)
	}
	m.screenActual = screen

	// if screen == 8 {
	// 	index := m.dev.SetPropertyValueU16(49, gtt43a.VisualBitmap_SourceIndex)
	// 	time.Sleep(2 * time.Second)
	// 	index(1)
	// 	time.Sleep(2 * time.Second)
	// 	index(2)
	// 	time.Sleep(2 * time.Second)
	// 	index(3)
	// 	time.Sleep(2 * time.Second)
	// 	index(4)
	// 	time.Sleep(2 * time.Second)
	// 	index(5)
	// 	time.Sleep(2 * time.Second)
	// 	index(6)
	// }
	// if screen == 7 {
	// 	m.dev.SetBacklightLegcay(60)
	// }
	return nil
}

func (m *display) mainScreen() error {
	if m.screenActual != 1 {
		if err := m.switchScreen(1, true); err != nil {
			return err
		}
	}

	m.ingresos(m.inputsCash, m.inputsApp, 0)
	m.route(m.routeName)
	m.driver(m.driverID)
	m.counters(m.counterInput, 0, m.counterOutput, 0)
	m.gpsstate(m.gpsState)
	m.netstate(m.netState)
	m.updateDate(0)
	return nil
}

func (m *display) disableSelectButton() {
	setState := m.dev.SetPropertyValueU8(int(buttonSelectPaso), gtt43a.ButtonState)
	setState(0x00)
}

func (m *display) updateDate(period int) error {

	if m.screenActual != 1 {
		return nil
	}
	tNow := time.Now()
	if m.lastUpdateDate.Minute() == tNow.Minute() && m.lastUpdateDate.Hour() == tNow.Hour() {
		return nil
	}
	m.lastUpdateDate = tNow
	textTime := m.dev.SetPropertyText(addrTimeDate, gtt43a.LabelText)
	if err := textTime(tNow.Format("2006/01/02 15:04")); err != nil {
		return err
	}
	return nil
}

func (m *display) reset() {
	m.dev.Reset()
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

func (m *display) doors(value [2]int) error {

	return nil
}

func (m *display) recorridoPercent(value int) error {
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

func (m *display) keyNum(text string) {
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

func (m *display) inputValue(initialText string, screen int) {
	screen2label := map[int]int{2: 52, 3: 48}
	fmt.Printf("current screen: %d\n", m.screenActual)
	if m.screenActual == screen {
		return
	}
	if err := m.switchScreen(screen, true); err != nil {
		logs.LogWarn.Println(err)
		return
	}
	if v, ok := screen2label[screen]; ok && len(initialText) > 0 {
		if err := m.dev.SetPropertyText(v, gtt43a.LabelText)(initialText); err != nil {
			logs.LogWarn.Println(err)
		}
	}
}

func (m *display) driver(driver string) error {
	if len(driver) <= 0 {
		return nil
	}

	m.driverID = driver
	if m.screenActual != 1 {
		return nil
	}
	text := m.dev.SetPropertyText(addrIDDriver,
		gtt43a.LabelText)

	if err := text(driver); err != nil {
		return err
	}

	return nil
}

func (m *display) counters(inputs1, outputs1, inputs2, outputs2 int64) error {
	/**/
	m.counterInput = inputs1 + inputs2
	m.counterInput = outputs1 + outputs2

	if m.screenActual != 1 {
		return nil
	}
	textCounterInputs := m.dev.SetPropertyText(int(addrCounterInputs), gtt43a.LabelText)
	if err := textCounterInputs(fmt.Sprintf("%d", inputs1+inputs2)); err != nil {
		return err
	}

	textCounterOutputs := m.dev.SetPropertyText(int(addrCounterOutputs), gtt43a.LabelText)
	if err := textCounterOutputs(fmt.Sprintf("%d", outputs1+outputs2)); err != nil {
		return err
	}
	return nil
}

func (m *display) gpsstate(state int) error {
	m.gpsState = state
	index := m.dev.SetPropertyValueU16(addrIconGPS, gtt43a.VisualBitmap_SourceIndex)
	if state == 0 {
		if err := index(0); err != nil {
			return err
		}
	} else {
		if err := index(1); err != nil {
			return err
		}
	}
	return nil

}

func (m *display) netstate(state int) error {
	m.netState = state
	index := m.dev.SetPropertyValueU16(addrIconNET, gtt43a.VisualBitmap_SourceIndex)
	if state == 0 {
		if err := index(0); err != nil {
			return err
		}
	} else {
		if err := index(1); err != nil {
			return err
		}
	}
	return nil
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
	// fmt.Printf("notifs: %v\n", m.notifications)
	return nil
}

func (m *display) shownotifications() error {
	if m.screenActual == 7 {
		return nil
	}
	if err := m.switchScreen(7, true); err != nil {
		return err
	}
	text := m.dev.SetPropertyText(addrTextAlarms, gtt43a.LabelText)
	arrayText := make([]string, 0)
	for _, v := range m.notifications {
		arrayText = append(arrayText, fmt.Sprintf("==> %s", v))
	}
	if err := text(strings.Join(arrayText, "\n")); err != nil {
		log.Println(err)
	}
	return nil
}

func (m *display) setBrightness(percent int) error {
	data := percent * 255 / 100
	return m.dev.SetBacklightLegcay(data)
}

func (m *display) eventCount(input, output int) error {

	m.counterInput += int64(input)
	m.counterOutput += int64(output)

	if m.screenActual != 1 {
		return nil
	}
	textCounterInputs := m.dev.SetPropertyText(int(addrCounterInputs), gtt43a.LabelText)
	if err := textCounterInputs(fmt.Sprintf("%d", m.counterInput)); err != nil {
		return err
	}

	textCounterOutputs := m.dev.SetPropertyText(int(addrCounterOutputs), gtt43a.LabelText)
	if err := textCounterOutputs(fmt.Sprintf("%d", m.counterOutput)); err != nil {
		return err
	}
	return nil
}
