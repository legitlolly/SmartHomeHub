package simulator

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/legitlolly/SmartHomeHub/internal/device"
)

type SimulatedDevice struct {
	id         device.ID
	state      device.State
	stateMutex sync.RWMutex // Doesn't block on reads like Mutex does.
}

func NewSimulatedDevice(id device.ID) *SimulatedDevice {
	return &SimulatedDevice{
		id: id,
		state: device.State{
			Power:      "off",
			Brightness: 0,
			UpdatedAt:  time.Now(),
		},
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
		d.state.Power = "on"
	case "turn_off":
		d.state.Power = "off"
	case "set_brightness":
		val, ok := cmd.Params["value"].(int)
		if !ok {
			return errors.New("invalid value for brightness")
		}
		if val < 0 || val > 100 {
			return errors.New("brightness must be 0-100")
		}
		d.state.Brightness = val
	default:
		return errors.New("unknown command")
	}

	d.state.UpdatedAt = time.Now()
	return nil
}

func (d *SimulatedDevice) State(ctx context.Context) (device.State, error) {
	d.stateMutex.RLock()
	defer d.stateMutex.RUnlock()

	return device.State{
		Power:      d.state.Power,
		Brightness: d.state.Brightness,
		UpdatedAt:  d.state.UpdatedAt,
	}, nil
}
