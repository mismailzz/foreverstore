package main 

import (
	"fmt"
	"github.com/mismailzz/foreverstore/p2p"
)

func main(){
	fmt.Println("Hello, AnthonyGC!")
	transport := p2p.NewTCPTransport(":3000")
	if err := transport.ListenAndAccept(); err != nil {
		fmt.Printf("Error starting TCP Transport: %s\n", err)
	}
	select {} // Block forever
}