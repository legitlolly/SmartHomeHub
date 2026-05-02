package device

import "context"

type ID string

type Device interface {
	ID() ID
	Execute(ctx context.Context, cmd Command) error
	State(ctx context.Context) (State, error)
}
