package p2p

// Some libs needed the hanshake, so lets define here
// HandShakeFunc defines the handshake function signature used by peers.
// It takes a Peer and returns an error if the handshake fails.
type HandShakeFunc func(Peer) error

// NoHandShakeFunc is a no-op handshake function that always succeeds.
func NoHandShakeFunc(p Peer) error {
	return nil
}