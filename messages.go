package main

import (
	"encoding/binary"
	"fmt"
)

type Message interface {
	WriteTo(c *Client) error
	ReadFrom(c *Client) error
	fmt.Stringer
}

type MessageFactory func() Message

type MessageType uint8

const (
	ServerMessageTypeFramebufferUpdate MessageType = iota
	ServerMessageTypeSetColorMapEntries
	ServerMessageTypeBell
	ServerMessageTypeServerCutText
)

var (
	defaultMessageTypes = map[MessageType]MessageFactory{
		ServerMessageTypeFramebufferUpdate: FramebufferUpdateMessageFactory,
		ServerMessageTypeBell:              BellMessageFactory,
		ServerMessageTypeServerCutText:     ServerCutTextMessageFactory,
	}
)

type PixelFormat struct {
	BitsPerPixel                    int
	Depth                           int
	BigEndian                       bool
	TrueColor                       bool
	RedMax, GreenMax, BlueMax       int
	RedShift, GreenShift, BlueShift int
}

type rawPixelFormat struct {
	BitsPerPixel                    uint8
	Depth                           uint8
	BigEndian                       uint8
	TrueColor                       uint8
	RedMax, GreenMax, BlueMax       uint16
	RedShift, GreenShift, BlueShift uint8
	Padding                         [3]uint8
}

func (pf PixelFormat) WriteTo(c *Client) error {
	panic("Not implemented")
}

func (pf *PixelFormat) ReadFrom(c *Client) error {
	var raw rawPixelFormat

	if err := binary.Read(c, binary.BigEndian, &raw); err != nil {
		return err
	}

	pf.BitsPerPixel = int(raw.BitsPerPixel)
	pf.Depth = int(raw.Depth)
	pf.BigEndian = raw.BigEndian != 0
	pf.TrueColor = raw.TrueColor != 0
	pf.RedMax, pf.GreenMax, pf.BlueMax = int(raw.RedMax), int(raw.GreenMax), int(raw.BlueMax)
	pf.RedShift, pf.GreenShift, pf.BlueShift = int(raw.RedShift), int(raw.GreenShift), int(raw.BlueShift)
	return nil
}

type Rectangle struct {
	X, Y          int
	Width, Height int
	RectangleData
}

type rawRectangleHeader struct {
	X, Y          uint16
	Width, Height uint16
	EncodingType  int32
}
