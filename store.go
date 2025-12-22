package main 

import (
	"io"
	"os"
	"fmt"
	"crypto/sha1"
	"encoding/hex"
	"strings"
	"bytes"
)

type PathTransformFunc func(string) PathKey

type StoreOpts struct {
	PathTransformFunc PathTransformFunc
}


var DefaultPathTransformFunc = func (key string) string{
	return key 
}

func CASPathTransformFunc (key string) PathKey {

	// Create determistic hash from same key using SHA1 
	hash := sha1.Sum([]byte(key))
	// Convert the bytes to hex string for hash 
	hashStr := hex.EncodeToString(hash[:])

	// Split the hash string into multiple parts for directory structure (depth levels)
	blocksize := 5
	sliceLen := len(hashStr) / blocksize
	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
    	from, to := i*blocksize, (i*blocksize)+blocksize
    	paths[i] = hashStr[from:to]
	}	
	
	return PathKey {
		PathName: strings.Join(paths, "/"), // Join the parts with "/" to form the final path
		Filename: hashStr,
	}
}

type PathKey struct {
	PathName string 
	Filename string 
}

func (p PathKey) FullPath() string { 
	return fmt.Sprintf("%s/%s", p.PathName, p.Filename)
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) writeStream (key string, r io.Reader) error {

	pathKey := s.PathTransformFunc(key)

	// Check permissions in the provided pathname
	if err := os.MkdirAll(pathKey.PathName, os.ModePerm); err != nil {
		return err
	}

	fileNameFullPath := pathKey.FullPath()

	// Create file
	f, err := os.Create(fileNameFullPath)
	if err != nil {
		return err
	}
	
	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	fmt.Printf("Written %d bytes to %s\n", n, fileNameFullPath)

	return nil
}

func (s *Store) readStream (key string) (io.Reader, error){
	pathKey := s.PathTransformFunc(key)
	f, err := os.Open(pathKey.FullPath())
	if err != nil {
		return nil, err
	}
	defer f.Close()
	
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)
	return buf, err
}