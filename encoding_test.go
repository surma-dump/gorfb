package main

import (
	"image"
	"testing"
)

func TestRawRectangleData_Apply_32BPP_BigEndian(t *testing.T) {
	rrd := &RawRectangleData{
		X:      0,
		Y:      0,
		Width:  3,
		Height: 1,
		PixelFormat: PixelFormat{
			BitsPerPixel: 32,
			Depth:        24,
			BigEndian:    true,
			TrueColor:    true,
			RedMax:       255,
			GreenMax:     255,
			BlueMax:      255,
			RedShift:     24,
			GreenShift:   16,
			BlueShift:    8,
		},
		Data: []byte{
			0xFF, 0x00, 0x00, 0x00,
			0x00, 0xFF, 0x00, 0x80,
			0x00, 0x00, 0xFF, 0xFF,
		},
	}

	img := image.NewRGBA(image.Rect(0, 0, 3, 1))
	rrd.Apply(img)

	if r, g, b, _ := img.At(0, 0).RGBA(); !(r == 0xFFFF && g == 0x00 && b == 0x00) {
		t.Fatalf("Expected red got %#v", img.At(0, 0))
	}
	if r, g, b, _ := img.At(1, 0).RGBA(); !(r == 0x00 && g == 0xFFFF && b == 0x00) {
		t.Fatalf("Expected green got %#v", img.At(1, 0))
	}
	if r, g, b, _ := img.At(2, 0).RGBA(); !(r == 0x00 && g == 0x00 && b == 0xFFFF) {
		t.Fatalf("Expected blue got %#v", img.At(2, 0))
	}
}

func TestRawRectangleData_Apply_32BPP_LittleEndian(t *testing.T) {
	rrd := &RawRectangleData{
		X:      0,
		Y:      0,
		Width:  3,
		Height: 1,
		PixelFormat: PixelFormat{
			BitsPerPixel: 32,
			Depth:        24,
			BigEndian:    false,
			TrueColor:    true,
			RedMax:       255,
			GreenMax:     255,
			BlueMax:      255,
			RedShift:     24,
			GreenShift:   16,
			BlueShift:    8,
		},
		Data: []byte{
			0x00, 0x00, 0x00, 0xFF,
			0x00, 0x00, 0xFF, 0x00,
			0x00, 0xFF, 0x00, 0x00,
		},
	}

	img := image.NewRGBA(image.Rect(0, 0, 3, 1))
	rrd.Apply(img)

	if r, g, b, _ := img.At(0, 0).RGBA(); !(r == 0xFFFF && g == 0x00 && b == 0x00) {
		t.Fatalf("Expected red got %#v", img.At(0, 0))
	}
	if r, g, b, _ := img.At(1, 0).RGBA(); !(r == 0x00 && g == 0xFFFF && b == 0x00) {
		t.Fatalf("Expected green got %#v", img.At(1, 0))
	}
	if r, g, b, _ := img.At(2, 0).RGBA(); !(r == 0x00 && g == 0x00 && b == 0xFFFF) {
		t.Fatalf("Expected blue got %#v", img.At(2, 0))
	}
}

func TestRawRectangleData_Apply_Shifts(t *testing.T) {
	rrd := &RawRectangleData{
		X:      0,
		Y:      0,
		Width:  3,
		Height: 1,
		PixelFormat: PixelFormat{
			BitsPerPixel: 32,
			Depth:        24,
			BigEndian:    true,
			TrueColor:    true,
			RedMax:       255,
			GreenMax:     255,
			BlueMax:      255,
			RedShift:     8,
			GreenShift:   16,
			BlueShift:    24,
		},
		Data: []byte{
			0x00, 0x00, 0xFF, 0x00,
			0x00, 0xFF, 0x00, 0x00,
			0xFF, 0x00, 0x00, 0x00,
		},
	}

	img := image.NewRGBA(image.Rect(0, 0, 3, 1))
	rrd.Apply(img)

	if r, g, b, _ := img.At(0, 0).RGBA(); !(r == 0xFFFF && g == 0x00 && b == 0x00) {
		t.Fatalf("Expected red got %#v", img.At(0, 0))
	}
	if r, g, b, _ := img.At(1, 0).RGBA(); !(r == 0x00 && g == 0xFFFF && b == 0x00) {
		t.Fatalf("Expected green got %#v", img.At(1, 0))
	}
	if r, g, b, _ := img.At(2, 0).RGBA(); !(r == 0x00 && g == 0x00 && b == 0xFFFF) {
		t.Fatalf("Expected blue got %#v", img.At(2, 0))
	}
}

func TestRawRectangleData_Apply_16BPP(t *testing.T) {
	rrd := &RawRectangleData{
		X:      0,
		Y:      0,
		Width:  3,
		Height: 1,
		PixelFormat: PixelFormat{
			BitsPerPixel: 16,
			Depth:        15,
			BigEndian:    true,
			TrueColor:    true,
			RedMax:       31,
			GreenMax:     31,
			BlueMax:      31,
			RedShift:     11,
			GreenShift:   6,
			BlueShift:    0,
		},
		Data: []byte{
			0xF8, 0x00,
			0x07, 0xE0,
			0x00, 0x1F,
		},
	}

	img := image.NewRGBA(image.Rect(0, 0, 3, 1))
	rrd.Apply(img)

	if r, g, b, _ := img.At(0, 0).RGBA(); !(r == 0xFFFF && g == 0x00 && b == 0x00) {
		t.Fatalf("Expected red got %#v", img.At(0, 0))
	}
	if r, g, b, _ := img.At(1, 0).RGBA(); !(r == 0x00 && g == 0xFFFF && b == 0x00) {
		t.Fatalf("Expected green got %#v", img.At(1, 0))
	}
	if r, g, b, _ := img.At(2, 0).RGBA(); !(r == 0x00 && g == 0x00 && b == 0xFFFF) {
		t.Fatalf("Expected blue got %#v", img.At(2, 0))
	}
}
