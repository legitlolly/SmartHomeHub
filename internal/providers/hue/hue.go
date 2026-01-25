package hue

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/amimof/huego"
	"github.com/legitlolly/SmartHomeHub/internal/device"
)

type HueDevice struct {
	id      device.ID
	lightID int          // Hue bridge light ID
	client  BridgeClient // Interface for testing

	// Cached state
	lastState  *cachedState
	stateMutex sync.RWMutex
}

type cachedState struct {
	power      string
	brightness int
	hue        *int
	saturation *int
	colorTemp  *int
	updatedAt  time.Time
}

func NewHueDevice(id device.ID, lightID int, client BridgeClient) *HueDevice {
	return &HueDevice{
		id:      id,
		lightID: lightID,
		client:  client,
		lastState: &cachedState{
			power:      "unknown",
			brightness: 0,
			updatedAt:  time.Now(),
		},
	}
}

func (d *HueDevice) ID() device.ID {
	return d.id
}

func (d *HueDevice) Execute(ctx context.Context, cmd device.Command) error {
	var state huego.State

	switch cmd.Action {
	case "turn_on":
		state.On = true

	case "turn_off":
		state.On = false

	case "set_brightness":
		var brightness int
		switch v := cmd.Params["value"].(type) {
		case float64:
			brightness = int(v)
		case int:
			brightness = v
		default:
			return fmt.Errorf("%w: brightness value must be a number", ErrInvalidParameter)
		}

		if brightness < 0 || brightness > 100 {
			return fmt.Errorf("%w: brightness must be 0-100", ErrInvalidParameter)
		}

		state.Bri = uint8((brightness * 254) / 100)

		if brightness > 0 {
			state.On = true
		}

	default:
		return fmt.Errorf("unknown command: %s", cmd.Action)
	}

	_, err := d.client.SetLightStateContext(ctx, d.lightID, state)
	if err != nil {
		return MapHueError(err)
	}

	d.updateCacheAfterCommand(cmd)

	return nil
}

func (d *HueDevice) State(ctx context.Context) (device.State, error) {
	light, err := d.client.GetLightContext(ctx, d.lightID)
	if err != nil {
		d.stateMutex.RLock()
		defer d.stateMutex.RUnlock()

		return device.State{
			DeviceType: "light",
			UpdatedAt:  d.lastState.updatedAt,
			Attributes: map[string]interface{}{
				"power":      "unknown",
				"brightness": 0,
				"error":      MapHueError(err).Error(),
			},
		}, MapHueError(err)
	}

	d.stateMutex.Lock()
	defer d.stateMutex.Unlock()

	if light.State.On {
		d.lastState.power = "on"
	} else {
		d.lastState.power = "off"
	}

	d.lastState.brightness = int(light.State.Bri) * 100 / 254

	if light.State.Hue > 0 {
		hue := int(light.State.Hue)
		d.lastState.hue = &hue
	}
	if light.State.Sat > 0 {
		sat := int(light.State.Sat)
		d.lastState.saturation = &sat
	}
	if light.State.Ct > 0 {
		ct := int(light.State.Ct)
		d.lastState.colorTemp = &ct
	}

	d.lastState.updatedAt = time.Now()

	attributes := map[string]interface{}{
		"power":      d.lastState.power,
		"brightness": d.lastState.brightness,
		"model":      light.ModelID,
		"name":       light.Name,
	}

	if d.lastState.hue != nil {
		attributes["hue"] = *d.lastState.hue
	}
	if d.lastState.saturation != nil {
		attributes["saturation"] = *d.lastState.saturation
	}
	if d.lastState.colorTemp != nil {
		attributes["color_temperature"] = *d.lastState.colorTemp
	}

	return device.State{
		DeviceType: "light",
		UpdatedAt:  d.lastState.updatedAt,
		Attributes: attributes,
	}, nil
}

func (d *HueDevice) updateCacheAfterCommand(cmd device.Command) {
	d.stateMutex.Lock()
	defer d.stateMutex.Unlock()

	switch cmd.Action {
	case "turn_on":
		d.lastState.power = "on"
	case "turn_off":
		d.lastState.power = "off"
	case "set_brightness":
		if val, ok := cmd.Params["value"].(int); ok {
			d.lastState.brightness = val
			if val > 0 {
				d.lastState.power = "on"
			}
		} else if val, ok := cmd.Params["value"].(float64); ok {
			d.lastState.brightness = int(val)
			if val > 0 {
				d.lastState.power = "on"
			}
		}
	}

	d.lastState.updatedAt = time.Now()
}
