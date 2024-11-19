package ignition

import "time"

type State int

const (
	NA State = iota
	UP
	DOWN
)

func (s State) String() string {
	switch s {
	case DOWN:
		return "DOWN"
	case UP:
		return "UP"
	}

	return "NA"
}

type IgnitionEvent struct {
	StateEvent State
	Timestamp  time.Time
}
