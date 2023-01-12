package display

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

type Display interface {
	init() error
	close()
	mainScreen() error
	screenError(sError ...string) error
	textError(sError ...string) error
	textConfirmation(sError ...string) error
	textConfirmationMainScreen(timeout time.Duration, sError ...string) error
	warningInMainScreen(timeout time.Duration, sText ...string) error
	ingresos(usoEfectivo, usoCivica, usoParcial int) error
	ingresosPartial(usoParcial int)
	selectionRuta()
	updateRuta(ruta, subruta string)
	alertBeep(repeat int)
	verifyReset(chan int, actor.Context)
	timeRecorrido(int)
	disableSelectButton()
	updateDate(int) error
	switchScreen(int, bool) error
	screen() int
	reset()
	setPuerta(int, int)
	textInput(text ...string)
	addTextInput(text string)
	delTextInput(int)
	clearTextInput()
	keyNum(text string)
	inputValue(initialText string, screen int)
	doors(value [2]int) error
	recorridoPercent(int) error
	route(string) error
	driver(string) error
	counters(inputs1, outputs1, inputs2, outputs2 int64) error
	eventCount(input, output int) error
	gpsstate(state int) error
	netstate(state int) error
	addnotification(msg string) error
	shownotifications() error
	setBrightness(percent int) error
}
