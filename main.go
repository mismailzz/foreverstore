package main

import (
	"log"
	"github.com/mismailzz/foreverstore/p2p"
)

func main() {

	s1 := makeServer(":3000", "")
	s2 := makeServer(":4000", ":3000")

	go func(){
		log.Fatal(s1.Start())
	}()

	s2.Start()

}

func makeServer(listenAddr string, peerNodes ...string) *FileServer {

	// Initialize TCP transport
	transportOpts := p2p.TCPTransportOpts{
		ListenAddress: listenAddr,
		Shakehand:     p2p.NoHandShakeFunc,
		Decoder:       &p2p.DefaultDecoder{},
		// TODO: OnPeer func
	}
	tcpTransport := p2p.NewTCPTransport(transportOpts)

	// File server options configuration
	fileServerOpts := FileServerOpts{
		StorageRootDir:    listenAddr + "_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
		BootstrapNodes:    peerNodes, 
	}
	server := NewFileServer(fileServerOpts)

	// Define OnPeer function to handle new peer connections
	tcpTransport.OnPeer = server.OnPeer

	return server
}