package service

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/dumacp/go-driverconsole/internal/itinerary"
	"github.com/dumacp/go-driverconsole/internal/utils"
	"github.com/golang/geo/s2"
)

type actorFollowService struct {
	iti         itinerary.Itinerary
	maxDistance int
}

func NewActorFollowService(iti itinerary.Itinerary, initialPoint s2.Point, maxDistance int) actor.Actor {

	a := &actorFollowService{}
	a.iti = iti
	a.maxDistance = maxDistance
	return a
}

func (a *actorFollowService) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:

	case s2.Point:
		point := msg
		pPoint, _ := a.iti.ProjectPoint(point)

		a := pPoint.Distance(point)
		distance := utils.AngleToMeters(a)

		fmt.Printf("distance: %d", distance)
	}
}
