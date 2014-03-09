package main

import (
	"encoding/binary"
	"fmt"
	"image"
	"io"
	"log"
)

type Client interface {
	io.ReadWriteCloser
	SendMessage(msg Message) error
	MessageChannel() <-chan Message
	Init() error
	PixelFormat() PixelFormat
	FramebufferSize() image.Rectangle
	LastMousePosition() image.Point
	Encoding(EncodingType) Encoding
	Message(MessageType) MessageFactory

	// Can only be called before Init()
	RegisterEncoding(EncodingType, Encoding)
	RegisterMessageType(MessageType, MessageFactory)
}

type defaultClient struct {
	io.ReadWriteCloser
	// All received messages will be sent to this channel
	// provided it is non-nil.
	Messages chan Message

	MessageTypes map[MessageType]MessageFactory
	Encodings    map[EncodingType]Encoding

	pixelFormat       PixelFormat
	name              string
	framebufferWidth  int
	framebufferHeight int
	mousePosition     image.Point

	unreadByte    byte
	hasUnreadByte bool
}

func NewClient(rwc io.ReadWriteCloser) Client {
	c := &defaultClient{
		ReadWriteCloser: rwc,
		Messages:        make(chan Message),

		MessageTypes: map[MessageType]MessageFactory{},
		Encodings:    map[EncodingType]Encoding{},
	}
	for k, v := range defaultEncodings {
		c.Encodings[k] = v
	}
	for k, v := range defaultMessageTypes {
		c.MessageTypes[k] = v
	}
	return c
}

func (c *defaultClient) SendMessage(msg Message) error {
	return msg.WriteTo(c)
}

func (c *defaultClient) MessageChannel() <-chan Message {
	return c.Messages
}

func (c *defaultClient) Init() error {
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

func (c *defaultClient) PixelFormat() PixelFormat {
	return c.pixelFormat
}

func (c *defaultClient) FramebufferSize() image.Rectangle {
	return image.Rect(0, 0, c.framebufferWidth, c.framebufferHeight)
}

func (c *defaultClient) LastMousePosition() image.Point {
	return c.mousePosition
}

func (c *defaultClient) Read(d []byte) (int, error) {
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

func (c *defaultClient) RegisterEncoding(typ EncodingType, enc Encoding) {
	c.Encodings[typ] = enc
}

func (c *defaultClient) RegisterMessageType(typ MessageType, f MessageFactory) {
	c.MessageTypes[typ] = f
}

func (c *defaultClient) Encoding(et EncodingType) Encoding {
	return c.Encodings[et]
}

func (c *defaultClient) Message(mt MessageType) MessageFactory {
	return c.MessageTypes[mt]
}

func (c *defaultClient) unread(d byte) {
	if c.hasUnreadByte {
		panic("Can only unread one byte")
	}
	c.hasUnreadByte = true
	c.unreadByte = d
}

func (c *defaultClient) readError() error {
	em := &ErrorMessage{}
	em.ReadFrom(c)
	return fmt.Errorf("%s", em.Message)
}

func (c *defaultClient) worker() {
	defer c.Close()
	// TODO: Better error handling (error channel?)
	for {
		var messageType uint8
		err := binary.Read(c, binary.BigEndian, &messageType)
		if err != nil {
			log.Printf("Failed reading message: %s", err)
			return
		}
		c.unread(messageType)

		factory := c.Message(MessageType(messageType))
		if factory == nil {
			log.Printf("Unknown message type %d", messageType)
			return
		}
		msg := factory()
		err = msg.ReadFrom(c)
		if err != nil {
			log.Printf("Could not parse message of type %d: %s", messageType, err)
			return
		}
		if c.Messages != nil {
			c.Messages <- msg
		}
	}
}
