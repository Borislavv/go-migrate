package storage

import "context"

type Factorier interface {
	Make(ctx context.Context) ([]Storager, error)
}

type Storager interface {
	Name() string
	Up() error
	Down() error
}
