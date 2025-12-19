package p2p

import (
	"fmt"
	"net"
	"sync"
)

// TCPPeer represents a remote node in an established TCP connection.
type TCPPeer struct {
	// conn is the underlying TCP connection to the peer (remote node).
	conn net.Conn

	// if outbound is true, then the connection was initiated by this node.
	// It means we dialed out to the remote node and retrieved this connection.
	// if outbound is false, then the connection was accepted from a remote node.
	// It means the remote node dialed us and we accepted the connection.
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

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
	peer := NewTCPPeer(conn, true) // outbound = true
	fmt.Printf("New incoming connection from:%+v\n", peer)
}
