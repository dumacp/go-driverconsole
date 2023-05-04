package ui

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/dumacp/go-driverconsole/internal/display"
)

func AckResponse(action *actor.Future) error {
	res, err := action.Result()
	if err != nil {
		return err
	}
	switch v := res.(type) {
	case error:
		return err
	case *display.AckMsg:
		return v.Error
	default:
		return nil
	}
}
