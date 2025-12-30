package p2p

import (
	"encoding/gob"
	"io"
)

// Any Protocol can define its own Decoder to decode bytes into Message.
type Decoder interface {
	Decode(io.Reader, *RPC) error
}

type GOBDecoder struct{}

// Instead of a fixed buffer, decode ONLY the GOB message
// so the underlying reader's cursor stays right at the start of the file data.
func (d *GOBDecoder) Decode(r io.Reader, rpc *RPC) error {
	return gob.NewDecoder(r).Decode(rpc)
}

// DefaultDecoder is a basic implementation of Decoder interface.
// Its not decoding the message
// just reading bytes into Message.Payload struct
type DefaultDecoder struct{}

func (d *DefaultDecoder) Decode(r io.Reader, rpc *RPC) error {

	buffer := make([]byte, 1024) // assuming max message size is 1024 bytes
	// connection read
	n, err := r.Read(buffer)
	if err != nil {
		return err
	}

	rpc.Payload = buffer[:n]

	return nil
}
