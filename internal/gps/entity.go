package gps

import (
	"time"

	"github.com/golang/geo/s2"
)

type GPSData struct {
	Validity bool
	LatLon   s2.LatLng
	HDop     float32
	Speed    float32
	Time     time.Time
}
