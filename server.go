package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/mismailzz/foreverstore/p2p"
)

type FileServerOpts struct {
	StorageRootDir    string
	PathTransformFunc PathTransformFunc
	Transport         p2p.Transport
	BootstrapNodes    []string
}

type FileServer struct {
	FileServerOpts
	store  *Store
	quitch chan struct{}

	peerLock sync.Mutex
	peers    map[string]p2p.Peer
}

func NewFileServer(opts FileServerOpts) *FileServer {
	storeOpts := StoreOpts{
		PathTransformFunc: opts.PathTransformFunc,
		RootDir:           opts.StorageRootDir,
	}
	return &FileServer{
		FileServerOpts: opts,
		store:          NewStore(storeOpts),
		quitch:         make(chan struct{}),
		peers:          make(map[string]p2p.Peer),
	}
}

type Payload struct {
	Key  string
	Data []byte
}

func (s *FileServer) broadcast(p *Payload) error {
	peers := []io.Writer{}
	for _, peer := range s.peers {
		peers = append(peers, peer)
	}

	mw := io.MultiWriter(peers...) // writes to all peers

	return gob.NewEncoder(mw).Encode(p) // encode the payload and write to all peers
}

type Message struct {
	Payload any
}

type MessageStoreFile struct {
	Key string
}

func (s *FileServer) StoreData(key string, r io.Reader) error {
	// 1. Store this file to the disk - using the store package
	// 2. Broadcast this file content (or stream it) to all known peers in the network - using the transport package

	buf := new(bytes.Buffer)
	msg := &Message{
		Payload: MessageStoreFile{
			Key: key,
		},
	}

	// Encoding
	if err := gob.NewEncoder(buf).Encode(msg); err != nil {
		return err
	}

	for _, peer := range s.peers {
		if err := peer.Send(buf.Bytes()); err != nil {
			return err
		}
	}

	time.Sleep(3 * time.Second)
	// because otherwise both of the strings are being send to TCP channel at once
	// and the Read inside the loop() for the file is failing
	// if we remove this - then race condition will happen between calls - which cause deadlock

	payload := []byte("THIS IS A BIG FILE")
	for _, peer := range s.peers {
		if err := peer.Send(payload); err != nil {
			return err
		}
	}

	return nil

	// buf := new(bytes.Buffer)
	// tee := io.TeeReader(r, buf)
	// // if we read to writestream directly, then buf will be empty
	// // due to which we are using TeeReader to write to both store and buf simultaneously
	// // we can verify this by printing buf.Bytes() before and after writeStream call without TeeReader

	// if err := s.store.writeStream(key, tee); err != nil {
	// 	return err
	// }
	// payload := &Payload{
	// 	Key:  key,
	// 	Data: buf.Bytes(),
	// }

	// // fmt.Println(buf.Bytes())

	// return s.broadcast(payload)
}

func (s *FileServer) Start() error {

	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}

	s.bootstrapNetwork()

	// as we used to select{} as blocker in the main function,
	// now we are doing it here with proper cleaning of listener
	s.loop()

	return nil
}

func (s *FileServer) loop() {

	defer func() {
		log.Println("file server stopped due to user quit action")
		s.Transport.Close() // stop listening too when the user quit action
	}()

	for {
		select {
		case rpc := <-s.Transport.Consume():

			var m Message
			if err := gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(&m); err != nil {
				log.Println(err)
				return
			}

			log.Printf("recv: %+v\n", m.Payload)

			if err := s.handleMessage(rpc.From.String(), &m); err != nil {
				log.Println(err)
				return
			}

		case <-s.quitch:
			return
		}
	}

}

func (s *FileServer) handleMessage(from string, msg *Message) error {

	switch v := msg.Payload.(type) {
	case MessageStoreFile:
		return s.handleMessageStoreFile(from, v)
	default:
		return nil
	}
}

func (s *FileServer) handleMessageStoreFile(from string, msg MessageStoreFile) error {

	peer, ok := s.peers[from]
	if !ok {
		return fmt.Errorf("peer (%s) couldnt be found in the peer list", from)
	}

	if err := s.store.writeStream(msg.Key, peer); err != nil {
		return err
	}

	peer.(*p2p.TCPPeer).Wg.Done()

	return nil
}

func (s *FileServer) Stop() {
	close(s.quitch)
}

func (s *FileServer) bootstrapNetwork() error {
	for _, addr := range s.BootstrapNodes {
		go func(addr string) {
			if err := s.Transport.Dial(addr); err != nil {
				fmt.Printf("Error dialing %s: %s\n", addr, err)
			}
		}(addr)
	}
	return nil
}

func (s *FileServer) OnPeer(peer p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()

	// Add peer to the map
	s.peers[peer.RemoteAddr().String()] = peer

	log.Printf("All Peers List: %+v\n", s.peers)

	log.Printf("Peer connected: %s\n", peer.RemoteAddr().String())
	return nil
}

func init() {
	// if you are putting any time inside ANY then we have to register in gob
	gob.Register(MessageStoreFile{})
}
