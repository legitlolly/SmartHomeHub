package api

import (
	"encoding/json"
	"net/http"

	"github.com/legitlolly/SmartHomeHub/internal/device"
)

type Handler struct {
	registry *device.Registry
}

func NewHandler(registry *device.Registry) *Handler {
	return &Handler{
		registry: registry,
	}
}

func (h *Handler) ListDevices(w http.ResponseWriter, r *http.Request) {
	devices := h.registry.List()

	response := map[string]interface{}{
		"count":   len(devices),
		"devices": devices,
	}
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetDeviceState(w http.ResponseWriter, r *http.Request) {
	//URL stores id i.e /devices/light-1/state the id is light-1
	deviceID := r.PathValue("id")

	dev, err := h.registry.Get(device.ID(deviceID))
	if err != nil {
		http.Error(w, "Device not found", http.StatusNotFound)
		return
	}

	state, err := dev.State(r.Context())
	if err != nil {
		http.Error(w, "Failed to get device state", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id":          deviceID,
		"device_type": state.DeviceType,
		"updated_at":  state.UpdatedAt,
		"state":       state.Attributes,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
