package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {

	ListenAddress := ":4000"
	tr := NewTCPTransport(ListenAddress)

	assert.Equal(t, tr.listenAddress, ListenAddress)

}
