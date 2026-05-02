package device

import (
	"errors"
	"maps"
	"sync"
)

var ErrDeviceNotFound = errors.New("device not found")
var ErrDeviceAlreadyRegistered = errors.New("device already registered")

type Registry struct {
	mu      sync.RWMutex
	devices map[ID]Device
}

func NewRegistry() *Registry {
	return &Registry{
		devices: make(map[ID]Device),
	}
}

func (r *Registry) Register(d Device) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.devices[d.ID()]; exists {
		return ErrDeviceAlreadyRegistered
	}
	r.devices[d.ID()] = d
	return nil
}

func (r *Registry) Unregister(id ID) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.devices, id)
}

func (r *Registry) Get(id ID) (Device, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	d, ok := r.devices[id]
	if !ok {
		return nil, ErrDeviceNotFound
	}
	return d, nil
}

func (r *Registry) List() map[ID]Device {
	r.mu.RLock()
	defer r.mu.RUnlock()

	copy := make(map[ID]Device, len(r.devices))
	maps.Copy(copy, r.devices)

	return copy
}
