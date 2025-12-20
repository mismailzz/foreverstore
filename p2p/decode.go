package p2p 

import (
	"io"	
)

// Any Protocol can define its own Decoder to decode bytes into Message.
type Decoder interface {
	Decode(io.Reader, *Message) error
}


// DefaultDecoder is a basic implementation of Decoder interface.
// Its not decoding the message
// just reading bytes into Message.Payload struct 
type DefaultDecoder struct{}

func (d *DefaultDecoder) Decode(r io.Reader, msg *Message) error {

	buffer := make([]byte, 1024) // assuming max message size is 1024 bytes
	// connection read 
	n, err := r.Read(buffer)
	if err != nil {
		return err
	}
	
	msg.Payload = buffer[:n]
	return nil
}