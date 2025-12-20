package main

import (
	"log"

	"github.com/mismailzz/foreverstore/p2p"
)

func main() {

	tr := p2p.NewTCPTransport(
		p2p.TCPTransportOpts{
			ListenAddress: ":4000",
			Shakehand:     p2p.NoHandShakeFunc,
		},
	)
	// tr := p2p.NewTCPTransport(":4000")
	err := tr.ListenAndAccept()
	if err != nil {
		log.Fatal(err)
	}

	select {} // Block forever

}
