package app

import (
	"reflect"

	"github.com/dumacp/go-schservices/api/services"
)

// func UpdateService(prev, current *services.ScheduleService) *services.ScheduleService {
// 	for i := 0; i < reflect.TypeOf(*prev).NumField(); i++ {
// 		// field := reflect.TypeOf(prev).Field(i)
// 		value1 := reflect.ValueOf(*prev).Field(i)
// 		value2 := reflect.ValueOf(*current).Field(i)

// 		if value2.IsZero() {
// 			continue
// 		}
// 		value1.Set(value2.Elem())
// 	}
// 	return prev
// }

func UpdateService(prev, current *services.ScheduleService) *services.ScheduleService {
	sourceValue := reflect.ValueOf(current).Elem()
	destinationValue := reflect.ValueOf(prev).Elem()

	for i := 0; i < sourceValue.NumField(); i++ {
		sourceFieldValue := sourceValue.Field(i)
		destinationFieldValue := destinationValue.Field(i)

		if destinationFieldValue.CanSet() {
			// fmt.Println("can set")
			// Check if the destination field is different from zero

			if sourceFieldValue.IsZero() {
				// if reflect.DeepEqual(destinationFieldValue.Interface(), reflect.Zero(destinationFieldValue.Type()).Interface()) {
				// fmt.Println("skip")
				continue // Skip updating if it's not zero
			}

			// fmt.Println("set", sourceFieldValue)
			// Update the destination field with the value from the source
			destinationFieldValue.Set(sourceFieldValue)
		}
	}
	return prev
}
