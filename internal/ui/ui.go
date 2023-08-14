package ui

import (
	"context"
	"fmt"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/dumacp/go-driverconsole/internal/buttons"
	"github.com/dumacp/go-driverconsole/internal/display"
)

type ui struct {
	lastUpdateDate time.Time
	screen         Screen
	disp           display.Display
	notif          []string
	// chEvents       chan *Event
	pid     *actor.PID
	rootctx *actor.RootContext
}

type UI interface {
	Init() error
	MainScreen() error
	TextWarning(text ...string) error
	TextConfirmation(text ...string) error
	TextConfirmationPopup(timeout time.Duration, text ...string) error
	TextWarningPopup(timeout time.Duration, sText ...string) error
	Inputs(in int32) error
	Outputs(out int32) error
	CashInputs(in int32) error
	ElectronicInputs(out int32) error
	DeviationInputs(dev int32) error
	Route(route ...string) error
	Driver(data string) error
	Beep(repeat, duty int, period time.Duration) error
	Date(date time.Time) error
	Screen(num int, force bool) error
	GetScreen() (int, error)
	KeyNum(ctx context.Context, prompt string) (chan int, error)
	Keyboard(ctx context.Context, prompt string) (chan string, error)
	Doors(state ...bool) error
	Gps(state bool) error
	StepEnable(state bool) error
	Network(state bool) error
	AddNotifications(add string) error
	ShowNotifications(...string) error
	ShowProgDriver(...string) error
	ShowProgVeh(...string) error
	ShowStats() error
	Brightness(percent int) error
	ServiceCurrentState(state int, prompt string) error
	InputHandler(inputs actor.Actor, callback func(evt *buttons.InputEvent)) error
	ReadBytesRawDisplay(label int) ([]byte, error)
	// Events() chan *Event
}

func New(ctx actor.Context, dev, disp actor.Actor) (UI, error) {

	props := actor.PropsFromFunc(NewActor(dev, disp).Receive)
	pid, err := ctx.SpawnNamed(props, "ui-actor")
	if err != nil {
		return nil, fmt.Errorf("init UI actor error: %s", err)
	}

	u := &ui{}
	u.notif = make([]string, 0)
	u.pid = pid
	u.rootctx = ctx.ActorSystem().Root

	return u, nil
}

func (u *ui) Init() error {

	res, err := u.rootctx.RequestFuture(u.pid, &InitUIMsg{}, 1*time.Second).Result()
	if err != nil {
		return err
	}
	if v, ok := res.(*AckMsg); ok && v.Error != nil {
		return v.Error
	} else if ok {
		return nil
	}
	return fmt.Errorf("init with response form display")

}

// func (u *ui) Events() chan *Event {
// 	return u.chEvents
// }

func (u *ui) InputHandler(inputs actor.Actor, callback func(evt *buttons.InputEvent)) error {
	res, err := u.rootctx.RequestFuture(u.pid, &AddInputsHandlerMsg{Handler: inputs, Evt2Func: callback}, 2*time.Second).Result()
	if err != nil {
		return err
	}
	if v, ok := res.(*AckMsg); ok && v.Error != nil {
		return v.Error
	} else if ok {
		return nil
	}
	return fmt.Errorf("inputHandler without response form display")
}

func (u *ui) MainScreen() error {
	res, err := u.rootctx.RequestFuture(u.pid, &MainScreenMsg{}, 2*time.Second).Result()
	if err != nil {
		return err
	}
	if v, ok := res.(*AckMsg); ok && v.Error != nil {
		return v.Error
	} else if ok {
		u.screen = MAIN
		return nil
	}
	return fmt.Errorf("mainScreen without response form display")
}

func (u *ui) TextWarning(text ...string) error {
	res, err := u.rootctx.RequestFuture(u.pid, &TextWarningMsg{Text: text}, 2*time.Second).Result()
	if err != nil {
		return err
	}
	if v, ok := res.(*AckMsg); ok && v.Error != nil {
		return v.Error
	} else if ok {
		return nil
	}
	return fmt.Errorf("textWarning with response form display")
}

func (u *ui) TextConfirmation(text ...string) error {
	res, err := u.rootctx.RequestFuture(u.pid, &TextConfirmationMsg{Text: text}, 2*time.Second).Result()
	if err != nil {
		return err
	}
	if v, ok := res.(*AckMsg); ok && v.Error != nil {
		return v.Error
	} else if ok {
		return nil
	}
	return fmt.Errorf("textConfirmation with response form display")
}

