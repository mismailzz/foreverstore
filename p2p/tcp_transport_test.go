package p2p 

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {

	listenAddr := ":4000"
	transport := NewTCPTransport(listenAddr)

	assert.Equal(t, listenAddr, transport.listernAddress)

	assert.Nil(t, transport.ListenAndAccept())
	
}