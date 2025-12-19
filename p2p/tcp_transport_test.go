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

func TestTCPTransport_ListenAndAccept(t *testing.T) {

	ListenAddress := ":4001" // Failure "4000" or port already in use, etc.
	tr := NewTCPTransport(ListenAddress)

	err := tr.ListenAndAccept()
	assert.Nil(t, err) // Expect no error

}
