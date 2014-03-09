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
type Encoding func(c Client, r *Rectangle) error

type RectangleData interface {
	Apply(draw.Image)
}

type EncodingType int32

const (
	EncodingTypeRaw EncodingType = iota

	EncodingTypePseudoCursor EncodingType = -239
)

func RawEncoding(c Client, r *Rectangle) error {
	bytesPerPixel := c.PixelFormat().BitsPerPixel / 8
	numPixels := r.Width * r.Height
	numBytes := numPixels * bytesPerPixel

	rd := &RawRectangleData{
		X:           r.X,
		Y:           r.Y,
		Width:       r.Width,
		Height:      r.Height,
		PixelFormat: c.PixelFormat(),
		Data:        make([]byte, numBytes),
	}
	r.RectangleData = rd

	_, err := io.ReadFull(c, rd.Data)
	return err
}

type RawRectangleData struct {
	X, Y          int
	Width, Height int
	PixelFormat   PixelFormat
	Data          []byte
}

func (rrd RawRectangleData) Apply(img draw.Image) {
	endian := binary.ByteOrder(binary.LittleEndian)
	if rrd.PixelFormat.BigEndian {
		endian = binary.BigEndian
	}

	bytesPerPixel := rrd.PixelFormat.BitsPerPixel / 8
	for y := 0; y < rrd.Height; y++ {
		for x := 0; x < rrd.Width; x++ {
			pixelIdx := y*rrd.Width + x
			pixelSlice := rrd.Data[pixelIdx*bytesPerPixel : (pixelIdx+1)*bytesPerPixel]
			var c color.Color
			switch bytesPerPixel {
			case 4:
				pixelValue := endian.Uint32(pixelSlice)
				c = color.RGBA{
					R: uint8(ShiftAndSlerp(pixelValue, rrd.PixelFormat.RedShift, rrd.PixelFormat.RedMax, 0xFF)),
					G: uint8(ShiftAndSlerp(pixelValue, rrd.PixelFormat.GreenShift, rrd.PixelFormat.GreenMax, 0xFF)),
					B: uint8(ShiftAndSlerp(pixelValue, rrd.PixelFormat.BlueShift, rrd.PixelFormat.BlueMax, 0xFF)),
					A: 255,
				}
			case 2:
				pixelValue := endian.Uint16(pixelSlice)
				c = color.RGBA{
					R: uint8(ShiftAndSlerp(uint32(pixelValue), rrd.PixelFormat.RedShift, rrd.PixelFormat.RedMax, 0xFF)),
					G: uint8(ShiftAndSlerp(uint32(pixelValue), rrd.PixelFormat.GreenShift, rrd.PixelFormat.GreenMax, 0xFF)),
					B: uint8(ShiftAndSlerp(uint32(pixelValue), rrd.PixelFormat.BlueShift, rrd.PixelFormat.BlueMax, 0xFF)),
					A: 255,
				}
			case 1:
				pixelValue := uint8(pixelSlice[0])
				c = color.RGBA{
					R: uint8(ShiftAndSlerp(uint32(pixelValue), rrd.PixelFormat.RedShift, rrd.PixelFormat.RedMax, 0xFF)),
					G: uint8(ShiftAndSlerp(uint32(pixelValue), rrd.PixelFormat.GreenShift, rrd.PixelFormat.GreenMax, 0xFF)),
					B: uint8(ShiftAndSlerp(uint32(pixelValue), rrd.PixelFormat.BlueShift, rrd.PixelFormat.BlueMax, 0xFF)),
					A: 255,
				}
			default:
				log.Printf("Unsupported BPP %d", rrd.PixelFormat.BitsPerPixel)
			}
			img.Set(x+rrd.X, y+rrd.Y, c)
		}
	}
}

func ShiftAndSlerp(val uint32, shift, inMax, outMax int) uint32 {
	in := val >> uint(shift) & uint32(inMax)
	out := float64(in) / float64(inMax) * float64(outMax)
	return uint32(out)
}

func CursorPseudoEncoding(c Client, r *Rectangle) error {
	bytesPerPixel := c.PixelFormat().BitsPerPixel / 8
	// TODO: Don't discard cursor image and mask
	buf := make([]byte, r.Width*r.Height*bytesPerPixel)
	if _, err := io.ReadFull(c, buf); err != nil {
		return err
	}
	buf = buf[0 : (r.Width+7)/8*r.Height]
	if _, err := io.ReadFull(c, buf); err != nil {
		return err
	}
	r.RectangleData = &CursorRectangleData{}
	return nil
}

type CursorRectangleData struct{}

func (crd *CursorRectangleData) Apply(img draw.Image) {}

var defaultEncodings = map[EncodingType]Encoding{
	EncodingTypeRaw:          RawEncoding,
	EncodingTypePseudoCursor: CursorPseudoEncoding,
}
