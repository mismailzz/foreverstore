package p2p 

import (
	"net"
)

// Message reprsents a generic P2P message in the TCP transport communication.
type Message struct {
	From	net.Addr // The sender's network address (i.e can be a simple string)
	Payload []byte // The actual data of the message, can be any message. 
}