package main

import (
	"encoding/binary"
	"fmt"
	"image"
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

type PointerEventMessage struct {
	MouseState MouseState
	Position   image.Point
}

type rawPointerEventMessage struct {
	MessageType MessageType
	ButtonMask  uint8
	X, Y        uint16
}

func (pem PointerEventMessage) WriteTo(c *Client) error {
	raw := &rawPointerEventMessage{
		MessageType: ClientMessageTypePointerEvent,
		ButtonMask:  pem.MouseState.Mask(),
		X:           uint16(pem.Position.X),
		Y:           uint16(pem.Position.Y),
	}
	return binary.Write(c, binary.BigEndian, raw)
}

func (pem *PointerEventMessage) ReadFrom(c *Client) error {
	panic("Not implemented")
}

func (pem PointerEventMessage) String() string {
	return fmt.Sprintf("%#v", pem)
}

const (
	MouseButtonLeft = iota
	MouseButtonMiddle
	MouseButtonRight
	MouseButtonWheelUp
	MouseButtonWheelDown
	MouseButton6
	MouseButton7
)

type MouseState struct {
	Buttons [8]bool
}

func (ms MouseState) Mask() uint8 {
	mask := uint8(0)
	for i, b := range ms.Buttons {
		if b {
			mask |= 1 << uint(i)
		}
	}

	return mask
}

func (ms MouseState) Set(idx int) MouseState {
	ms.Buttons[idx] = true
	return ms
}

func (ms MouseState) Unset(idx int) MouseState {
	ms.Buttons[idx] = false
	return ms
}
