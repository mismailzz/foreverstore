package p2p 

// Message reprsents a generic P2P message in the TCP transport communication.
type Message struct {
	Payload []byte // The actual data of the message, can be any message. 
}