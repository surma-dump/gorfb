package main

import (
	"image"
	"image/png"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	conn, err := net.Dial("tcp4", os.Args[1])
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}
	defer conn.Close()

	c := &Client{
		ReadWriteCloser: conn,
		Messages:        make(chan Message),
	}

	if err := c.Init(); err != nil {
		log.Fatalf("Could not initialize connection: %s", err)
	}
	defer c.Close()

	log.Printf("Size: %#v", c.FramebufferSize())
	log.Printf("BPP: %d, Depth: %d, Name: %s", c.PixelFormat().BitsPerPixel, c.PixelFormat().Depth, c.DesktopName())

	c.SetEncodings(EncodingTypePseudoCursor, EncodingTypeRaw)

	c.RequestFramebufferUpdate(c.FramebufferSize(), false)
	go func() {
		for {
			time.Sleep(5 * time.Second)
			c.RequestFramebufferUpdate(c.FramebufferSize(), true)
		}
	}()

	img := image.NewRGBA(c.FramebufferSize())
	for msg := range c.Messages {
		log.Printf("Received event: %s", msg)

		switch x := msg.(type) {
		case *FramebufferUpdateMessage:
			x.ApplyAll(img)
			f, err := os.Create(time.Now().String() + ".png")
			if err != nil {
				log.Fatalf("Could not open file: %s", err)
			}
			defer f.Close()
			png.Encode(f, img)
		default:
			log.Printf("Unhandled message")
		}
	}
}
