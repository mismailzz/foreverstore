package p2p

import (
	"fmt"
	"net"
	"sync"
)

type TCPTransport struct {
	listenAddress string
	listener      net.Listener

	// Better to Mutex should be above the object -- need to be protect
	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTCPTransport(listenAdrr string) *TCPTransport {
	return &TCPTransport{
		listenAddress: listenAdrr,
	}
}

// What Transport do? Listen/Accept
func (t *TCPTransport) ListenAndAccept() error {

	var err error
	// define listener
	t.listener, err = net.Listen("tcp", t.listenAddress)
	if err != nil {
		return err
	}

	// Accept loop
	go t.startAcceptLoop()

	return nil
}

// Transport Accept loop
func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
		}
		// Handle new connection
		go t.handleNewConnection(conn)
	}
}

func (t *TCPTransport) handleNewConnection(conn net.Conn) {
	fmt.Println("Handling new connection from:", conn.RemoteAddr())
}
