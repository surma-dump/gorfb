package main

import (
	"encoding/binary"
	"fmt"
)

type FramebufferUpdateRequestMessage struct {
	Incremental   bool
	X, Y          int
	Width, Height int
}

type rawFramebufferUpdateRequestMessage struct {
	MessageType         uint8
	Incremental         uint8
	X, Y, Width, Height uint16
}

func (rfm FramebufferUpdateRequestMessage) WriteTo(c *Client) error {
	raw := rawFramebufferUpdateRequestMessage{
		MessageType: 3,
		Incremental: 0,
		X:           uint16(rfm.X),
		Y:           uint16(rfm.Y),
		Width:       uint16(rfm.Width),
		Height:      uint16(rfm.Height),
	}
	if rfm.Incremental {
		raw.Incremental = 1
	}

	err := binary.Write(c, binary.BigEndian, raw)
	if err != nil {
		return err
	}
	return nil
}

func (rfm *FramebufferUpdateRequestMessage) ReadFrom(c *Client) error {
	panic("Not implemented")
}

func (rfm FramebufferUpdateRequestMessage) String() string {
	return fmt.Sprintf("%#v", rfm)
}

type SetEncodingsMessage struct {
	EncodingTypes []EncodingType
}

type rawSetEncodingsMessage struct {
	MessageType  uint8
	Padding      byte
	NumEncodings uint16
}

func (sem SetEncodingsMessage) WriteTo(c *Client) error {
	raw := &rawSetEncodingsMessage{
		MessageType:  2,
		NumEncodings: uint16(len(sem.EncodingTypes)),
	}
	if err := binary.Write(c, binary.BigEndian, &raw); err != nil {
		return err
	}

	return binary.Write(c, binary.BigEndian, sem.EncodingTypes)
}

func (sem *SetEncodingsMessage) ReadFrom(c *Client) error {
	panic("Not implemented")
}

func (sem SetEncodingsMessage) String() string {
	return fmt.Sprintf("%#v", sem)
}
