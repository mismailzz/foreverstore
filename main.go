package main

import (
	"log"
	"github.com/mismailzz/foreverstore/p2p"
	"time"
)

func main() {

	transportOpts := p2p.TCPTransportOpts {
		ListenAddress: ":3000",
		Shakehand: p2p.NoHandShakeFunc,
		Decoder: &p2p.DefaultDecoder{},
		// TODO: OnPeer func
	}
	tcpTransport := p2p.NewTCPTransport(transportOpts)

	fileServerOpts := FileServerOpts{
		StorageRootDir: "3000_network", // to differentiate and storage in the future client files in different port
		PathTransformFunc: CASPathTransformFunc,
		Transport:  tcpTransport,
	}
	
	server := NewFileServer(fileServerOpts)


	go func(){ // stop server after 3 secoond from starting server
		time.Sleep(time.Second * 3) //after 3 sec
		server.Stop()
	}()


	if err := server.Start(); err != nil {
		log.Fatalf("Error occurred during starting the server: %s\n", err)
	}

	

}
