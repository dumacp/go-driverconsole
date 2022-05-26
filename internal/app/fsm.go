package app

import (
	"github.com/dumacp/go-logs/pkg/logs"
	"github.com/looplab/fsm"
)

const (
	sStart = "sStart"
	sOpen  = "sOpen"
	sClose = "sClose"
	sWait  = "sWait"
)

const (
	eStarted = "eStarted"
	eOpenned = "eOpenned"
	eClosed  = "eClosed"
	eError   = "eError"
)

func Fsm() *fsm.FSM {

	callbacksfsm := fsm.Callbacks{
		"before_event": func(e *fsm.Event) {
			if e.Err != nil {
				// log.Println(e.Err)
				e.Cancel(e.Err)
			}
		},
		"leave_state": func(e *fsm.Event) {
			if e.Err != nil {
				// log.Println(e.Err)
				e.Cancel(e.Err)
			}
		},
		"enter_state": func(e *fsm.Event) {
			logs.LogBuild.Printf("FSM APP, state src: %s, state dst: %s", e.Src, e.Dst)
		},
	}

	f := fsm.NewFSM(
		sStart,
		fsm.Events{
			{
				Name: eStarted,
				Src:  []string{sStart},
				Dst:  sOpen,
			},
			{
				Name: eOpenned,
				Src:  []string{sOpen},
				Dst:  sWait,
			},
			{
				Name: eError,
				Src:  []string{sStart, sOpen, sWait},
				Dst:  sClose,
			},
			{
				Name: eClosed,
				Src:  []string{sClose},
				Dst:  sOpen,
			},
		},
		callbacksfsm,
	)

	return f
}
