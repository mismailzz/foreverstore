package main 

import (
	"testing"
	"bytes"
	"fmt"
	"io/ioutil"
)

func TestStore(t *testing.T){
	
	// 1. Initialize Store
	opts := StoreOpts {
		PathTransformFunc: CASPathTransformFunc,
		RootDir: RootDirPathname,
	}
	s := NewStore(opts)

	// 2. Define a file 
	key := "exampleFile"
	data := []byte("hello world!")

	// 3. Write to a file 
	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Errorf("writeStream %s\n", err)
	}

	// 4. Read to a file 
	r, err := s.readStream(key)
	if err != nil { 
		t.Errorf("readStream %s\n", err)
	}

	b, _ := ioutil.ReadAll(r)
	fmt.Println(string(b))

	if string(b) != string(data){
		t.Errorf("want %s have %s", data, b)
	}

	// 5. Delete file 

	if err = s.Delete(key); err != nil {
		t.Errorf("Error deleting file %s\n", err)
	}
}

// func TestCASPathTransformFunc(t *testing.T) {
// 	key := "examplekey"
// 	expectedPath := "67743/9fbbb/305f1/d04e0/3730e/29d2c/78498/e5231" // Example expected path, adjust as needed
// 	expectedFileName := "677439fbbb305f1d04e03730e29d2c78498e5231"
// 	pathKey := CASPathTransformFunc(key)
	
// 	if pathKey.PathName != expectedPath {
// 		t.Errorf("have %s want %s ", pathKey.PathName, expectedPath)
// 	}
// 	if pathKey.Filename != expectedFileName {
// 		t.Errorf("have %s want %s ", pathKey.Filename, expectedFileName)
// 	}
// }

// func TestStoreDefaultPathTransform(t *testing.T) {
// 	opts := StoreOpts{
// 		PathTransformFunc: DefaultPathTransformFunc,
// 	}

// 	store := NewStore(opts)

// 	data := bytes.NewReader([]byte("Hello, World!"))
// 	if err := store.writeStream("testkey", data); err != nil {
// 		t.Fatalf("writeStream failed: %v", err)
// 	}
// }