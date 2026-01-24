package device_test

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/legitlolly/SmartHomeHub/internal/device"
	"github.com/legitlolly/SmartHomeHub/internal/providers/simulator"
)

func TestRegistry_RegisterAndGet(t *testing.T) {
	testRegistry := device.NewRegistry()
	testDev1 := simulator.NewSimulatedDevice("test-light-1")
	if err := testRegistry.Register(testDev1); err != nil {
		t.Fatalf("Failed to register device: %v", err)
	}

	if dList := testRegistry.List(); len(dList) != 1 {
		t.Fatalf("Registry is not as expected, recieved : %v", dList)
	}

	retrieved, err := testRegistry.Get(testDev1.ID())
	if err != nil {
		t.Fatalf("Device not found: %v", err)
	}

	if retrieved.ID() != testDev1.ID() {
		t.Fatalf("Expected device ID %s, got %s", testDev1.ID(), retrieved.ID())
	}

	testRegistry.Unregister(testDev1.ID())

	if _, err := testRegistry.Get(testDev1.ID()); err == nil {
		t.Fatalf("Device should not exist after deletion")
	}

}

func TestRegistry_Unregister(t *testing.T) {
	testRegistry := device.NewRegistry()
	testDev1 := simulator.NewSimulatedDevice("test-light-1")

	if err := testRegistry.Register(testDev1); err != nil {
		t.Fatalf("Failed to register device: %v", err)
	}
	testRegistry.Unregister(testDev1.ID())

	_, err := testRegistry.Get(testDev1.ID())
	if err != device.ErrDeviceNotFound {
		t.Fatalf("Expected ErrDeviceNotFound, got %v", err)
	}
}

func TestRegistry_GetNonExistent(t *testing.T) {
	testRegistry := device.NewRegistry()

	_, err := testRegistry.Get("non-existent-device-id")
	if err != device.ErrDeviceNotFound {
		t.Fatalf("Expected ErrDeviceNotFound, got %v", err)
	}
}

func TestRegistry_DuplicateDevice(t *testing.T) {
	testRegistry := device.NewRegistry()
	testDev1 := simulator.NewSimulatedDevice("test-light-1")

	if err := testRegistry.Register(testDev1); err != nil {
		t.Fatalf("Failed to register device: %v", err)
	}

	if err := testRegistry.Register(testDev1); err != device.ErrDeviceAlreadyRegistered {
		t.Fatalf("Duplicate device registered expected ErrDeviceAlreadyRegistered")
	}
}

func TestRegistry_UpdateDevice(t *testing.T) {
	ctx := context.Background()
	testRegistry := device.NewRegistry()
	testDev1 := simulator.NewSimulatedDevice("test-light-1")

	if err := testRegistry.Register(testDev1); err != nil {
		t.Fatalf("Failed to register device: %v", err)
	}

	cmdBrightness := device.Command{
		DeviceID: testDev1.ID(),
		Action:   "set_brightness",
		Params:   map[string]any{"value": 50},
	}

	if err := testDev1.Execute(ctx, cmdBrightness); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	retrieved, err := testRegistry.Get(testDev1.ID())

	if err != nil {
		t.Fatalf("Expected to find device instead got: %v", err)
	}

	state, _ := retrieved.State(ctx)

	if state.Attributes["brightness"] != 50 {
		t.Fatalf("Stored device brightness is incorrect, expected 50 got %d", state.Attributes["brightness"])
	}

}

func TestRegistry_MultipleDevices(t *testing.T) {
	testRegistry := device.NewRegistry()
	testDev1 := simulator.NewSimulatedDevice("test-light-1")
	testDev2 := simulator.NewSimulatedDevice("test-light-2")

	testRegistry.Register(testDev1)
	testRegistry.Register(testDev2)

	devices := testRegistry.List()
	if len(devices) != 2 {
		t.Fatalf("Expected 2 devices, got %d", len(devices))
	}
}

func TestRegistry_ConcurrentAccess(t *testing.T) {
	testRegistry := device.NewRegistry()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			dev := simulator.NewSimulatedDevice(device.ID(fmt.Sprintf("light-%d", id)))
			testRegistry.Register(dev)
		}(i)
	}
	wg.Wait()

	if len(testRegistry.List()) != 100 {
		t.Fatalf("Expected 100 devices after concurrent registration. Got %d", len(testRegistry.List()))
	}
}
