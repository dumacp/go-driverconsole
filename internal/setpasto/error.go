package app

import "time"

type errorDisplay struct {
	Timestamp time.Time
	Error     error
}
