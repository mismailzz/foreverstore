package p2p

import (
	"errors"
	"fmt"
	"net"
)

// TCPPeer represents a remote node in an established TCP connection.
type TCPPeer struct {
	// The underlying TCP connection of the peer
	net.Conn

	// if outbound is true, then the connection was initiated by this node.
	// It means we dialed out to the remote node and retrieved this connection.
	// if outbound is false, then the connection was accepted from a remote node.
	// It means the remote node dialed us and we accepted the connection.
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		Conn:     conn,
		outbound: outbound,
	}
}

// Send method implements the Peer interface to send data to the remote peer.
func (p *TCPPeer) Send(data []byte) error {
	_, err := p.Write(data)
	return err
}

// Defined to reduce the size of TCPTransport struct
type TCPTransportOpts struct {
	// public fields so callers from other packages can set options
	ListenAddress string
	Shakehand     HandShakeFunc
	Decoder       Decoder
	OnPeer        func(Peer) error
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener
	rpcchan  chan RPC
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcchan:          make(chan RPC),
	}
}

// Consume implements the Transport interface, and returns a read-only channel of RPC message coming from peers.
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcchan
}

// What Transport do? Listen/Accept
func (t *TCPTransport) ListenAndAccept() error {

	var err error
	// define listener
	t.listener, err = net.Listen("tcp", t.ListenAddress)
	if err != nil {
		return err
	}

	// Accept loop
	go t.startAcceptLoop()
	fmt.Printf("Starting a connection on port:%s\n", t.ListenAddress)

	return nil
}

// Transport Accept loop
func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()

		if errors.Is(err, net.ErrClosed) {
			return
		} // to stop accepting when the listener is closed otherwise Close for listener is happening outside causing panic

		if err != nil {
			fmt.Println("Accept error:", err)
		}
		// Handle new connection
		go t.handleNewConnection(conn, false) // false means inbound connection
	}
}

func (t *TCPTransport) handleNewConnection(conn net.Conn, outbound bool) {
	peer := NewTCPPeer(conn, outbound)
	fmt.Printf("New incoming connection from:%+v\n", peer)

	defer func() { // ensure connection is closed on exit
		conn.Close()
	}()

	// 1. Handshake
	if err := t.Shakehand(peer); err != nil {
		fmt.Printf("Error happens during handshake %s\n", err)
		return
	}

	// 2. Notify OnPeer
	if t.OnPeer != nil { // if OnPeer function is defined
		if err := t.OnPeer(peer); err != nil {
			fmt.Printf("OnPeer error: %s\n", err)
			return
		}
	}

	// 3. Read messages in loop from connection
	rpc := RPC{}
	// Read
	for {
		if err := t.Decoder.Decode(conn, &rpc); err != nil {
			fmt.Printf("Error decoding message: %s\n", err)
			return
		}

		rpc.From = conn.RemoteAddr()
		// Send to channel
		t.rpcchan <- rpc

	}
}

// implementing the Transport interface
// stop listening too when the user quit action
func (t *TCPTransport) Close() error {
	return t.listener.Close()
}

// Dial to a remote address and establish connection
func (t *TCPTransport) Dial(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	// Handle new outbound connection
	go t.handleNewConnection(conn, true) // true means outbound connection

	return nil
}
