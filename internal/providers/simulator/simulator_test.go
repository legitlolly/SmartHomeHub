package simulator

import (
	"context"
	"testing"

	"github.com/legitlolly/SmartHomeHub/internal/device"
)

func TestSimulatedDevice_ExecuteandState(t *testing.T) {
	ctx := context.Background()
	testDev1 := NewSimulatedDevice("test-light-1")

	// Test turning on command
	cmdOn := device.Command{
		DeviceID: testDev1.ID(),
		Action:   "turn_on",
	}

	if err := testDev1.Execute(ctx, cmdOn); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	state, _ := testDev1.State(ctx)
	if state.Power != "on" {
		t.Errorf("Test failed -> expected on, got %s", state.Power)
	}

	cmdBrightness := device.Command{
		DeviceID: testDev1.ID(),
		Action:   "set_brightness",
		Params:   map[string]any{"value": 50},
	}

	if err := testDev1.Execute(ctx, cmdBrightness); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	state, _ = testDev1.State(ctx)
	if state.Brightness != 50 {
		t.Errorf("Test failed -> expected brightness=50, got %d", state.Brightness)
	}
}

func TestSimulatedDevice_InvalidCommand(t *testing.T) {
	ctx := context.Background()
	dev := NewSimulatedDevice("light-2")

	cmd := device.Command{
		DeviceID: dev.ID(),
		Action:   "invalid",
	}

	if err := dev.Execute(ctx, cmd); err == nil {
		t.Fatal("expected error for invalid command, got nil")
	}
}
