package storage

type Storager interface {
	Name() string
	Up() error
	Down() error
}
