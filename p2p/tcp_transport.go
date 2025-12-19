package p2p

import (
	"net"
	"sync"
)

type TCPTransport struct {
	listenAddress string
	listener      net.Addr

	// Better to Mutex should be above the object -- need to be protect
	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTCPTransport(listenAdrr string) *TCPTransport {
	return &TCPTransport{
		listenAddress: listenAdrr,
	}
}

