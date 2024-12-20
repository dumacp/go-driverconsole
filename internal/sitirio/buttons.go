package app

type Event struct {
	Label EventLabel
	Value interface{}
}

type EventLabel int

const (
	PROGRAMATION_DRIVER EventLabel = iota
	PROGRAMATION_VEH
	STATS
	SHOW_NOTIF
	ROUTE
	DRIVER
	SCREEN_SWITCH
	STEP_ENABLE
	STEP_APPLY
	BRIGHT_ADD
	BRIGHT_SUB
	ROUTE_BUTTON
	DRIVER_BUTTON
)
