package buttons

import (
	"context"
)

type InputEvent struct {
	TypeEvent TypeEvent
	KeyCode   KeyCode
	Value     interface{}
	Error     error
}

type ButtonDevice interface {
	Init(dev interface{}) error
	Close() error
	ListenButtons(ctx context.Context) (<-chan *InputEvent, error)
}

type KeyCode int

type TypeEvent int

const (
	ButtonEvent TypeEvent = iota
	DataEvent
)

type ButtonValue int

const (
	TextNumRoute  = "textNumRoute"
	TextNumDriver = "textNumDriver"
)

const (
	GridButton0 ButtonValue = iota
	GridButton1
	GridButton2
	GridButton3
	GridButton4
	GridButton5
	GridButton6
	GridButton7
	GridButton8
	GridButton9
	GridButtonEnter
	GridButtonDel
	UpVoc
	EnterVoc
	SelectPASO
	EnterPASO
	ResetCounter
	ResetRecorrido
)

type StopButtons struct{}
type StartButtons struct{}
