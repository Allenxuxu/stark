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
	Options() Options
	Next(service string) (*registry.Node, error)
	Close() error
	String() string
}
