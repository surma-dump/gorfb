package main

import (
	"encoding/binary"
	"image/color"
	"image/draw"
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
	for y := 0; y < r.Height; y++ {
		log.Printf("Reading row %d...", y)
		binary.Read(c, binary.BigEndian, rd.Data[y*r.Width*bytesPerPixel:(y+1)*r.Width*bytesPerPixel])
	}
	return binary.Read(c, binary.BigEndian, rd.Data)
}

type RawRectangleData struct {
	X, Y          int
	Width, Height int
	PixelFormat   PixelFormat
	Data          []byte
}

func (r RawRectangleData) Apply(img draw.Image) {
	bytesPerPixel := r.PixelFormat.BitsPerPixel / 8
	for y := 0; y < r.Height; y++ {
		for x := 0; x < r.Width; x++ {
			pixelIdx := y*r.Width + x
			pixelSlice := r.Data[pixelIdx*bytesPerPixel : (pixelIdx+1)*bytesPerPixel]
			var c color.Color
			switch bytesPerPixel {
			case 4:
				pixelValue := binary.BigEndian.Uint32(pixelSlice)
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

var defaultEncodings = map[int32]Encoding{
	0: RawEncoding,
}
