package hue

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/amimof/huego"
	"github.com/legitlolly/SmartHomeHub/internal/device"
)

// MockBridgeClient implements BridgeClient for testing
type MockBridgeClient struct {
	lights      map[int]*huego.Light
	getError    error
	setError    error
	callHistory []string
}

func NewMockBridgeClient() *MockBridgeClient {
	return &MockBridgeClient{
		lights:      make(map[int]*huego.Light),
		callHistory: []string{},
	}
}

func (m *MockBridgeClient) GetLightContext(ctx context.Context, id int) (*huego.Light, error) {
	m.callHistory = append(m.callHistory, fmt.Sprintf("GetLight(%d)", id))

	if m.getError != nil {
		return nil, m.getError
	}

	light, ok := m.lights[id]
	if !ok {
		return nil, errors.New("light not found")
	}

	return light, nil
}

func (m *MockBridgeClient) SetLightStateContext(ctx context.Context, id int, state huego.State) (*huego.Response, error) {
	m.callHistory = append(m.callHistory, fmt.Sprintf("SetLightState(%d)", id))

	if m.setError != nil {
		return nil, m.setError
	}

	if light, ok := m.lights[id]; ok {
		light.State.On = state.On
		if state.Bri > 0 {
			light.State.Bri = state.Bri
		}
		if state.Hue > 0 {
			light.State.Hue = state.Hue
		}
	}

	return &huego.Response{}, nil
}

func (m *MockBridgeClient) AddLight(id int, name string, on bool, bri uint8) {
	m.lights[id] = &huego.Light{
		ID:      id,
		Name:    name,
		State:   &huego.State{On: on, Bri: bri},
		ModelID: "LCT001",
	}
}

func (m *MockBridgeClient) SimulateError(getErr, setErr error) {
	m.getError = getErr
	m.setError = setErr
}

func TestHueDevice_Execute_TurnOn(t *testing.T) {
	ctx := context.Background()
	mock := NewMockBridgeClient()
	mock.AddLight(1, "Test Light", false, 0)

	dev := NewHueDevice("test-hue-1", 1, mock)

	cmd := device.Command{
		DeviceID: dev.ID(),
		Action:   "turn_on",
	}

	err := dev.Execute(ctx, cmd)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	light, _ := mock.GetLightContext(ctx, 1)
	if !light.State.On {
		t.Error("Expected light to be on")
	}

	state, err := dev.State(ctx)
	if err != nil {
		t.Fatalf("State failed: %v", err)
	}
	if state.Attributes["power"] != "on" {
		t.Errorf("Expected power=on, got %v", state.Attributes["power"])
	}
}

func TestHueDevice_Execute_TurnOff(t *testing.T) {
	ctx := context.Background()
	mock := NewMockBridgeClient()
	mock.AddLight(1, "Test Light", true, 254)

	dev := NewHueDevice("test-hue-1", 1, mock)

	cmd := device.Command{
		DeviceID: dev.ID(),
		Action:   "turn_off",
	}

	err := dev.Execute(ctx, cmd)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	light, _ := mock.GetLightContext(ctx, 1)
	if light.State.On {
		t.Error("Expected light to be off")
	}

	state, err := dev.State(ctx)
	if err != nil {
		t.Fatalf("State failed: %v", err)
	}
	if state.Attributes["power"] != "off" {
		t.Errorf("Expected power=off, got %v", state.Attributes["power"])
	}
}

func TestHueDevice_Execute_SetBrightness(t *testing.T) {
	ctx := context.Background()
	mock := NewMockBridgeClient()
	mock.AddLight(1, "Test Light", true, 254)

	dev := NewHueDevice("test-hue-1", 1, mock)

	cmd := device.Command{
		DeviceID: dev.ID(),
		Action:   "set_brightness",
		Params:   map[string]any{"value": 50},
	}

	err := dev.Execute(ctx, cmd)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	light, _ := mock.GetLightContext(ctx, 1)
	expected := uint8(127) // (50 * 254) / 100
	if light.State.Bri != expected {
		t.Errorf("Expected Hue bri=%d, got %d", expected, light.State.Bri)
	}

	state, _ := dev.State(ctx)
	if state.Attributes["brightness"] != 50 {
		t.Errorf("Expected brightness=50, got %v", state.Attributes["brightness"])
	}
}

func TestHueDevice_Execute_SetBrightness_Float(t *testing.T) {
	ctx := context.Background()
	mock := NewMockBridgeClient()
	mock.AddLight(1, "Test Light", true, 127)

	dev := NewHueDevice("test-hue-1", 1, mock)

	cmd := device.Command{
		DeviceID: dev.ID(),
		Action:   "set_brightness",
		Params:   map[string]any{"value": float64(75)},
	}

	err := dev.Execute(ctx, cmd)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	light, _ := mock.GetLightContext(ctx, 1)
	expected := uint8(190) // (75 * 254) / 100
	if light.State.Bri != expected {
		t.Errorf("Expected Hue bri=%d, got %d", expected, light.State.Bri)
	}
}

