package main 

import (
	"testing"
	"bytes"
)

func TestCASPathTransformFunc(t *testing.T) {
	key := "examplekey"
	expectedPath := "67743/9fbbb/305f1/d04e0/3730e/29d2c/78498/e5231" // Example expected path, adjust as needed
	expectedOriginal := "677439fbbb305f1d04e03730e29d2c78498e5231"
	pathKey := CASPathTransformFunc(key)
	
	if pathKey.PathName != expectedPath {
		t.Errorf("have %s want %s ", pathKey.PathName, expectedPath)
	}
	if pathKey.Original != expectedOriginal {
		t.Errorf("have %s want %s ", pathKey.Original, expectedOriginal)
	}
}

func TestStore_writeStream(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}

	store := NewStore(opts)

	data := bytes.NewReader([]byte("Hello, World!"))
	if err := store.writeStream("testkey", data); err != nil {
		t.Fatalf("writeStream failed: %v", err)
	}

}