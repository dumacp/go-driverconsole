package device

import (
	"fmt"
	"time"

	"github.com/dumacp/go-logs/pkg/logs"
	"github.com/looplab/fsm"
)

const (
	sStart = "sStart"
	sOpen  = "sOpen"
	sClose = "sClose"
	sWait  = "sWait"
	sStop  = "sStop"
)

const (
	eStarted = "eStarted"
	eOpenned = "eOpenned"
	eClosed  = "eClosed"
	eError   = "eError"
	eStop    = "eStop"
)

func beforeEvent(event string) string {
	return fmt.Sprintf("before_%s", event)
}
func enterState(state string) string {
	return fmt.Sprintf("enter_%s", state)
}
func leaveState(state string) string {
	return fmt.Sprintf("leave_%s", state)
}

func (a *Actor) Fsm() {

	var disp interface{}
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
		beforeEvent(eStarted): func(e *fsm.Event) {
			var err error

			for range []int{0, 1, 2} {
				disp, err = NewDevice(a.portSerial, a.speedBaud)
				if err == nil {
					break
				}
				dev, ok := disp.(Device)
				if ok {
					dev.Close()
				}
				time.Sleep(3 * time.Second)
			}
			if err != nil {
				e.Cancel(err)
				return
			}
			a.ctx.Send(a.ctx.Self(), &MsgDevice{Device: disp})
		},
		enterState(sClose): func(e *fsm.Event) {
			if disp == nil {
				return
			}
			dev, ok := disp.(Device)
			if ok {
				dev.Close()
			}
		},
	}

	f := fsm.NewFSM(
		sStart,
		fsm.Events{
			{
				Name: eStarted,
				Src:  []string{sStart, sClose},
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
				Name: eStop,
				Src:  []string{sStart, sOpen, sWait},
				Dst:  sStop,
			},
		},
		callbacksfsm,
	)

	a.fmachinae = f
}
