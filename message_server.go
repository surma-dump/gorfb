package main

import (
	"encoding/binary"
	"fmt"
	"image/draw"
	"io"
)

type FramebufferUpdateMessage struct {
	Rectangles []Rectangle
}

func FramebufferUpdateMessageFactory() Message {
	return &FramebufferUpdateMessage{}
}

type rawFramebufferUpdateMessage struct {
	MessageType   MessageType
	Padding       byte
	NumRectangles uint16
}

func (fum FramebufferUpdateMessage) WriteTo(c *Client) error {
	panic("Not implemented")
}

func (fum *FramebufferUpdateMessage) ReadFrom(c *Client) error {
	var raw rawFramebufferUpdateMessage
	err := binary.Read(c, binary.BigEndian, &raw)
	if err != nil {
		return err
	}

	fum.Rectangles = make([]Rectangle, raw.NumRectangles)
	for i := range fum.Rectangles {
		var raw rawRectangleHeader
		if err := binary.Read(c, binary.BigEndian, &raw); err != nil {
			return err
		}
		r := &fum.Rectangles[i]
		r.X, r.Y = int(raw.X), int(raw.Y)
		r.Width, r.Height = int(raw.Width), int(raw.Height)

		enc, ok := c.AdditionalEncodings[EncodingType(raw.EncodingType)]
		if !ok {
			enc, ok = defaultEncodings[EncodingType(raw.EncodingType)]
			if !ok {
				return fmt.Errorf("Unknown encoding")
			}
		}
		err = enc(c, r)
		if err != nil {
			return err
		}
	}
	return nil
}

func (fum FramebufferUpdateMessage) String() string {
	return fmt.Sprintf("%#v", fum)
}

func (fum FramebufferUpdateMessage) ApplyAll(img draw.Image) {
	for _, r := range fum.Rectangles {
		r.Apply(img)
	}
}

type BellMessage struct{}

func BellMessageFactory() Message {
	return &BellMessage{}
}

type rawBellMessage struct {
	MessageType MessageType
}

func (bm BellMessage) WriteTo(c *Client) error {
	panic("Not implemented")
}

func (bm *BellMessage) ReadFrom(c *Client) error {
	var raw rawBellMessage
	if err := binary.Read(c, binary.BigEndian, &raw); err != nil {
		return err
	}
	return nil
}

func (bm BellMessage) String() string {
	return fmt.Sprintf("BELL!")
}

type ServerCutTextMessage struct {
	Text string
}

func ServerCutTextMessageFactory() Message {
	return &ServerCutTextMessage{}
}

type rawServerCutTextMessage struct {
	MessageType MessageType
	Padding     [3]byte
	TextLength  uint32
}

func (sctm ServerCutTextMessage) WriteTo(c *Client) error {
	panic("Not implemented")
}

func (sctm *ServerCutTextMessage) ReadFrom(c *Client) error {
	var raw rawServerCutTextMessage
	if err := binary.Read(c, binary.BigEndian, &raw); err != nil {
		return err
	}

	txt := make([]byte, raw.TextLength)
	if _, err := io.ReadFull(c, txt); err != nil {
		return err
	}
	sctm.Text = string(txt)
	return nil
}

func (sctm ServerCutTextMessage) String() string {
	return fmt.Sprintf("%#v", sctm)
}
