package main 

import (
	"log"
	"github.com/mismailzz/foreverstore/p2p"
)

func main(){
	transportOpts := p2p.TCPTransportOpts{
		ListenAddress: ":3000",
		HandshakeFunc: p2p.NOPEHandshake,
		Decoder:       p2p.DefatultDecoder{},
	}

	transport := p2p.NewTCPTransport(transportOpts)

	if err := transport.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}
 
	select {} // Block forever
}