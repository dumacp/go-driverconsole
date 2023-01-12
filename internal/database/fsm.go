package database

import (
	"fmt"
	"time"

	"github.com/dumacp/go-logs/pkg/logs"
	"github.com/looplab/fsm"
	"go.etcd.io/bbolt"
)

const (
	sInit      = "sInit"
	sOpen      = "sOpen"
	sClose     = "sClose"
	sRestart   = "sRestart"
	sWaitEvent = "sWait"
)

const (
	eOpenCmd = "eOpenCmd"
	eOpened  = "eOpened"
	eClosed  = "eClosed"
	eError   = "eError"
	eStart   = "eStart"
)

// func beforeEvent(event string) string {
// 	return fmt.Sprintf("before_%s", event)
// }

func enterState(state string) string {
	return fmt.Sprintf("enter_%s", state)
}

// func leaveState(state string) string {
// 	return fmt.Sprintf("leave_%s", state)
// }

func (a *dbActor) initFSM() *fsm.FSM {

	f := fsm.NewFSM(
		sInit,
		fsm.Events{
			{Name: eOpenCmd, Src: []string{sInit, sClose, sRestart}, Dst: sOpen},
			{Name: eOpened, Src: []string{sOpen}, Dst: sWaitEvent},
			{Name: eClosed, Src: []string{sOpen, sWaitEvent}, Dst: sClose},
			{Name: eError, Src: []string{sOpen, sWaitEvent}, Dst: sRestart},
		},
		fsm.Callbacks{
			"enter_state": func(e *fsm.Event) {
				logs.LogBuild.Printf("FSM DB state Src: %v, state Dst: %v", e.Src, e.Dst)
			},
			"leave_state": func(e *fsm.Event) {
				if e.Err != nil {
					e.Cancel(e.Err)
				}
			},
			"before_event": func(e *fsm.Event) {
				if e.Err != nil {
					e.Cancel(e.Err)
				}
			},
			enterState(sOpen): func(e *fsm.Event) {
				if a.db != nil {
					a.db.Close()
				}
				db, err := bbolt.Open(a.pathDB, 0666, nil)
				if err != nil {
					logs.LogError.Println(err)
					a.ctx.Send(a.ctx.Self(), &MsgErrorDB{})
					e.Err = err
					return
				}
				a.db = db
				a.ctx.Send(a.ctx.Self(), &MsgOpenedDB{})
			},
			enterState(sClose): func(e *fsm.Event) {
				if a.db != nil {
					a.db.Close()
					a.db = nil
				}
				a.behavior.Become(a.CloseState)
			},
			enterState(sRestart): func(e *fsm.Event) {
				time.Sleep(30 * time.Second)
				a.behavior.Become(a.CloseState)
				a.ctx.Send(a.ctx.Self(), &MsgOpenDB{})
			},
		})

	return f
}
