package device

import "time"

type State struct {
	DeviceType string
	UpdatedAt  time.Time
	Attributes map[string]interface{}
}
