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
			Decoder:       &p2p.DefaultDecoder{},
			OnPeer: func(peer p2p.Peer) error { 
				log.Printf("New peer connected: %v\n", peer)
				return nil
			},
		},
	)
	// tr := p2p.NewTCPTransport(":4000")
	err := tr.ListenAndAccept()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			rpc := <-tr.Consume()
			log.Printf("Received message from %s: %s\n", rpc.From.String(), string(rpc.Payload))
		}
	}()

	select {} // Block forever

}
