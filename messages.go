package main

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Message interface {
	WriteTo(c *Client) error
	ReadFrom(c *Client) error
	fmt.Stringer
}

type ProtocolVersionMessage struct {
	Major, Minor int
}

func (pvm ProtocolVersionMessage) WriteTo(c *Client) error {
	_, err := fmt.Fprintf(c, "RFB %03d.%03d\n", pvm.Major, pvm.Minor)
	return err
}

func (pvm *ProtocolVersionMessage) ReadFrom(c *Client) error {
	_, err := fmt.Fscanf(c, "RFB %03d.%03d\n", &pvm.Major, &pvm.Minor)
	return err
}

func (pvm ProtocolVersionMessage) String() string {
	return fmt.Sprintf("%#v", pvm)
}

type SecurityType byte

const (
	SecurityTypeInvalid SecurityType = iota
	SecurityTypeNone
	SecurityTypeVNCAuthentication
)

type SecurityTypeList []SecurityType

func (stl SecurityTypeList) Contains(dst SecurityType) bool {
	for _, st := range stl {
		if st == dst {
			return true
		}
	}
	return false
}

type SupportedSecurityTypesMessage struct {
	SecurityTypeList
}

func (sstm SupportedSecurityTypesMessage) WriteTo(c *Client) error {
	panic("Not implemented")
}

func (sstm *SupportedSecurityTypesMessage) ReadFrom(c *Client) error {
	var numSecurityTypes byte
	if err := binary.Read(c, binary.BigEndian, &numSecurityTypes); err != nil {
		return err
	}

	if numSecurityTypes == 0 {
		return fmt.Errorf("List of supported security types empty")
	}

	sstm.SecurityTypeList = make(SecurityTypeList, numSecurityTypes)
	return binary.Read(c, binary.BigEndian, sstm.SecurityTypeList)
}

func (sstm SupportedSecurityTypesMessage) String() string {
	return fmt.Sprintf("%#v", sstm)
}

type ErrorMessage struct {
	Message string
}

func (em ErrorMessage) WriteTo(c *Client) error {
	if err := binary.Write(c, binary.BigEndian, int32(len(em.Message))); err != nil {
		return err
	}

	_, err := io.WriteString(c, em.Message)
	return err
}

func (em *ErrorMessage) ReadFrom(c *Client) error {
	var messageLength byte
	if err := binary.Read(c, binary.BigEndian, &messageLength); err != nil {
		return err
	}

	raw := make([]byte, messageLength)
	_, err := io.ReadFull(c, raw)
	em.Message = string(raw)
	return err
}

func (em ErrorMessage) String() string {
	return fmt.Sprintf("Error: %s", em.Message)
}

type ChooseSecurityTypeMessage struct {
	SecurityType
}

func (cstm ChooseSecurityTypeMessage) WriteTo(c *Client) error {
	return binary.Write(c, binary.BigEndian, cstm.SecurityType)
}

func (cstm *ChooseSecurityTypeMessage) ReadFrom(c *Client) error {
	panic("Not implemented")
}

func (cstm ChooseSecurityTypeMessage) String() string {
	return fmt.Sprintf("%#v", cstm)
}

type SecurityResult uint32

const (
	SecurityResultOK SecurityResult = iota
	SecurityResultFailed
)

type SecurityResultMessage struct {
	SecurityResult
}

func (srm SecurityResultMessage) WriteTo(c *Client) error {
	panic("Not implemented")
}

func (srm *SecurityResultMessage) ReadFrom(c *Client) error {
	if err := binary.Read(c, binary.BigEndian, &srm.SecurityResult); err != nil {
		return err
	}

	return nil
}

func (srm SecurityResultMessage) String() string {
	return fmt.Sprintf("%#v", srm)
}

type ClientInitMessage struct {
	Share bool
}

func (cim ClientInitMessage) WriteTo(c *Client) error {
	i := uint8(0)
	if cim.Share {
		i = 1
	}
	return binary.Write(c, binary.BigEndian, i)
}

func (cim *ClientInitMessage) ReadFrom(c *Client) error {
	panic("Not implemented")
}

func (cim ClientInitMessage) String() string {
	return fmt.Sprintf("%#v", cim)
}

type ServerInitMessage struct {
	FramebufferWidth, FramebufferHeight int
	PixelFormat
	Name string
}

func (sim ServerInitMessage) WriteTo(c *Client) error {
	panic("Not implemented")
}

func (sim *ServerInitMessage) ReadFrom(c *Client) error {
	var width, height uint16
	err := binary.Read(c, binary.BigEndian, &width)
	if err != nil {
		return err
	}

	err = binary.Read(c, binary.BigEndian, &height)
	if err != nil {
		return err
	}

	sim.FramebufferWidth, sim.FramebufferHeight = int(width), int(height)

	if err := (&sim.PixelFormat).ReadFrom(c); err != nil {
		return err
	}

	var nameLength uint32
	if err := binary.Read(c, binary.BigEndian, &nameLength); err != nil {
		return err
	}

	rawName := make([]byte, nameLength)
	if err := binary.Read(c, binary.BigEndian, rawName); err != nil {
		return err
	}
	sim.Name = string(rawName)
	return nil
}

