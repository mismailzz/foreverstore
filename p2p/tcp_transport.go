package p2p 

import (
	"net"
	"sync"
	"fmt"
)

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
	fmt.Printf("New incomfing connection %v\n", conn)
}