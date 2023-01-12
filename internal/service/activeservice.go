package service

import (
	"time"
)

type ActiveService struct {
	BeginTime   time.Time
	StartedTime time.Time
	EndTime     time.Time
	StoppedTime time.Time
	Driver      string
	Itinerary   string
}
