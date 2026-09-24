package maintenance

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dumacp/go-driverconsole/internal/pubsub"
	"github.com/dumacp/go-gwiot/pkg/events"
	"github.com/dumacp/go-logs/pkg/logs"
)

const (
	TOPIC      = "EVENTS/driverconsole"
	EVENT_TYPE = "DRIVERCONSOLE_EVT"
)

// dispositivos que generan eventos de mantenimiento
const (
	DEV_HMI = "HMI"
)

// códigos de evento
const (
	CODE_DISCONNECTED = "DISCONNECTED"
	CODE_CONNECTED    = "CONNECTED"
)

// HMIReportInterval intervalo mínimo entre reportes de fallo de la HMI
var HMIReportInterval = 15 * time.Minute

// Value contenido de un evento DRIVERCONSOLE_EVT
type Value struct {
	Dev   string  `json:"dev"`
	Code  string  `json:"code"`
	State bool    `json:"state"`
	Error string  `json:"error,omitempty"`
	Since float64 `json:"since,omitempty"` // inicio del fallo (segundos epoch)
	Count int     `json:"count,omitempty"` // verificaciones fallidas desde el inicio del fallo
}

// NewEvent construye un evento DRIVERCONSOLE_EVT
func NewEvent(t time.Time, value *Value) *events.Event {
	return &events.Event{
		Timestamp: float64(t.UnixMilli()) / 1000,
		Type:      EVENT_TYPE,
		Value:     value,
	}
}

// Publish envía el evento al topic de eventos (procesado por go-gwiot)
func Publish(evt *events.Event) error {
	data, err := json.Marshal(evt)
	if err != nil {
		return fmt.Errorf("marshal %s error: %w", evt.Type, err)
	}
	logs.LogInfo.Printf("maintenance event: %s", data)
	pubsub.Publish(TOPIC, data)
	return nil
}

// Monitor detecta fallos de un dispositivo y decide cuándo reportarlos:
// el primer fallo se reporta de inmediato; fallos posteriores (persistentes
// o repetidos) solo si pasó "interval" desde el último reporte de fallo.
// La recuperación se reporta solo si el fallo en curso fue reportado.
type Monitor struct {
	dev        string
	interval   time.Duration
	failing    bool
	reported   bool
	since      time.Time
	count      int
	lastReport time.Time
}

func NewMonitor(dev string, interval time.Duration) *Monitor {
	return &Monitor{
		dev:      dev,
		interval: interval,
	}
}

// Update recibe el resultado de una verificación y devuelve el evento a
// publicar, o nil si no hay nada que reportar
func (m *Monitor) Update(err error, now time.Time) *events.Event {
	if err != nil {
		if !m.failing {
			m.failing = true
			m.since = now
			m.count = 0
		}
		m.count++
		if !m.lastReport.IsZero() && now.Sub(m.lastReport) < m.interval {
			return nil
		}
		m.lastReport = now
		m.reported = true
		return NewEvent(now, &Value{
			Dev:   m.dev,
			Code:  CODE_DISCONNECTED,
			State: false,
			Error: err.Error(),
			Since: float64(m.since.UnixMilli()) / 1000,
			Count: m.count,
		})
	}
	if !m.failing {
		return nil
	}
	m.failing = false
	if !m.reported {
		return nil
	}
	m.reported = false
	return NewEvent(now, &Value{
		Dev:   m.dev,
		Code:  CODE_CONNECTED,
		State: true,
		Since: float64(m.since.UnixMilli()) / 1000,
	})
}
