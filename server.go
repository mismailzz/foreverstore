package main 

import "github.com/mismailzz/foreverstore/p2p"

type FileServerOpts struct {
	StorageRootDir string 
	PathTransformFunc PathTransformFunc
	Transport p2p.Transport
}

type FileServer struct {
	FileServerOpts
	store *Store 
}

func NewFileServer(opts FileServerOpts) *FileServer {
	storeOpts := StoreOpts{
		PathTransformFunc: opts.PathTransformFunc,
		RootDir: opts.StorageRootDir,
	}
	return &FileServer{
		FileServerOpts: opts,
		store: NewStore(storeOpts),
	}
}

func (s *FileServer) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}
	return nil 
}