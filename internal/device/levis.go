package device

import "github.com/dumacp/go-levis"

type devPi struct {
	port  string
	speed int
	dev   interface{}
}

func (d *devPi) Init() (interface{}, error) {
	dev, err := levis.NewDevice(d.port, d.speed)
	if err != nil {
		return nil, err
	}

	return dev, nil
}

func (d *devPi) Close() error {
	if v, ok := d.dev.(levis.Device); ok {
		return v.Close()
	}
	return nil
}

func NewPiDevice(port string, speed int) Device {
	dev := &devPi{
		port:  port,
		speed: speed,
	}

	return dev
}