type PixelFormat struct {
	BitsPerPixel                    int
	Depth                           int
	BigEndian                       bool
	TrueColor                       bool
	RedMax, GreenMax, BlueMax       int
	RedShift, GreenShift, BlueShift int
}

type rawPixelFormat struct {
	BitsPerPixel                    uint8
	Depth                           uint8
	BigEndian                       uint8
	TrueColor                       uint8
	RedMax, GreenMax, BlueMax       uint16
	RedShift, GreenShift, BlueShift uint8
	Padding                         [3]uint8
}

func (pf PixelFormat) WriteTo(c *Client) error {
	panic("Not implemented")
}

func (pf *PixelFormat) ReadFrom(c *Client) error {
	var raw rawPixelFormat

	if err := binary.Read(c, binary.BigEndian, &raw); err != nil {
		return err
	}

	pf.BitsPerPixel = int(raw.BitsPerPixel)
	pf.Depth = int(raw.Depth)
	pf.BigEndian = raw.BigEndian != 0
	pf.TrueColor = raw.TrueColor != 0
	pf.RedMax, pf.GreenMax, pf.BlueMax = int(raw.RedMax), int(raw.GreenMax), int(raw.BlueMax)
	pf.RedShift, pf.GreenShift, pf.BlueShift = int(raw.RedShift), int(raw.GreenShift), int(raw.BlueShift)
	return nil
}

type FramebufferUpdateRequestMessage struct {
	Incremental   bool
	X, Y          int
	Width, Height int
}

type rawFramebufferUpdateRequestMessage struct {
	MessageType         uint8
	Incremental         uint8
	X, Y, Width, Height uint16
}

func (rfm FramebufferUpdateRequestMessage) WriteTo(c *Client) error {
	raw := rawFramebufferUpdateRequestMessage{
		MessageType: 3,
		Incremental: 0,
		X:           uint16(rfm.X),
		Y:           uint16(rfm.Y),
		Width:       uint16(rfm.Width),
		Height:      uint16(rfm.Height),
	}
	if rfm.Incremental {
		raw.Incremental = 1
	}

	err := binary.Write(c, binary.BigEndian, raw)
	if err != nil {
		return err
	}
	return nil
}

func (rfm *FramebufferUpdateRequestMessage) ReadFrom(c *Client) error {
	panic("Not implemented")
}

func (rfm FramebufferUpdateRequestMessage) String() string {
	return fmt.Sprintf("%#v", rfm)
}

type FramebufferUpdateMessage struct {
	Rectangles []Rectangle
}

type rawFramebufferUpdateMessage struct {
	MessageType   uint8
	Padding       byte
	NumRectangles uint16
}

func (fum FramebufferUpdateMessage) WriteTo(c *Client) error {
	panic("Not implemented")
}

func (fum *FramebufferUpdateMessage) ReadFrom(c *Client) error {
	var raw rawFramebufferUpdateMessage
	err := binary.Read(c, binary.BigEndian, &raw)
	if err != nil {
		return err
	}

	fum.Rectangles = make([]Rectangle, raw.NumRectangles)
	for i := range fum.Rectangles {
		var raw rawRectangleHeader
		if err := binary.Read(c, binary.BigEndian, &raw); err != nil {
			return err
		}
		r := &fum.Rectangles[i]
		r.X, r.Y = int(raw.X), int(raw.Y)
		r.Width, r.Height = int(raw.Width), int(raw.Height)

		enc, ok := defaultEncodings[EncodingType(raw.EncodingType)]
		if !ok {
			return fmt.Errorf("Unknown encoding")
		}
		err = enc(c, r)
		if err != nil {
			return err
		}
	}
	return nil
}

func (fum FramebufferUpdateMessage) String() string {
	return fmt.Sprintf("%#v", fum)
}

type Rectangle struct {
	X, Y          int
	Width, Height int
	RectangleData
}

type rawRectangleHeader struct {
	X, Y          uint16
	Width, Height uint16
	EncodingType  int32
}

type SetEncodingsMessage struct {
	EncodingTypes []EncodingType
}

type rawSetEncodingsMessage struct {
	MessageType  uint8
	Padding      byte
	NumEncodings uint16
}

func (sem SetEncodingsMessage) WriteTo(c *Client) error {
	raw := &rawSetEncodingsMessage{
		MessageType:  2,
		NumEncodings: uint16(len(sem.EncodingTypes)),
	}
	if err := binary.Write(c, binary.BigEndian, &raw); err != nil {
		return err
	}

	return binary.Write(c, binary.BigEndian, sem.EncodingTypes)
}

func (sem *SetEncodingsMessage) ReadFrom(c *Client) error {
	panic("Not implemented")
}

func (sem SetEncodingsMessage) String() string {
	return fmt.Sprintf("%#v", sem)
}
