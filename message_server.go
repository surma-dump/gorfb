package main

import (
	"encoding/binary"
	"fmt"
	"image/draw"
)

type FramebufferUpdateMessage struct {
	Rectangles []Rectangle
}

func FramebufferUpdateMessageFactory() Message {
	return &FramebufferUpdateMessage{}
}

type rawFramebufferUpdateMessage struct {
	MessageType   uint8
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
