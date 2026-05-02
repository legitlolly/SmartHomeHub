package hue

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrBridgeUnreachable    = errors.New("hue bridge is unreachable")
	ErrLightNotFound        = errors.New("hue light not found")
	ErrInvalidParameter     = errors.New("invalid parameter value")
	ErrAuthenticationFailed = errors.New("authentication failed - check HUE_USERNAME")
)

// MapHueError converts huego library errors to user-friendly messages
func MapHueError(err error) error {
	if err == nil {
		return nil
	}

	errMsg := err.Error()

	if strings.Contains(errMsg, "unauthorized") || strings.Contains(errMsg, "authentication") {
		return fmt.Errorf("%w: %v", ErrAuthenticationFailed, err)
	}

	if strings.Contains(errMsg, "not found") || strings.Contains(errMsg, "resource") {
		return fmt.Errorf("%w: %v", ErrLightNotFound, err)
	}

	if strings.Contains(errMsg, "connection") || strings.Contains(errMsg, "timeout") {
		return fmt.Errorf("%w: %v", ErrBridgeUnreachable, err)
	}

	return err
}
