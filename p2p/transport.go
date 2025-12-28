package p2p

import (
	"net"
)

// Peer is an interface that represents the remote node in the p2p network.
type Peer interface {
	RemoteAddress() net.Addr
	Close() error
}

// Transport is anything that handles the communication betweend nodes/peers
// in the network. This can be of the form (TCP, UDP, websockets, gRPC, etc).
type Transport interface {
	Dial(address string) error
	ListenAndAccept() error
	Consume() <-chan RPC //only read channel
	Close() error
}
