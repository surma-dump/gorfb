package main

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/draw"
	"io"
	"log"
)

type Client struct {
	io.ReadWriteCloser
	EvCh chan bool
	// To be set by server messages
	PixelFormat PixelFormat
	Name        string
	Framebuffer draw.Image
}

func (c *Client) Init() error {
	pvm := &ProtocolVersionMessage{}
	if err := pvm.ReadFrom(c); err != nil {
		return fmt.Errorf("Could not get version number of server: %s", err)
	}
	if pvm.Major != 3 || pvm.Minor != 8 {
		return fmt.Errorf("Unsupported protocol version: %d.%d", pvm.Major, pvm.Minor)
	}

	if err := pvm.WriteTo(c); err != nil {
		return fmt.Errorf("Could not sent version number to server: %s", err)
	}

	sstm := &SupportedSecurityTypesMessage{}
	if err := sstm.ReadFrom(c); err != nil {
		return fmt.Errorf("Enumeration of security methods failed: %s", err)
	}
	if len(sstm.SecurityTypeList) == 0 {
		return fmt.Errorf("Securty type list is empty: %s", c.readError())
	}
	if !sstm.SecurityTypeList.Contains(SecurityTypeNone) {
		return fmt.Errorf("Desired security type not supported")
	}

	cstm := ChooseSecurityTypeMessage{SecurityTypeNone}
	if err := cstm.WriteTo(c); err != nil {
		return fmt.Errorf("Could not send credentials: %s", err)
	}

	srm := &SecurityResultMessage{}
	if err := srm.ReadFrom(c); err != nil {
		return fmt.Errorf("Could not obtain security result: %s", err)
	}
	if srm.SecurityResult == SecurityResultFailed {
		return fmt.Errorf("Authentication failed: %s", c.readError())
	}

	cim := ClientInitMessage{false}
	if err := cim.WriteTo(c); err != nil {
		return fmt.Errorf("Could not initialize client: %s", err)
	}

	sim := &ServerInitMessage{}
	if err := sim.ReadFrom(c); err != nil {
		return fmt.Errorf("Could not read server init message: %s", err)
	}

	c.PixelFormat = sim.PixelFormat
	c.Name = sim.Name
	c.Framebuffer = image.NewRGBA(image.Rect(0, 0, sim.FramebufferWidth, sim.FramebufferHeight))

	go c.worker()

	return nil
}

func (c *Client) readError() error {
	em := &ErrorMessage{}
	em.ReadFrom(c)
	return fmt.Errorf("%s", em.Message)
}

func (c *Client) worker() {
	defer c.Close()
	// TODO: Better error handling (error channel?)
	for {
		var messageType uint8
		err := binary.Read(c, binary.BigEndian, &messageType)
		if err != nil {
			log.Printf("Failed reading message: %s", err)
			return
		}

		switch messageType {
		case 0:
			fum := &FramebufferUpdateMessage{}
			err := fum.ReadFrom(c)
			if err != nil {
				log.Printf("Could not parse message: %s", err)
				return
			}
			for _, r := range fum.Rectangles {
				r.RectangleData.Apply(c.Framebuffer)
			}
			c.EvCh <- true
		default:
			log.Printf("Unknown message type %d", messageType)
			return
		}
	}
}

func (c *Client) RequestFramebufferUpdate(r image.Rectangle) {
	r = r.Canon()
	(&FramebufferUpdateRequestMessage{
		X:           r.Min.X,
		Y:           r.Min.Y,
		Width:       r.Dx(),
		Height:      r.Dy(),
		Incremental: true,
	}).WriteTo(c)
}
