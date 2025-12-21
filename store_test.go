package main 

import (
	"testing"
	"bytes"
)

func TestCASPathTransformFunc(t *testing.T) {
	key := "examplekey"
	expectedPath := "67743/9fbbb/305f1/d04e0/3730e/29d2c/78498/e5231" // Example expected path, adjust as needed
	result := CASPathTransformFunc(key)
	
	if result != expectedPath {
		t.Errorf("CASPathTransformFunc(%s) = %s; want %s", key, result, expectedPath)
	}
}

func TestStore_writeStream(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: DefaultPathTransformFunc,
	}

	store := NewStore(opts)

	data := bytes.NewReader([]byte("Hello, World!"))
	if err := store.writeStream("testkey", data); err != nil {
		t.Fatalf("writeStream failed: %v", err)
	}

}