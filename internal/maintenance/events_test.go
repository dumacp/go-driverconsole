package maintenance

import (
	"errors"
	"testing"
	"time"
)

func code(t *testing.T, m *Monitor, err error, now time.Time) string {
	t.Helper()
	evt := m.Update(err, now)
	if evt == nil {
		return ""
	}
	if evt.Type != EVENT_TYPE {
		t.Fatalf("type = %q, want %q", evt.Type, EVENT_TYPE)
	}
	return evt.Value.(*Value).Code
}

func TestMonitor(t *testing.T) {
	errHMI := errors.New("device not ready")
	t0 := time.Date(2026, 9, 24, 8, 0, 0, 0, time.UTC)
	step := 10 * time.Second

	tests := []struct {
		name  string
		steps []struct {
			err  error
			at   time.Duration
			want string
		}
	}{
		{
			name: "sin fallos no reporta",
			steps: []struct {
				err  error
				at   time.Duration
				want string
			}{
				{nil, 0, ""},
				{nil, step, ""},
			},
		},
		{
			name: "primer fallo, persistente y recuperación",
			steps: []struct {
				err  error
				at   time.Duration
				want string
			}{
				{errHMI, 0, CODE_DISCONNECTED},
				{errHMI, step, ""},
				{errHMI, 14 * time.Minute, ""},
				{errHMI, 15 * time.Minute, CODE_DISCONNECTED},
				{errHMI, 16 * time.Minute, ""},
				{nil, 17 * time.Minute, CODE_CONNECTED},
				{nil, 18 * time.Minute, ""},
			},
		},
		{
			name: "fallos repetidos dentro del intervalo no se reportan",
			steps: []struct {
				err  error
				at   time.Duration
				want string
			}{
				{errHMI, 0, CODE_DISCONNECTED},
				{nil, time.Minute, CODE_CONNECTED},
				{errHMI, 2 * time.Minute, ""},
				{nil, 3 * time.Minute, ""},
				{errHMI, 5 * time.Minute, ""},
				{errHMI, 15 * time.Minute, CODE_DISCONNECTED},
				{nil, 16 * time.Minute, CODE_CONNECTED},
			},
		},
		{
			name: "fallo repetido después del intervalo se reporta",
			steps: []struct {
				err  error
				at   time.Duration
				want string
			}{
				{errHMI, 0, CODE_DISCONNECTED},
				{nil, time.Minute, CODE_CONNECTED},
				{errHMI, 20 * time.Minute, CODE_DISCONNECTED},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMonitor(DEV_HMI, 15*time.Minute)
			for i, s := range tt.steps {
				if got := code(t, m, s.err, t0.Add(s.at)); got != s.want {
					t.Errorf("step %d (%s): got %q, want %q", i, s.at, got, s.want)
				}
			}
		})
	}
}

func TestMonitorCountIncreases(t *testing.T) {
	// el contador creciente evita que go-gwiot descarte el reporte persistente como duplicado
	m := NewMonitor(DEV_HMI, 15*time.Minute)
	t0 := time.Now()
	first := m.Update(errors.New("x"), t0).Value.(*Value)
	for i := 1; i < 90; i++ {
		m.Update(errors.New("x"), t0.Add(time.Duration(i)*10*time.Second))
	}
	second := m.Update(errors.New("x"), t0.Add(15*time.Minute)).Value.(*Value)
	if second.Count <= first.Count || second.Since != first.Since {
		t.Errorf("first=%+v second=%+v", first, second)
	}
}
