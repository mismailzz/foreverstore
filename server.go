package main 

import (
	"github.com/mismailzz/foreverstore/p2p"
	"log"
	"fmt"
)

type FileServerOpts struct {
	StorageRootDir string 
	PathTransformFunc PathTransformFunc
	Transport p2p.Transport
}

type FileServer struct {
	FileServerOpts
	store *Store 
	quitch chan struct{}
}

func NewFileServer(opts FileServerOpts) *FileServer {
	storeOpts := StoreOpts{
		PathTransformFunc: opts.PathTransformFunc,
		RootDir: opts.StorageRootDir,
	}
	return &FileServer{
		FileServerOpts: opts,
		store: NewStore(storeOpts),
		quitch: make(chan struct{}),
	}
}

func (s *FileServer) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}

	// as we used to select{} as blocker in the main function, now we are doing it here with proper cleaning of listener 
	s.loop()

	return nil 
}

func (s *FileServer) loop(){

	defer func(){
		log.Println("file server stopped due to user quit action")
		s.Transport.Close() // stop listening too when the user quit action
	}()

	for {
		select {
			case msg := <- s.Transport.Consume():
				fmt.Println(msg)
			case <-s.quitch:
				return
		}
	}

}

func (s *FileServer) Stop() {
	close(s.quitch)
}