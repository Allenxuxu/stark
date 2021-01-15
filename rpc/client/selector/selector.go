package selector

import (
	"errors"

	"github.com/Allenxuxu/stark/registry"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrNoneAvailable = errors.New("none available")
)

type Selector interface {
	GetService(service string) ([]*registry.Service, error)
	Watch(service string) (registry.Watcher, error)
	Address(service string) string
	Options() Options
	Close() error
	String() string
}
