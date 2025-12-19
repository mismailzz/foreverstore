package p2p

// Peer is an interface that represents the remote node in the p2p network.
type Peer interface {}

// Transport is anything that handles the communication betweend nodes/peers
// in the network. This can be of the form (TCP, UDP, websockets, gRPC, etc).
type Transport interface {
	ListenAndAccept() error
}
