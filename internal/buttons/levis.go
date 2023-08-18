package buttons

import (
	"context"
	"fmt"
	"time"

	"github.com/dumacp/go-levis"
	"github.com/dumacp/go-logs/pkg/logs"
)

const (
	// AddrSelectPaso  = 0
	// AddrEnterPaso   = 1
	// AddrEnterRuta   = 2
	// AddrEnterDriver = 3

	// AddrScreenSwitch     = 4
	// AddrScreenAlarms     = 7
	// AddrScreenProgVeh    = 8
	// AddrScreenProgDriver = 9
	// AddrScreenMore       = 10
	AddrReset = 20
	AddrBeep  = 23
	// AddrAddBright        = 21
	// AddrSubBright        = 22
)

type pi3070g struct {
	listenStart int
	listenEnd   int
	buttons     []int
	dev         levis.Device
}

func NewConfPiButtons(startAddrToListen, endAddrToListen int, buttonsAddress []int) ButtonDevice {
	pi := &pi3070g{
		listenStart: startAddrToListen,
		listenEnd:   endAddrToListen,
		buttons:     buttonsAddress,
	}
	return pi
}

func (p *pi3070g) Init(dev interface{}) error {
	pi, ok := dev.(levis.Device)
	if !ok {
		var ii levis.Device
		return fmt.Errorf("device is not %T", ii)
	}
	p.dev = pi

	p.dev.Conf().SetButtonMem(p.listenStart, p.listenEnd)

	for _, v := range p.buttons {
		if err := p.dev.AddButton(v); err != nil {
			return err
		}
	}

	fmt.Printf("//////////////// conf buttons: %+v\n", p)

	return nil
}

func (p *pi3070g) Close() error {

	if p.dev != nil {
		return p.dev.Close()
	}
	return nil
}

// addrNoRoute    = 120
// addrNameRoute  = 100
// addrNoDriver   = 160
// addrNameDriver = 140

func (p *pi3070g) ListenButtons(contxt context.Context) (<-chan *InputEvent, error) {

	if p.dev == nil {
		return nil, fmt.Errorf("device is not iniatilize")
	}
	ch := p.dev.ListenButtonsWithContext(contxt)
	chEvt := make(chan *InputEvent, 1)

	go func() {
		defer close(chEvt)
		// lastStep := time.Now()
		// enableStep := time.NewTimer(5 * time.Second)
		// activeStep := false

		for {
			select {
			case <-contxt.Done():
				logs.LogWarn.Println("ListenButtons context is closed")
				return
			// case <-enableStep.C:
			// 	diff := time.Since(lastStep)
			// 	if diff < 3*time.Second && activeStep {
			// 		enableStep.Reset(diff)
			// 		break
			// 	}
			// 	if activeStep {
			// 		fmt.Println("reset addrSelectPaso")
			// 		activeStep = false
			// 		if err := p.dev.SetIndicator(addr, false); err != nil {
			// 			fmt.Println(err)
			// 		}
			// 	}
			case button, ok := <-ch:
				// if button.Value == 0 {
				// 	break
				// }
				// if err := p.dev.SetIndicator(button.Addr, false); err != nil {
				// 	fmt.Println(err)
				// }
				if !ok {
					evt := &InputEvent{
						Error: fmt.Errorf("device closed"),
					}
					select {
					case chEvt <- evt:
					case <-time.After(100 * time.Millisecond):
					}
					return
				}
				evt := &InputEvent{
					TypeEvent: ButtonEvent,
					KeyCode:   KeyCode(button.Addr),
					Value:     button.Value == 0,
					Error:     nil,
				}
				select {
				case chEvt <- evt:
				case <-time.After(100 * time.Millisecond):
				}
			}
		}
	}()
	return chEvt, nil
}
