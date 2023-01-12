package platform

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/golang/geo/s2"
	"go.etcd.io/bbolt"
)

const (
	ItinerariosName = "ITINERARIES"
)

type ItineraryPathPoint struct {
	Coords   []float64 `json:"coords"`
	Name     string    `json:"name"`
	Radius   int       `json:"radius"`
	Type     string    `json:"type"`
	MaxSpeed int       `json:"maxSpeed"`
	ETA      int       `json:"eta"`
}
type ItineraryMetadata struct {
	UpdateAt int64 `json:"updatedAt"`
}
type RouteMngItinerary struct {
	Active         bool                  `json:"active"`
	ID             string                `json:"id"`
	Name           string                `json:"name"`
	OrganizationID string                `json:"organizationId"`
	RouteID        string                `json:"routeId"`
	Metadata       *ItineraryMetadata    `json:"metadata"`
	Path           []*ItineraryPathPoint `json:"path"`
}
type DataItinerary struct {
	Data *RouteMngItinerary `json:"RouteMngItinerary"`
}

func (r *RouteMngItinerary) Coords() [][]float64 {
	coords := make([][]float64, 0)
	for _, v := range r.Path {
		coords = append(coords, v.Coords)
	}
	return coords
}

func (r *RouteMngItinerary) Points() []s2.Point {
	coords := make([]s2.Point, 0)
	for _, v := range r.Path {
		if len(v.Coords) < 2 {
			continue
		}
		lalon := s2.LatLngFromDegrees(v.Coords[1], v.Coords[0])
		coords = append(coords, s2.PointFromLatLng(lalon))
	}
	return coords
}

func (r *RouteMngItinerary) ControlPoints() []s2.Point {
	coords := make([]s2.Point, 0)
	for _, v := range r.Path {
		for _, subs := range []string{"Stop", "Control"} {
			if strings.Contains(v.Type, subs) {
				if len(v.Coords) < 2 {
					continue
				}
				lalon := s2.LatLngFromDegrees(v.Coords[1], v.Coords[0])
				coords = append(coords, s2.PointFromLatLng(lalon))
				continue
			}
		}

	}
	return coords
}

func (r *RouteMngItinerary) GetPolyline() s2.Polyline {
	points := make([]s2.LatLng, 0)
	for _, v := range r.Path {
		coord := v.Coords
		if len(coord) >= 2 {
			lalon := s2.LatLngFromDegrees(coord[1], coord[0])
			points = append(points, lalon)
		}
	}
	return *s2.PolylineFromLatLngs(points)
}

func (r *RouteMngItinerary) Save(db *bbolt.DB) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}
	if err := db.Update(func(tx *bbolt.Tx) error {
		bk, err := tx.CreateBucketIfNotExists([]byte(ItinerariosName))
		if err != nil {
			return err
		}
		data, err := json.Marshal(r)
		if err != nil {
			return err
		}
		if err := bk.Put([]byte(r.ID), data); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (r *RouteMngItinerary) Delete(db *bbolt.DB) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}
	if err := db.Update(func(tx *bbolt.Tx) error {
		bk, err := tx.CreateBucketIfNotExists([]byte(ItinerariosName))
		if err != nil {
			return err
		}
		if err := bk.Delete([]byte(r.ID)); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
