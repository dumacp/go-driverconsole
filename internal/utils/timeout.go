package utils

import "time"

// timeout function to compare current;y trigger with the previous execution
// ti is initial time (for previous trigger)
// seq is sequency of timeouts to await in consecutive triggers. IFf len(seq) == 0
// then only one timeout is 0.
func Timeout(ti time.Time, seq []time.Duration) func() bool {
	mem := ti
	idx := 0
	var seqm []time.Duration
	if len(seq) <= 0 {
		seqm = []time.Duration{0}
	} else {
		seqm = seq
	}
	return func() bool {
		if idx >= len(seqm) {
			idx = len(seqm) - 1
		}
		duration := seqm[idx]
		result := mem.Add(duration).Before(time.Now())
		mem = time.Now()
		idx++
		return result
	}
}
