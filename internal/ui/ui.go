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
	lastUpdateDate  time.Time
	screen          Screen
	disp            display.Display
	chEvents        chan *Event
	pid             *actor.PID
	rootctx         *actor.RootContext
	button2LabelApp map[buttons.KeyCode]EventType
}

type UI interface {
	Init() error
	MainScreen() error
	TextWarning(text ...string) error
	TextConfirmation(text ...string) error
	TextConfirmationPopup(timeout time.Duration, text ...string) error
	TextWarningPopup(timeout time.Duration, sText ...string) error
	Inputs(in int) error
	Outputs(out int) error
	DeviationInputs(dev int) error
	Route(route ...string) error
	Driver(data string) error
	Beep(repeat int, timeout time.Duration) error
	Date(date time.Time) error
	Screen(num int, force bool) error
	GetScreen() (int, error)
	KeyNum(ctx context.Context, prompt string) (chan int, error)
	Keyboard(ctx context.Context, prompt string) (chan string, error)
	Doors(state ...bool) error
	Gps(state bool) error
	Network(state bool) error
	AddNotifications(add string) error
	ShowNotifications() error
	ShowProgDriver() error
	ShowProgVeh() error
	ShowStats() error
	Brightness(percent int) error
	ServiceCurrentState(state int, prompt string) error
	InputHandler(inputs actor.Actor, button2LabelApp map[buttons.KeyCode]EventType) error
	Events() chan *Event
}

func New(ctx actor.Context, dev, disp actor.Actor) (UI, error) {

	props := actor.PropsFromFunc(NewActor(dev, disp).Receive)
	pid, err := ctx.SpawnNamed(props, "ui-actor")
	if err != nil {
		return nil, fmt.Errorf("init UI actor error: %s", err)
	}

	u := &ui{}
	u.pid = pid
	u.rootctx = ctx.ActorSystem().Root

	return u, nil
}

func (u *ui) Init() error {

	if err := u.MainScreen(); err != nil {
		return err
	}

	return nil
}

func (u *ui) Events() chan *Event {
	return u.chEvents
}

func (u *ui) InputHandler(inputs actor.Actor, button2LabelApp map[buttons.KeyCode]EventType) error {
	u.button2LabelApp = button2LabelApp
	res, err := u.rootctx.RequestFuture(u.pid, &AddInputsHandlerMsg{Handler: inputs}, 3*time.Second).Result()
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
	res, err := u.rootctx.RequestFuture(u.pid, &MainScreenMsg{}, 3*time.Second).Result()
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
	res, err := u.rootctx.RequestFuture(u.pid, &TextWarningMsg{Text: text}, 3*time.Second).Result()
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
	res, err := u.rootctx.RequestFuture(u.pid, &TextConfirmationMsg{Text: text}, 3*time.Second).Result()
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

func (u *ui) Inputs(in int) error {
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

func (u *ui) Outputs(out int) error {
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

func (u *ui) DeviationInputs(dev int) error {
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

func (u *ui) Route(route ...string) error {
	res, err := u.rootctx.RequestFuture(u.pid, &RouteMsg{Route: route}, 3*time.Second).Result()
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
	res, err := u.rootctx.RequestFuture(u.pid, &DriverMsg{Data: data}, 3*time.Second).Result()
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

func (u *ui) Beep(repeat int, timeout time.Duration) error {
	return nil
}

func (u *ui) Date(date time.Time) error {

	res, err := u.rootctx.RequestFuture(u.pid, &DateMsg{Date: date}, 3*time.Second).Result()
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
	panic("not implemented") // TODO: Implement
}

func (u *ui) Network(state bool) error {
	panic("not implemented") // TODO: Implement
}

func (u *ui) AddNotifications(add string) error {
	panic("not implemented") // TODO: Implement
}

func (u *ui) ShowNotifications() error {
	return nil
}

func (u *ui) ShowProgVeh() error {
	return nil
}

func (u *ui) ShowProgDriver() error {
	return nil
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
