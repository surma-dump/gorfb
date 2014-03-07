package main

import (
	"encoding/binary"
	"fmt"
	"image"
	"io"
	"log"
)

type Client struct {
	io.ReadWriteCloser
	// All received messages will be sent to this channel
	// provided it is non-nil.
	Messages chan Message

	AdditionalMessageTypes map[MessageType]MessageFactory
	AdditionalEncodings    map[EncodingType]Encoding

	pixelFormat       PixelFormat
	name              string
	framebufferWidth  int
	framebufferHeight int
	mousePosition     image.Point

	unreadByte    byte
	hasUnreadByte bool
}

func (c *Client) PixelFormat() PixelFormat {
	return c.pixelFormat
}

func (c *Client) DesktopName() string {
	return c.name
}

func (c *Client) FramebufferSize() image.Rectangle {
	return image.Rect(0, 0, c.framebufferWidth, c.framebufferHeight)
}

func (c *Client) MousePosition() image.Point {
	return c.mousePosition
}

func (c *Client) Read(d []byte) (int, error) {
	if len(d) == 0 {
		return 0, nil
	}
	n1 := 0
	if c.hasUnreadByte {
		d[0] = c.unreadByte
		c.hasUnreadByte = false
		n1 += 1
	}
	n2, err := c.ReadWriteCloser.Read(d[n1:])
	return n1 + n2, err
}

func (c *Client) Unread(d byte) {
	if c.hasUnreadByte {
		panic("Can only unread one byte")
	}
	c.hasUnreadByte = true
	c.unreadByte = d
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

	c.pixelFormat = sim.PixelFormat
	c.name = sim.Name
	c.framebufferWidth, c.framebufferHeight = int(sim.FramebufferWidth), int(sim.FramebufferHeight)

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
		c.Unread(messageType)

		factory, ok := c.AdditionalMessageTypes[MessageType(messageType)]
		if !ok {
			factory, ok = defaultMessageTypes[MessageType(messageType)]
			if !ok {
				log.Printf("Unknown message type %d", messageType)
				return
			}
		}
		msg := factory()
		err = msg.ReadFrom(c)
		if err != nil {
			log.Printf("Could not parse message: %s", err)
			return
		}
		if c.Messages != nil {
			c.Messages <- msg
		}
	}
}

func (c *Client) RequestFramebufferUpdate(r image.Rectangle, incremental bool) {
	r = r.Canon()
	(&FramebufferUpdateRequestMessage{
		X:           r.Min.X,
		Y:           r.Min.Y,
		Width:       r.Dx(),
		Height:      r.Dy(),
		Incremental: incremental,
	}).WriteTo(c)
}

func (c *Client) SetEncodings(et ...EncodingType) {
	(&SetEncodingsMessage{
		EncodingTypes: et,
	}).WriteTo(c)
}

func (c *Client) SetClipboard(text string) {
	(&ClientCutTextMessage{
		Text: text,
	}).WriteTo(c)
}

func (c *Client) SetMouseState(pos image.Point, ms MouseState) {
	c.mousePosition = pos
	(&PointerEventMessage{
		Position:   pos,
		MouseState: ms,
	}).WriteTo(c)
}

func (c *Client) PressKey(key int) {
	(&KeyEventMessage{
		Key:     key,
		Pressed: true,
	}).WriteTo(c)
}

func (c *Client) ReleaseKey(key int) {
	(&KeyEventMessage{
		Key:     key,
		Pressed: false,
	}).WriteTo(c)
}