func TestHueDevice_Execute_InvalidBrightness(t *testing.T) {
	ctx := context.Background()
	mock := NewMockBridgeClient()
	mock.AddLight(1, "Test Light", true, 127)

	dev := NewHueDevice("test-hue-1", 1, mock)

	cmd := device.Command{
		DeviceID: dev.ID(),
		Action:   "set_brightness",
		Params:   map[string]any{"value": 150}, // Invalid: > 100
	}

	err := dev.Execute(ctx, cmd)
	if err == nil {
		t.Fatal("Expected error for invalid brightness, got nil")
	}

	if !errors.Is(err, ErrInvalidParameter) {
		t.Errorf("Expected ErrInvalidParameter, got %v", err)
	}
}

func TestHueDevice_Execute_InvalidBrightnessType(t *testing.T) {
	ctx := context.Background()
	mock := NewMockBridgeClient()
	mock.AddLight(1, "Test Light", true, 127)

	dev := NewHueDevice("test-hue-1", 1, mock)

	cmd := device.Command{
		DeviceID: dev.ID(),
		Action:   "set_brightness",
		Params:   map[string]any{"value": "fifty"}, // Invalid: string instead of number
	}

	err := dev.Execute(ctx, cmd)
	if err == nil {
		t.Fatal("Expected error for invalid brightness type, got nil")
	}

	if !errors.Is(err, ErrInvalidParameter) {
		t.Errorf("Expected ErrInvalidParameter, got %v", err)
	}
}

func TestHueDevice_State_BridgeUnreachable(t *testing.T) {
	ctx := context.Background()
	mock := NewMockBridgeClient()
	mock.SimulateError(errors.New("connection timeout"), nil)

	dev := NewHueDevice("test-hue-1", 1, mock)

	state, err := dev.State(ctx)

	if err == nil {
		t.Fatal("Expected error when bridge unreachable")
	}

	if !errors.Is(err, ErrBridgeUnreachable) {
		t.Errorf("Expected ErrBridgeUnreachable, got %v", err)
	}

	if _, hasError := state.Attributes["error"]; !hasError {
		t.Error("Expected error attribute in state")
	}
}

func TestHueDevice_Execute_UnknownCommand(t *testing.T) {
	ctx := context.Background()
	mock := NewMockBridgeClient()
	mock.AddLight(1, "Test Light", false, 0)

	dev := NewHueDevice("test-hue-1", 1, mock)

	cmd := device.Command{
		DeviceID: dev.ID(),
		Action:   "invalid_action",
	}

	err := dev.Execute(ctx, cmd)
	if err == nil {
		t.Fatal("Expected error for unknown command, got nil")
	}
}

func TestHueDevice_ID(t *testing.T) {
	mock := NewMockBridgeClient()
	dev := NewHueDevice("custom-id-123", 5, mock)

	if dev.ID() != "custom-id-123" {
		t.Errorf("Expected ID=custom-id-123, got %s", dev.ID())
	}
}

func TestHueDevice_State_ExtendedAttributes(t *testing.T) {
	ctx := context.Background()
	mock := NewMockBridgeClient()

	hue := uint16(25500)
	sat := uint8(254)
	ct := uint16(300)
	on := true
	bri := uint8(200)

	mock.lights[1] = &huego.Light{
		ID:      1,
		Name:    "Color Light",
		ModelID: "LCT015",
		State: &huego.State{
			On:  on,
			Bri: bri,
			Hue: hue,
			Sat: sat,
			Ct:  ct,
		},
	}

	dev := NewHueDevice("color-light", 1, mock)
	state, err := dev.State(ctx)

	if err != nil {
		t.Fatalf("State failed: %v", err)
	}

	if state.Attributes["hue"] != int(hue) {
		t.Errorf("Expected hue=%d, got %v", hue, state.Attributes["hue"])
	}
	if state.Attributes["saturation"] != int(sat) {
		t.Errorf("Expected saturation=%d, got %v", sat, state.Attributes["saturation"])
	}
	if state.Attributes["color_temperature"] != int(ct) {
		t.Errorf("Expected color_temperature=%d, got %v", ct, state.Attributes["color_temperature"])
	}
	if state.Attributes["model"] != "LCT015" {
		t.Errorf("Expected model=LCT015, got %v", state.Attributes["model"])
	}
	if state.Attributes["name"] != "Color Light" {
		t.Errorf("Expected name=Color Light, got %v", state.Attributes["name"])
	}
}

func TestHueDevice_State_BasicAttributes(t *testing.T) {
	ctx := context.Background()
	mock := NewMockBridgeClient()
	mock.AddLight(1, "Simple Light", true, 127)

	dev := NewHueDevice("simple-light", 1, mock)
	state, err := dev.State(ctx)

	if err != nil {
		t.Fatalf("State failed: %v", err)
	}

	if state.DeviceType != "light" {
		t.Errorf("Expected DeviceType=light, got %s", state.DeviceType)
	}

	if state.Attributes["power"] != "on" {
		t.Errorf("Expected power=on, got %v", state.Attributes["power"])
	}

	brightness := state.Attributes["brightness"].(int)
	if brightness < 49 || brightness > 51 {
		t.Errorf("Expected brightness around 50, got %d", brightness)
	}

	if state.Attributes["model"] != "LCT001" {
		t.Errorf("Expected model=LCT001, got %v", state.Attributes["model"])
	}

	if state.Attributes["name"] != "Simple Light" {
		t.Errorf("Expected name=Simple Light, got %v", state.Attributes["name"])
	}
}
