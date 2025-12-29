package p2p

import (
	"net"
)

// Peer is an interface that represents the remote node in the p2p network.
type Peer interface {
	net.Conn // Embeds net.Conn interface for network connection methods (it has all the required methods)
	Send(data []byte) error
}

// Transport is anything that handles the communication betweend nodes/peers
// in the network. This can be of the form (TCP, UDP, websockets, gRPC, etc).
type Transport interface {
	Dial(address string) error
	ListenAndAccept() error
	Consume() <-chan RPC //only read channel
	Close() error
}
