package utils

import (
	"math"

	"github.com/golang/geo/s1"
)

// Equatorial Earth radius in km
const R = 6378.1

func AngleToMeters(a s1.Angle) int {

	l := math.Pi * R * a.Degrees() / 180
	return int(l * 1000)
}
