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

	c := NewClient(conn)
	if err := c.Init(); err != nil {
		log.Fatalf("Could not initialize connection: %s", err)
	}
	defer c.Close()

	log.Printf("Size: %#v", c.FramebufferSize())
	log.Printf("BPP: %d, Depth: %d", c.PixelFormat().BitsPerPixel, c.PixelFormat().Depth)

	c.SendMessage(&SetEncodingsMessage{
		EncodingTypes: []EncodingType{EncodingTypePseudoCursor, EncodingTypeRaw},
	})

	TypeString(c, "Hello World[Shift+Left][Shift+Left][Ctrl+X][Ctrl+V][Ctrl+V]!")

	c.SendMessage(&FramebufferUpdateRequestMessage{
		Incremental: false,
		Rectangle:   c.FramebufferSize(),
	})
	go func() {
		for {
			time.Sleep(5 * time.Second)
			c.SendMessage(&FramebufferUpdateRequestMessage{
				Incremental: true,
				Rectangle:   c.FramebufferSize(),
			})
		}
	}()

	img := image.NewRGBA(c.FramebufferSize())
	for msg := range c.MessageChannel() {
		switch x := msg.(type) {
		case *FramebufferUpdateMessage:
			log.Printf("Updating framebuffer")
			x.ApplyAll(img)
			func() {
				f, _ := os.Create(time.Now().String() + ".png")
				defer f.Close()
				png.Encode(f, img)
			}()
		case *BellMessage:
			log.Printf("Bell!")
		case *ServerCutTextMessage:
			log.Printf("New text in clipboard: %s", x.Text)
		default:
			log.Printf("Unhandled message")
		}
	}
}
