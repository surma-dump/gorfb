package main

import (
	"encoding/binary"
	"fmt"
)

const (
	ClientMessageTypeSetPixelFormat MessageType = iota
	_
	ClientMessageTypeSetEncodings
	ClientMessageTypeFramebufferUpdateRequest
	ClientMessageTypeKeyEvent
	ClientMessageTypePointerEvent
	ClientMessageTypeClientCutText
)

type FramebufferUpdateRequestMessage struct {
	Incremental   bool
	X, Y          int
	Width, Height int
}

type rawFramebufferUpdateRequestMessage struct {
	MessageType         MessageType
	Incremental         uint8
	X, Y, Width, Height uint16
}

func (rfm FramebufferUpdateRequestMessage) WriteTo(c *Client) error {
	raw := rawFramebufferUpdateRequestMessage{
		MessageType: ClientMessageTypeFramebufferUpdateRequest,
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
	MessageType  MessageType
	Padding      byte
	NumEncodings uint16
}

func (sem SetEncodingsMessage) WriteTo(c *Client) error {
	raw := &rawSetEncodingsMessage{
		MessageType:  ClientMessageTypeSetEncodings,
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

type ClientCutTextMessage struct {
	Text string
}

type rawClientCutTextMessage struct {
	MessageType MessageType
	Padding     [3]byte
	TextLength  uint32
}

func (cctm ClientCutTextMessage) WriteTo(c *Client) error {
	raw := &rawClientCutTextMessage{
		MessageType: ClientMessageTypeClientCutText,
		TextLength:  uint32(len(cctm.Text)),
	}
	if err := binary.Write(c, binary.BigEndian, raw); err != nil {
		return err
	}
	_, err := c.Write([]byte(cctm.Text))
	return err
}

func (cctm *ClientCutTextMessage) ReadFrom(c *Client) error {
	panic("Not implemented")
}

func (cctm ClientCutTextMessage) String() string {
	return fmt.Sprintf("%#v", cctm)
}
