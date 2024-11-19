package ui

// type Screen int

const (
	// MAIN          Screen = 0
	// DRIVER_SCREEN Screen = 1
	// VEHICLE       Screen = 2
	// ALARM         Screen = 3
	// SERVICE       Screen = 4
	// CONFIRMATION  Screen = 5
	// WARNING       Screen = 6
	// ERROR         Screen = 7
	// CHECK         Screen = 8

	MAIN_SCREEN int = iota
	PROGRAMATION_DRIVER_SCREEN
	PROGRAMATION_VEH_SCREEN
	ALARMS_SCREEN
	ADDITIONALS_SCREEN
	NOTIFICATIONS_SCREEN
	INFO_SCREEN
	WARN_SCREEN
	ERROR_SCREEN
	CHECK_SCREEN
	KEY_ROUTE_SCREEN
	KEY_DRIVER_SCREEN
)
