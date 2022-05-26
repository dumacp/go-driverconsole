package app

import "github.com/dumacp/matrixorbital/gtt43a"

type DisplayDevice struct {
	Device gtt43a.Display
}
type InputText struct {
	Text string
}

type SelectPaso struct{}
type EnterPaso struct{}
type ResetCounter struct{}

type MsgAppPaso struct {
	Value int
}
type MsgAppPercentRecorrido struct {
	Data int
}
type MsgMainScreen struct{}
type MsgScreen struct {
	ID     int
	Switch bool
}
type MsgInitRecorrido struct{}
type MsgStopRecorrido struct{}
type ResetRecorrido struct{}
type MsgSubscribe struct{}
type MsgChangeRoute struct {
	ID int
}
type MsgChangeDriver struct {
	ID int
}

type MsgCounters struct {
	Parcial  int
	Efectivo int
	App      int
}

type MsgRoute struct {
	Route string
}
type MsgDriver struct {
	Driver string
}

type MsgSetRoutes struct {
	Routes map[int]string
}

type MsgDoors struct {
	Value [2]int
}

type MsgConfirmationText struct {
	Text []byte
}
type MsgConfirmationButton struct{}

type MsgWarningText struct {
	Text []byte
}
type MsgWarningButton struct{}
