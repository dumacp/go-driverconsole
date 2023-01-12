package platform

import (
	"time"

	"github.com/golang/geo/s2"
)

// Route selected route
type RouteStarted struct {
	ID   string
	Name string
	//ItineraryIndex    int
	ItineraryID       string
	TimeStamp         time.Time
	StartedPoint      *s2.LatLng
	StartedPointIndex int
	Running           bool
	Online            bool
}
