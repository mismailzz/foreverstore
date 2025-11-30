package p2p

// HandshakeFunc is a function type that defines the signature for handshake functions.
type HandshakeFunc func(Peer) error

func NOPEHandshake(Peer) error {return nil}