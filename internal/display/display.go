package display

type Display interface {
	init()
	close()
	mainScreen()
	screenError(sError ...string)
	textError(sError ...string)
	textConfirmation(sError ...string)
	ingresos(usoEfectivo, usoCivica, usoParcial int)
	ingresosPartial(usoParcial int)
	selectionRuta()
	updateRuta(ruta, subruta string)
	alertBeep(repeat int)
	verifyReset() chan int
	// listenButtons() chan Button
	timeRecorrido(int)
	disableSelectButton()
	updateDate(int)
	switchScreen(int, bool) error
	screen() int
	reset()
	setPuerta(int, int)
	textInput(text ...string)
	addTextInput(text string)
	delTextInput(int)
	clearTextInput()
	keyNum(text string)

	doors(value [2]int) error
	recorridoPercent(int) error
	route(string) error
}
