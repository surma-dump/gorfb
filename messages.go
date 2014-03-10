package rfb

import (
	"encoding/binary"
	"fmt"
)

// Message is the generic interface for both client-to-server and
// server-to-client messages.
type Message interface {
	WriteTo(c Client) error
	ReadFrom(c Client) error
	fmt.Stringer
}

// MessageFactory is a function creating a new instance of a
// message. It is used by the parser to call ReadFrom() afterwards.
type MessageFactory func() Message

type MessageType uint8

//PixelFormat describes the format of the framebuffer.
type PixelFormat struct {
	// Bits per pixel. Has to be either 8, 16 or 32.
	BitsPerPixel int
	// Number of bits actually. Ignored.
	Depth int
	// True if the pixel bytes are given in big-endian order.
	BigEndian bool
	// True if {Red,Green,Blue}{Max,Shift} are used to extract
	// actual RGB values from the pixel data. Otherwise, a palette
	// is set and used by the server.
	TrueColor bool
	// Maximum value for each of the colors
	RedMax, GreenMax, BlueMax int
	// Number of bits to shift right to the colors LSB to the
	// bytes LSB.
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

func (pf PixelFormat) WriteTo(c Client) error {
	panic("Not implemented")
}

func (pf *PixelFormat) ReadFrom(c Client) error {
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

// Rectangle describes a rectangle as described in the RFC, containing
// encoded pixel data. Not to be confused with image.Rectangle.
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
