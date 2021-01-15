package rpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer_Stop(t *testing.T) {
	exit := make(chan struct{})
	close := func() error {
		select {
		case <-exit:
			return nil
		default:
			close(exit)
		}
		return nil
	}

	assert.Nil(t, close())
	assert.Nil(t, close())
	assert.Nil(t, close())
}