func (u *ui) TextConfirmationPopup(timeout time.Duration, sText ...string) error {
	res, err := u.rootctx.RequestFuture(u.pid, &TextConfirmationPopupMsg{
		Text:    sText,
		Timeout: timeout,
	}, 3*time.Second).Result()
	if err != nil {
		return err
	}
	if v, ok := res.(*AckMsg); ok && v.Error != nil {
		return v.Error
	} else if ok {
		return nil
	}
	return fmt.Errorf("textConfirmationPopup with response form display")
}

func (u *ui) TextWarningPopup(timeout time.Duration, sText ...string) error {
	res, err := u.rootctx.RequestFuture(u.pid, &TextWarningPopupMsg{
		Timeout: timeout,
		Text:    sText,
	}, 3*time.Second).Result()
	if err != nil {
		return err
	}
	if v, ok := res.(*AckMsg); ok && v.Error != nil {
		return v.Error
	} else if ok {
		return nil
	}
	return fmt.Errorf("textWarningPopup with response form display")
}

func (u *ui) Inputs(in int32) error {
	res, err := u.rootctx.RequestFuture(u.pid, &InputsMsg{
		In: in,
	}, 3*time.Second).Result()
	if err != nil {
		return err
	}
	if v, ok := res.(*AckMsg); ok && v.Error != nil {
		return v.Error
	} else if ok {
		return nil
	}
	return fmt.Errorf("inputs with response form display")
}

func (u *ui) Outputs(out int32) error {
	res, err := u.rootctx.RequestFuture(u.pid, &OutputsMsg{
		Out: out,
	}, 3*time.Second).Result()
	if err != nil {
		return err
	}
	if v, ok := res.(*AckMsg); ok && v.Error != nil {
		return v.Error
	} else if ok {
		return nil
	}
	return fmt.Errorf("outputs with response form display")
}

func (u *ui) DeviationInputs(dev int32) error {
	res, err := u.rootctx.RequestFuture(u.pid, &DeviationInputsMsg{
		Dev: dev,
	}, 3*time.Second).Result()
	if err != nil {
		return err
	}
	if v, ok := res.(*AckMsg); ok && v.Error != nil {
		return v.Error
	} else if ok {
		return nil
	}
	return fmt.Errorf("deviationInputs with response form display")
}

func (u *ui) CashInputs(in int32) error {
	res, err := u.rootctx.RequestFuture(u.pid, &CashInputsMsg{
		In: in,
	}, 3*time.Second).Result()
	if err != nil {
		return err
	}
	if v, ok := res.(*AckMsg); ok && v.Error != nil {
		return v.Error
	} else if ok {
		return nil
	}
	return fmt.Errorf("inputs with response form display")
}

func (u *ui) ElectronicInputs(in int32) error {
	res, err := u.rootctx.RequestFuture(u.pid, &ElectronicInputsMsg{
		In: in,
	}, 3*time.Second).Result()
	if err != nil {
		return err
	}
	if v, ok := res.(*AckMsg); ok && v.Error != nil {
		return v.Error
	} else if ok {
		return nil
	}
	return fmt.Errorf("inputs with response form display")
}

func (u *ui) Route(route ...string) error {
	res, err := u.rootctx.RequestFuture(u.pid, &RouteMsg{Route: route}, 2*time.Second).Result()
	if err != nil {
		return err
	}
	if v, ok := res.(*AckMsg); ok && v.Error != nil {
		return v.Error
	} else if ok {
		return nil
	}
	return fmt.Errorf("route without response from display")
}

func (u *ui) Driver(data string) error {
	res, err := u.rootctx.RequestFuture(u.pid, &DriverMsg{Data: data}, 2*time.Second).Result()
	if err != nil {
		return err
	}
	if v, ok := res.(*AckMsg); ok && v.Error != nil {
		return v.Error
	} else if ok {
		return nil
	}
	return fmt.Errorf("driver without response from display")
}

func (u *ui) Beep(repeat, duty int, period time.Duration) error {
	res, err := u.rootctx.RequestFuture(u.pid, &BeepMsg{Repeat: repeat, Duty: duty, Period: period}, 2*time.Second).Result()
	if err != nil {
		return err
	}
	if v, ok := res.(*AckMsg); ok && v.Error != nil {
		return v.Error
	} else if ok {
		return nil
	}
	return fmt.Errorf("beep without response from display")
}

func (u *ui) Date(date time.Time) error {

	res, err := u.rootctx.RequestFuture(u.pid, &DateMsg{Date: date}, 2*time.Second).Result()
	if err != nil {
		return err
	}
	if v, ok := res.(*AckMsg); ok && v.Error != nil {
		return v.Error
	} else if ok {
		return nil
	}
	return fmt.Errorf("date without response from display")
}

