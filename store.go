package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var (
	RootDirPathname = "rootDirForeverStore"
)

type PathTransformFunc func(string) PathKey

var DefaultPathTransformFunc = func(key string) PathKey {
	return PathKey{
		Pathname: key,
		Filename: key,
	}
}

func CASPathTransformFunc(key string) PathKey {

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

	return PathKey{
		Pathname: strings.Join(paths, "/"), // Join the parts with "/" to form the final path
		Filename: hashStr,
	}
}

type PathKey struct {
	Pathname string
	Filename string
}

func (p PathKey) FullPath(rootDir string) string {
	return fmt.Sprintf("%s/%s/%s", rootDir, p.Pathname, p.Filename)
}

type StoreOpts struct {
	PathTransformFunc PathTransformFunc
	RootDir           string
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) writeStream(key string, r io.Reader) error {

	pathKey := s.PathTransformFunc(key)

	pathnameWithRoot := fmt.Sprintf("%s/%s", s.RootDir, pathKey.Pathname)

	// Check permissions in the provided pathname
	if err := os.MkdirAll(pathnameWithRoot, os.ModePerm); err != nil {
		return err
	}

	filenameFullPathWithRoot := pathKey.FullPath(s.RootDir)

	// Create file
	f, err := os.Create(filenameFullPathWithRoot)
	if err != nil {
		return err
	}

	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	panic("ddd") // learning point here - even we defined panic here but its seem that when
	// tcp connection is reading it couldnt find to get the EOF from the READ IO but in the test
	// case things working fine bcause we used the bytes and intialize newReader for those bytes
	// that automatically add the EOF but in this case when the program runs and TCP connection to
	// to try to read from the stream then EOF doesnt found and its stuck or block here
	// in io.Copy call and pani didnt execute - similarly on the server.go handleMessageStoreFile()
	// doesnt have an error but its stuck on its writestream call and cant proceed further to release the
	// waitGroup - and the Copy is trying to stream from the reader of the connection to a file
	// mitigation could be we can apply the expected limit on the server side or others. will discuss next
	// let say if we apply limit then EOF happen and panic will be excuted here

	fmt.Printf("Written %d bytes to %s\n", n, filenameFullPathWithRoot)

	return nil
}

func (s *Store) readStream(key string) (io.Reader, error) {
	pathKey := s.PathTransformFunc(key)

	filenameFullPathWithRoot := pathKey.FullPath(s.RootDir)

	f, err := os.Open(filenameFullPathWithRoot)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)
	return buf, err
}

func (s *Store) Delete(key string) error {
	pathKey := s.PathTransformFunc(key)

	defer func() {
		log.Printf("deleted [%s] from disk", pathKey.Filename)
	}()

	filenameFullPathWithRoot := pathKey.FullPath(s.RootDir)

	if fileExist(filenameFullPathWithRoot) {
		// It can delete only file but not the path
		if err := os.RemoveAll(filenameFullPathWithRoot); err != nil {
			return err
		}
		// due to which we delete the parent directory - which is wierd but workaround
		// need double deletion
		parentDir := s.RootDir
		if err := os.RemoveAll(parentDir); err != nil {
			return err
		}
	}

	return nil
}

func fileExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// file doesnt exist
		return false
	}
	return true
}
