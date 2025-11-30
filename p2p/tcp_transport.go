package p2p 

import (
	"net"
	"sync"
	"fmt"
)

// TCPPeer represents the remote node over the TCP established connection.
type TCPPeer struct {
	// conn is the underlying connection of the peer.
	conn net.Conn
	// if we dial and retrieve a conn => outbound = true
	// if we accept and retrieve a conn => outbound = false
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}


type TCPTransport struct {
	listernAddress string
	listener       net.Listener

	mu             sync.Mutex
	peers          map[string]Peer
	
}

func NewTCPTransport(listenAddr string) *TCPTransport {
	return &TCPTransport{
		listernAddress: listenAddr,
	}
}

func (t *TCPTransport) ListenAndAccept() error {

	var err error
	// Create a TCP listener
	t.listener, err = net.Listen("tcp", t.listernAddress)
	if err != nil {
		return err
	}
	// Start accepting connections
	go t.startAcceptLoop()
	return nil

}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP Accept Error: %s\n", err)
		}
		go t.handleConnection(conn)
	}
}

func (t *TCPTransport) handleConnection(conn net.Conn) {

	peer := NewTCPPeer(conn, true) // outbound = true for accepted connections
	fmt.Printf("New incomfing connection %v\n", peer)
}