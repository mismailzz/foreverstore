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



type TCPTransportOpts struct {
	ListenAddress string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
}


type TCPTransport struct {
	TCPTransportOpts
	listener       net.Listener

	mu             sync.Mutex
	peers          map[string]Peer
	
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
	}
}

func (t *TCPTransport) ListenAndAccept() error {

	var err error
	// Create a TCP listener
	t.listener, err = net.Listen("tcp", t.ListenAddress)
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

		fmt.Printf("New incomfing connection %v\n", conn)
		go t.handleConnection(conn)
	}
}

func (t *TCPTransport) handleConnection(conn net.Conn) {

	peer := NewTCPPeer(conn, true) // outbound = true for accepted connections

	if err := t.HandshakeFunc(peer);  err != nil {
		fmt.Printf("TCP Handshake failed: %s\n", err)
		conn.Close()
		return
	}

	// Read loop 
	msg := &Message{}
	// buf := make([]byte, 1024)
	for {
		// n, err := conn.Read(buf)
		// if err != nil {
		// 	fmt.Printf("TCP Read error: %s\n", err)
		// 	conn.Close()
		// 	return
		// }
		// fmt.Printf("Received message: %+v\n", buf[:n])

		if err := t.Decoder.Decode(conn, msg); err != nil {
			fmt.Printf("TCP Decode error: %s\n", err)
			continue
		}
		msg.From = conn.RemoteAddr() // From where the message is coming
		fmt.Printf("Received message: %+v\n", msg)
	}




}