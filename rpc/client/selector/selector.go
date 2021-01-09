// Package selector is a way to pick a list of service nodes
package selector

import (
	"errors"

	"github.com/Allenxuxu/stark/pkg/registry"
)

// Selector builds on the registry as a mechanism to pick nodes
// and mark their status. This allows host pools and other things
// to be built using various algorithms.
type Selector interface {
	GetService(service string) ([]*registry.Service, error)
	Watch(service string) (registry.Watcher, error)
	Address(service string) string
	Options() Options
	// Close renders the selector unusable
	Close() error
	// Name of the selector
	String() string
}

var (
	ErrNotFound      = errors.New("not found")
	ErrNoneAvailable = errors.New("none available")
)
