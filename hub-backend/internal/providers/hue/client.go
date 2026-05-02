package hue

import (
	"context"

	"github.com/amimof/huego"
)

// BridgeClient interface allows for mocking in tests
type BridgeClient interface {
	GetLightContext(ctx context.Context, id int) (*huego.Light, error)
	SetLightStateContext(ctx context.Context, id int, state huego.State) (*huego.Response, error)
}

type HuegoBridge struct {
	bridge *huego.Bridge
}

func NewHuegoBridge(ip, username string) *HuegoBridge {
	return &HuegoBridge{
		bridge: huego.New(ip, username),
	}
}

func (h *HuegoBridge) GetLightContext(ctx context.Context, id int) (*huego.Light, error) {
	return h.bridge.GetLightContext(ctx, id)
}

func (h *HuegoBridge) SetLightStateContext(ctx context.Context, id int, state huego.State) (*huego.Response, error) {
	return h.bridge.SetLightStateContext(ctx, id, state)
}
