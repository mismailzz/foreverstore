package main 

import (
	"io"
	"os"
	"fmt"
	"crypto/sha1"
	"encoding/hex"
	"strings"
)

type PathTransformFunc func(string) string

type StoreOpts struct {
	PathTransformFunc PathTransformFunc
}

var DefaultPathTransformFunc = func (key string) string {
	return key
}

func CASPathTransformFunc (key string) string {

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
	
	// Join the parts with "/" to form the final path
	return strings.Join(paths, "/")
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

	pathname := s.PathTransformFunc(key)

	// Check permissions in the provided pathname
	if err := os.MkdirAll(pathname, os.ModePerm); err != nil {
		return err
	}

	filename := "myspecialfile"
	pathAndFilename := pathname + "/" + filename

	// Create file
	file, err := os.Create(pathAndFilename)
	if err != nil {
		return err
	}
	
	n, err := io.Copy(file, r)
	if err != nil {
		return err
	}

	fmt.Printf("Writted %d bytes to %s\n", n, pathAndFilename)

	return nil
}