//go:build (gtt50 || !levis) && (gtt50 || !gtt43)
// +build gtt50 !levis
// +build gtt50 !gtt43

package device

import (
	"time"

	"github.com/dumacp/matrixorbital/gtt43a"
)

func NewDevice(port string, speed int) (Device, error) {

	opts := &gtt43a.PortOptions{}
	opts.Baud = speed
	opts.Port = port
	opts.ReadTimeout = 30 * time.Millisecond
	dev := gtt43a.NewDisplay(opts)

	if err := dev.Open(); err != nil {
		return nil, err
	}

	return dev, nil
}
