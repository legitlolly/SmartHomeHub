package simulator

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/legitlolly/SmartHomeHub/internal/device"
)

type SimulatedDevice struct {
	id device.ID

	// Internal state - just the raw values
	power      string
	brightness int
	updatedAt  time.Time

	stateMutex sync.RWMutex
}

func NewSimulatedDevice(id device.ID) *SimulatedDevice {
	return &SimulatedDevice{
		id:         id,
		power:      "off", // Start off
		brightness: 0,
		updatedAt:  time.Now(),
	}
}

func (d *SimulatedDevice) ID() device.ID {
	return d.id
}

func (d *SimulatedDevice) Execute(ctx context.Context, cmd device.Command) error {
	d.stateMutex.Lock()
	defer d.stateMutex.Unlock()

	// Simulated latency
	time.Sleep(100 * time.Millisecond)

	switch cmd.Action {
	case "turn_on":
		d.power = "on"
	case "turn_off":
		d.power = "off"
	case "set_brightness":
		var brightness int
		switch v := cmd.Params["value"].(type) {
		case float64:
			brightness = int(v)
		case int:
			brightness = v
		default:
			return errors.New("Invalid value for brightness")
		}
		if brightness < 0 || brightness > 100 {
			return errors.New("Brightness must be 0-100")
		}
		d.brightness = brightness
	default:
		return errors.New("unknown command")
	}

	d.updatedAt = time.Now()
	return nil
}

func (d *SimulatedDevice) State(ctx context.Context) (device.State, error) {
	d.stateMutex.RLock()
	defer d.stateMutex.RUnlock()

	// Build the State struct from internal fields
	return device.State{
		DeviceType: "light",
		UpdatedAt:  d.updatedAt,
		Attributes: map[string]interface{}{
			"power":      d.power,
			"brightness": d.brightness,
		},
	}, nil
}
