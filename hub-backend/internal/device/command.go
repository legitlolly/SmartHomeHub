package device

type Command struct {
	DeviceID ID
	Action   string
	Params   map[string]any
}
