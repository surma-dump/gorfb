package main

import (
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
		EvCh:            make(chan bool),
	}

	if err := c.Init(); err != nil {
		log.Fatalf("Could not initialize connection: %s", err)
	}
	defer c.Close()

	r := c.Framebuffer.Bounds().Canon()
	log.Printf("%d x %d", r.Dx(), r.Dy())
	log.Printf("BPP: %d, Depth: %d, Name: %s", c.PixelFormat.BitsPerPixel, c.PixelFormat.Depth, c.Name)

	c.SetEncodings(EncodingTypeRaw)
	time.Sleep(5 * time.Second)
	c.RequestFramebufferUpdate(c.Framebuffer.Bounds().Canon(), false)
	<-c.EvCh
	log.Printf("Received event")

	f, err := os.Create(time.Now().String() + ".png")
	if err != nil {
		log.Fatalf("Could not open file: %s", err)
	}
	defer f.Close()
	png.Encode(f, c.Framebuffer)

	time.Sleep(10 * time.Second)
	c.RequestFramebufferUpdate(c.Framebuffer.Bounds().Canon(), false)
	<-c.EvCh
	log.Printf("Received event")

	f, err = os.Create(time.Now().String() + ".png")
	if err != nil {
		log.Fatalf("Could not open file: %s", err)
	}
	defer f.Close()
	png.Encode(f, c.Framebuffer)
}
