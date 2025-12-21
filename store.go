package main 

import (
	"io"
	"os"
	"fmt"
)

type PathTransformFunc func(string) string

type StoreOpts struct {
	PathTransformFunc PathTransformFunc
}

var DefaultPathTransformFunc = func (key string) string {
	return key
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