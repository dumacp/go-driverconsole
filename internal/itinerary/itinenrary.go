package itinerary

import (
	"github.com/dumacp/go-driverconsole/internal/utils"
	"github.com/golang/geo/s2"
)

type ControlPoint struct {
	Name     string
	Type     string
	Radius   int
	MaxSpeed int
	ETA      int
	idx      int
	Point    s2.Point
}

func (g ControlPoint) IsInto(point s2.Point) bool {
	distance := utils.AngleToMeters(g.Point.Distance(point))
	return distance <= g.Radius
}

type Itinerary struct {
	ID             string
	OrganizationID string
	Route          string
	Polyline       s2.Polyline
	ControlPoints  []ControlPoint
	PaymentID      int
}

func (iti Itinerary) ControlPoint(point s2.Point) (ControlPoint, bool) {
	for _, v := range iti.ControlPoints {
		if v.IsInto(point) {
			return v, true
		}
	}
	return ControlPoint{}, false
}

func (iti Itinerary) ProjectPoint(point s2.Point) (s2.Point, float32) {
	len := len(iti.Polyline)
	p, i := iti.Polyline.Project(point)

	return p, float32(i) / float32(len)
}
