package main

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

type ClosableBuffer struct {
	*bytes.Buffer
}

func (cb *ClosableBuffer) Close() error {
	return nil
}

func NewClosableBuffer() *ClosableBuffer {
	return &ClosableBuffer{
		Buffer: &bytes.Buffer{},
	}
}

func TestDefaultClient_Unread(t *testing.T) {
	c := &defaultClient{
		ReadWriteCloser: NewClosableBuffer(),
	}

	c.Write([]byte("ABC"))
	buf := make([]byte, 2)
	_, err := io.ReadFull(c, buf)
	if err != nil {
		t.Fatalf("Unexpected reading error: %s", err)
	}
	expected := []byte("AB")
	if !reflect.DeepEqual(buf, expected) {
		t.Fatalf("Expected %#v, got %#v", expected, buf)
	}

	c = &defaultClient{
		ReadWriteCloser: NewClosableBuffer(),
	}
	c.Write([]byte("ABC"))
	c.unread('0')

	buf = make([]byte, 3)
	_, err = io.ReadFull(c, buf)
	if err != nil {
		t.Fatalf("Unexpected reading error: %s", err)
	}
	expected = []byte("0AB")
	if !reflect.DeepEqual(buf, expected) {
		t.Fatalf("Expected %#v, got %#v", expected, buf)
	}
}
