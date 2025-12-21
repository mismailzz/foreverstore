package main 

import (
	"testing"
	"bytes"
)

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