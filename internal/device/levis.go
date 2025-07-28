package device

import (
	"fmt"

	"github.com/dumacp/go-levis"
)

type devPi struct {
	port  string
	speed int
	dev   interface{}
}

func (d *devPi) Init() (interface{}, error) {
	dev, err := levis.NewDevice(d.port, d.speed)
	if err != nil {
		fmt.Printf("error opening Pi device on port %s: %s\n", d.port, err)
		return nil, err
	}

	// verify readiness
	if _, err := dev.ReadRegister(0, 1); err != nil {
		dev.Close()
		speed := d.speed
		switch speed {
		case 115200:
			speed = 38400 // try with 38400 if 115200 fails
		case 38400:
			speed = 115200 // try with 115200 if 38400 fails
		default:
			return nil, fmt.Errorf("device not ready on port %s: %s", d.port, err)
		}
		devv, err := levis.NewDevice(d.port, speed)
		if err != nil {
			fmt.Printf("error opening Pi device on port %s: %s\n", d.port, err)
			return nil, err
		}
		if _, err := devv.ReadRegister(0, 1); err != nil {
			devv.Close()
			return nil, fmt.Errorf("device not ready on port %s: %s", d.port, err)
		} else {
			dev = devv
		}
		fmt.Printf("device ready on port %s with speed %d\n", d.port, speed)
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
