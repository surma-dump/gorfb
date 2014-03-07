package main

import (
	"encoding/binary"
	"image/color"
	"image/draw"
	"io"
	"log"
)

// An encoding reads the RectangleData for the given rectangle
// from the connection
type Encoding func(c *Client, r *Rectangle) error

type RectangleData interface {
	Apply(draw.Image)
}

func RawEncoding(c *Client, r *Rectangle) error {
	bytesPerPixel := c.PixelFormat.BitsPerPixel / 8
	numPixels := r.Width * r.Height
	numBytes := numPixels * bytesPerPixel

	rd := &RawRectangleData{
		X:           r.X,
		Y:           r.Y,
		Width:       r.Width,
		Height:      r.Height,
		PixelFormat: c.PixelFormat,
		Data:        make([]byte, numBytes),
	}
	r.RectangleData = rd

	_, err := io.ReadFull(c, rd.Data)
	return err
}

type EncodingType int32

const (
	EncodingTypeRaw EncodingType = iota
)

type RawRectangleData struct {
	X, Y          int
	Width, Height int
	PixelFormat   PixelFormat
	Data          []byte
}

func (r RawRectangleData) Apply(img draw.Image) {
	endian := binary.ByteOrder(binary.LittleEndian)
	if r.PixelFormat.BigEndian {
		endian = binary.BigEndian
	}

	bytesPerPixel := r.PixelFormat.BitsPerPixel / 8
	for y := 0; y < r.Height; y++ {
		for x := 0; x < r.Width; x++ {
			pixelIdx := y*r.Width + x
			pixelSlice := r.Data[pixelIdx*bytesPerPixel : (pixelIdx+1)*bytesPerPixel]
			var c color.Color
			switch bytesPerPixel {
			case 4:
				pixelValue := endian.Uint32(pixelSlice)
				c = color.RGBA{
					R: uint8((pixelValue >> uint32(r.PixelFormat.RedShift)) & uint32(r.PixelFormat.RedMax)),
					G: uint8((pixelValue >> uint32(r.PixelFormat.GreenShift)) & uint32(r.PixelFormat.GreenMax)),
					B: uint8((pixelValue >> uint32(r.PixelFormat.BlueShift)) & uint32(r.PixelFormat.BlueMax)),
					A: 255,
				}
			default:
				log.Printf("Unsupported BPP %d", r.PixelFormat.BitsPerPixel)
			}
			img.Set(x+r.X, y+r.Y, c)
		}
	}
}

var defaultEncodings = map[EncodingType]Encoding{
	EncodingTypeRaw: RawEncoding,
}
