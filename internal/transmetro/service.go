package app

import (
	"reflect"
	"time"

	"github.com/dumacp/go-schservices/api/services"
)

func UpdateServiceStable(prev, current *services.ScheduleService) *services.ScheduleService {
	if current == nil && prev != nil {
		return prev
	}
	if prev == nil && current != nil {
		return current
	}
	if current != nil && prev != nil && current.Id != prev.Id {
		return prev
	}

	if current.CheckpointTimingState != nil {
		prev.CheckpointTimingState = current.CheckpointTimingState
	}
	if current.Driver != nil {
		prev.Driver = current.Driver
	}
	if current.DriverIds != nil {
		prev.DriverIds = current.DriverIds
	}
	if current.Itinerary != nil {
		prev.Itinerary = current.Itinerary
	}
	if current.OrganizationId != "" {
		prev.OrganizationId = current.OrganizationId
	}
	if current.Route != nil {
		prev.Route = current.Route
	}
	if current.ScheduleDateTime > 0 {
		prev.ScheduleDateTime = current.ScheduleDateTime
	}
	if current.State != "" {
		prev.State = current.State
	}

	return prev
}

func UpdateService(prev, current *services.ScheduleService) *services.ScheduleService {
	if current == nil && prev != nil {
		return prev
	}
	if prev == nil && current != nil {
		return current
	}
	if current != nil && prev != nil && current.Id != prev.Id {
		return prev
	}
	sourceValue := reflect.ValueOf(current).Elem()
	destinationValue := reflect.ValueOf(prev).Elem()

	zeroTime := time.UnixMilli(0)

	for i := 0; i < sourceValue.NumField(); i++ {
		sourceFieldValue := sourceValue.Field(i)
		destinationFieldValue := destinationValue.Field(i)

		if destinationFieldValue.CanSet() {
			// fmt.Println("can set")
			// Check if the destination field is different from zero

			if sourceFieldValue.Type() == reflect.TypeOf(time.Time{}) {
				t := sourceFieldValue.Interface().(time.Time)
				if t.IsZero() || t.Equal(zeroTime) {
					continue // Omitir si es un valor considerado "cero"
				}
			} else {
				if sourceFieldValue.IsZero() {
					continue // Skip updating if it's not zero
				}
			}

			// fmt.Println("set", sourceFieldValue)
			// Update the destination field with the value from the source
			destinationFieldValue.Set(sourceFieldValue)
		}
	}
	return prev
}
