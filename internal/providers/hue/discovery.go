package hue

import (
	"context"
	"fmt"
	"log"

	"github.com/amimof/huego"
	"github.com/legitlolly/SmartHomeHub/internal/device"
)

// DiscoverAndRegisterLights discovers all lights on the bridge and registers them
func DiscoverAndRegisterLights(ctx context.Context, registry Registry, ip, username string) error {
	bridge := huego.New(ip, username)

	lights, err := bridge.GetLightsContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to discover lights: %w", MapHueError(err))
	}

	if len(lights) == 0 {
		log.Println("No Hue lights found on bridge")
		return nil
	}

	log.Printf("Discovered %d Hue light(s):", len(lights))
	for _, light := range lights {
		deviceID := device.ID(fmt.Sprintf("hue-light-%d", light.ID))

		log.Printf("  - Light %d: %s (Model: %s)", light.ID, light.Name, light.ModelID)

		bridgeClient := NewHuegoBridge(ip, username)

		hueDevice := NewHueDevice(deviceID, light.ID, bridgeClient)
		if err := registry.Register(hueDevice); err != nil {
			log.Printf("    Failed to register %s: %v", deviceID, err)
			continue
		}
		log.Printf("    Registered as: %s", deviceID)
	}

	return nil
}

type Registry interface {
	Register(dev device.Device) error
}