func (u *ui) Screen(num int, force bool) error {
	res, err := u.rootctx.RequestFuture(u.pid, &ScreenMsg{
		Num: num,
	}, 3*time.Second).Result()
	if err != nil {
		return err
	}
	if v, ok := res.(*AckMsg); ok && v.Error != nil {
		return v.Error
	} else if ok {
		return nil
	}
	return fmt.Errorf("screen without response from display")
}

func (u *ui) GetScreen() (int, error) {
	panic("not implemented") // TODO: Implement
}

func (u *ui) KeyNum(ctx context.Context, prompt string) (chan int, error) {
	panic("not implemented") // TODO: Implement
}

func (u *ui) Keyboard(ctx context.Context, prompt string) (chan string, error) {
	panic("not implemented") // TODO: Implement
}

func (u *ui) Doors(state ...bool) error {
	panic("not implemented") // TODO: Implement
}

func (u *ui) Gps(state bool) error {
	res, err := u.rootctx.RequestFuture(u.pid, &GpsMsg{State: state}, 1*time.Second).Result()
	if err != nil {
		return err
	}
	if v, ok := res.(*AckMsg); ok && v.Error != nil {
		return v.Error
	} else if ok {
		return nil
	}
	return fmt.Errorf("network without response from display")
}

func (u *ui) Network(state bool) error {
	res, err := u.rootctx.RequestFuture(u.pid, &NetworkMsg{State: state}, 1*time.Second).Result()
	if err != nil {
		return err
	}
	if v, ok := res.(*AckMsg); ok && v.Error != nil {
		return v.Error
	} else if ok {
		return nil
	}
	return fmt.Errorf("network without response from display")
}

func (u *ui) AddNotifications(add string) error {
	u.notif = append(u.notif, add)
	if len(u.notif) > 10 {
		copy(u.notif, u.notif[1:])
		u.notif = u.notif[:len(u.notif)-1]
	}
	return nil
}

func (u *ui) ShowNotifications(data ...string) error {
	res, err := u.rootctx.RequestFuture(u.pid, &ShowNotificationsMsg{Text: data}, 3*time.Second).Result()
	if err != nil {
		return err
	}
	if v, ok := res.(*AckMsg); ok && v.Error != nil {
		return v.Error
	} else if ok {
		return nil
	}
	return fmt.Errorf("showNotifications without response from display")
}

func (u *ui) ShowProgVeh(data ...string) error {
	res, err := u.rootctx.RequestFuture(u.pid, &ShowProgVehMsg{Text: data}, 3*time.Second).Result()
	if err != nil {
		return err
	}
	if v, ok := res.(*AckMsg); ok && v.Error != nil {
		return v.Error
	} else if ok {
		return nil
	}
	return fmt.Errorf("showProgVeh without response from display")
}

func (u *ui) ShowProgDriver(data ...string) error {
	res, err := u.rootctx.RequestFuture(u.pid, &ShowProgDriverMsg{Text: data}, 3*time.Second).Result()
	if err != nil {
		return err
	}
	if v, ok := res.(*AckMsg); ok && v.Error != nil {
		return v.Error
	} else if ok {
		return nil
	}
	return fmt.Errorf("showProgDriver without response from display")
}

func (u *ui) ShowStats() error {
	return nil
}

func (u *ui) Brightness(percent int) error {
	panic("not implemented") // TODO: Implement
}

func (u *ui) ServiceCurrentState(state int, prompt string) error {
	res, err := u.rootctx.RequestFuture(u.pid, &ServiceCurrentStateMsg{
		State:  state,
		Prompt: prompt,
	}, 2*time.Second).Result()
	if err != nil {
		return err
	}
	if v, ok := res.(*AckMsg); ok && v.Error != nil {
		return v.Error
	} else if ok {
		return nil
	}
	return fmt.Errorf("deviationInputs with response form display")
}

func (u *ui) ReadBytesRawDisplay(label int) ([]byte, error) {
	res, err := u.rootctx.RequestFuture(u.pid, &ReadBytesRawMsg{Label: label}, 2*time.Second).Result()
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*ReadBytesRawResponseMsg); ok && v.Error != nil {
		return nil, v.Error
	} else if ok {
		data := make([]byte, len(v.Value))
		copy(data, v.Value)
		return data, nil
	}
	return nil, fmt.Errorf("textWarning with response form display")
}

func (u *ui) StepEnable(state bool) error {
	res, err := u.rootctx.RequestFuture(u.pid, &StepEnableMsg{State: state}, 2*time.Second).Result()
	if err != nil {
		return err
	}
	if v, ok := res.(*AckMsg); ok && v.Error != nil {
		return v.Error
	} else if ok {
		return nil
	}
	return fmt.Errorf("stepEnable without response from display")
}
