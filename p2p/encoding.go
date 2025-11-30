package p2p 

import (
	"io"
	"encoding/gob"
)

type Decoder interface {
	Decode(io.Reader, *Message) error
}

type GOBDecoder struct{}
func (d GOBDecoder) Decode(r io.Reader, msg *Message) error {
	return gob.NewDecoder(r).Decode(msg)
}

type DefatultDecoder struct{}
func (d DefatultDecoder) Decode(r io.Reader, msg *Message) error {
	buf := make([]byte, 1024)
	n, err := r.Read(buf)
	if err != nil {
		return err
	}
	msg.Payload = buf[:n]
	return nil
}