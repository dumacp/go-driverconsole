//go:build levis
// +build levis

package device

import "github.com/dumacp/go-levis"

type dev struct {
	port  string
	speed int
	dev   interface{}
}

func (d *dev) Init() (interface{}, error) {
	dev, err := levis.NewDevice(d.port, d.speed)
	if err != nil {
		return nil, err
	}

	return dev, nil
}

func (d *dev) Close() error {
	if v, ok := d.dev.(levis.Device); ok {
		return v.Close()
	}
	return nil
}

func NewPiDevice(port string, speed int) Device {
	dev := &dev{
		port:  port,
		speed: speed,
	}

	return dev
}
