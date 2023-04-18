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
			logs.LogBuild.Printf("FSM DEVICE, state src: %s, state dst: %s", e.Src, e.Dst)
		},
		beforeEvent(eStarted): func(e *fsm.Event) {
			var err error

			var devi interface{}
			for _, v := range []int{0, 3, 3, 10, 30, 60} {
				if v > 0 {
					time.Sleep(time.Duration(v) * time.Second)
				}
				devi, err = a.dev.Init()
				if err == nil {
					break
				}
				fmt.Printf("open device error: %s\n", err)
				a.dev.Close()
			}
			if err != nil {
				e.Cancel(fmt.Errorf("%w", err))
				return
			}
			a.ctx.Send(a.ctx.Self(), &MsgDevice{Device: devi})
		},
		enterState(sClose): func(e *fsm.Event) {
			if a.dev == nil {
				return
			}
			//dev, ok := disp.(Device)
			//if ok {
			if err := a.dev.Close(); err != nil {
				fmt.Printf("error close device = %s\n", err)
			}
			fmt.Println("close device")
			//}
		},
	}

	f := fsm.NewFSM(
		sStart,
		fsm.Events{
			{
				Name: eStarted,
				Src:  []string{sStart, sClose, sStop, sWait},
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
