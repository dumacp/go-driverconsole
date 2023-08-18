package device

import (
	"time"

	"github.com/dumacp/matrixorbital/gtt43a"
)

type devGtt50 struct {
	port  string
	speed int
	dev   interface{}
}

func (d *devGtt50) Init() (interface{}, error) {
	opts := &gtt43a.PortOptions{}
	opts.Baud = d.speed
	opts.Port = d.port
	opts.ReadTimeout = 100 * time.Millisecond
	dev := gtt43a.NewDisplay(opts)

	if err := dev.Open(); err != nil {
		return nil, err
	}
	return dev, nil
}

func (d *devGtt50) Close() error {
	if v, ok := d.dev.(gtt43a.Display); ok {
		return v.Close()
	}
	return nil
}

func NewGtt50Device(port string, speed int) Device {
	dev := &devGtt50{
		port:  port,
		speed: speed,
	}
	return dev
}
