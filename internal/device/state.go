package device

import "time"

type State struct {
	Power      string
	Brightness int
	UpdatedAt  time.Time
}
