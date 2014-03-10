package rfb

import (
	"encoding/binary"
	"fmt"
	"image"
)

const (
	ClientMessageTypeSetPixelFormat MessageType = iota
	ClientMessageTypeSetEncodings
	ClientMessageTypeFramebufferUpdateRequest
	ClientMessageTypeKeyEvent
	ClientMessageTypePointerEvent
	ClientMessageTypeClientCutText
)

var (
	DefaultMessageTypes = map[MessageType]MessageFactory{
		ServerMessageTypeFramebufferUpdate: FramebufferUpdateMessageFactory,
		ServerMessageTypeBell:              BellMessageFactory,
		ServerMessageTypeServerCutText:     ServerCutTextMessageFactory,
	}
)

// FramebufferUpdateRequestMessage requests an update on the state
// of the framebuffer. If incremental is false, the complete
// framebuffer is sent. If incremental is true, the server assumes
// that the client is still in posession of the last framebuffer
// and only send the data needed to reconstruct the new content.
// Keep in mind that some server implementations send a black screen
// if incremental is false to assume a "well-defined state".
type FramebufferUpdateRequestMessage struct {
	Incremental bool
	Rectangle   image.Rectangle
}

type rawFramebufferUpdateRequestMessage struct {
	MessageType         MessageType
	Incremental         uint8
	X, Y, Width, Height uint16
}

func (rfm FramebufferUpdateRequestMessage) WriteTo(c Client) error {
	r := rfm.Rectangle.Canon()
	raw := rawFramebufferUpdateRequestMessage{
		MessageType: ClientMessageTypeFramebufferUpdateRequest,
		Incremental: 0,
		X:           uint16(r.Min.X),
		Y:           uint16(r.Min.Y),
		Width:       uint16(r.Dx()),
		Height:      uint16(r.Dy()),
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

func (rfm *FramebufferUpdateRequestMessage) ReadFrom(c Client) error {
	panic("Not implemented")
}

func (rfm FramebufferUpdateRequestMessage) String() string {
	return fmt.Sprintf("%#v", rfm)
}

// SetEncodingsMessage announces the supported encodings to the server.
// Order is correlated with preference by the client. The server is free
// to ignore it. RawEncoding is implicitly supported.
type SetEncodingsMessage struct {
	EncodingTypes []EncodingType
}

type rawSetEncodingsMessage struct {
	MessageType  MessageType
	Padding      byte
	NumEncodings uint16
}

func (sem SetEncodingsMessage) WriteTo(c Client) error {
	raw := &rawSetEncodingsMessage{
		MessageType:  ClientMessageTypeSetEncodings,
		NumEncodings: uint16(len(sem.EncodingTypes)),
	}
	if err := binary.Write(c, binary.BigEndian, &raw); err != nil {
		return err
	}

	return binary.Write(c, binary.BigEndian, sem.EncodingTypes)
}

func (sem *SetEncodingsMessage) ReadFrom(c Client) error {
	panic("Not implemented")
}

func (sem SetEncodingsMessage) String() string {
	return fmt.Sprintf("%#v", sem)
}

// ClientCutTextMessage sets the servers clipboard contents.
type ClientCutTextMessage struct {
	Text string
}

type rawClientCutTextMessage struct {
	MessageType MessageType
	Padding     [3]byte
	TextLength  uint32
}

func (cctm ClientCutTextMessage) WriteTo(c Client) error {
	raw := &rawClientCutTextMessage{
		MessageType: ClientMessageTypeClientCutText,
		TextLength:  uint32(len(cctm.Text)),
	}
	if err := binary.Write(c, binary.BigEndian, raw); err != nil {
		return err
	}
	_, err := c.Write([]byte(cctm.Text))
	return err
}

func (cctm ClientCutTextMessage) ReadFrom(c Client) error {
	panic("Not implemented")
}

func (cctm ClientCutTextMessage) String() string {
	return fmt.Sprintf("%#v", cctm)
}

// PointerEventMessage changes the state of the pointer device.
type PointerEventMessage struct {
	MouseState MouseState
	Position   image.Point
}

type rawPointerEventMessage struct {
	MessageType MessageType
	ButtonMask  uint8
	X, Y        uint16
}

func (pem PointerEventMessage) WriteTo(c Client) error {
	raw := &rawPointerEventMessage{
		MessageType: ClientMessageTypePointerEvent,
		ButtonMask:  pem.MouseState.Mask(),
		X:           uint16(pem.Position.X),
		Y:           uint16(pem.Position.Y),
	}
	return binary.Write(c, binary.BigEndian, raw)
}

func (pem *PointerEventMessage) ReadFrom(c Client) error {
	panic("Not implemented")
}

func (pem PointerEventMessage) String() string {
	return fmt.Sprintf("%#v", pem)
}

const (
	MouseButtonLeft = iota
	MouseButtonMiddle
	MouseButtonRight
	MouseButtonWheelUp
	MouseButtonWheelDown
	MouseButton6
	MouseButton7
)

// MouseState holds the state of the 8 mouse buttons.
type MouseState struct {
	Buttons [8]bool
}

// Mask converts the a MouseState to a bit mask as defined by the RFC.
func (ms MouseState) Mask() uint8 {
	mask := uint8(0)
	for i, b := range ms.Buttons {
		if b {
			mask |= 1 << uint(i)
		}
	}

	return mask
}

// Set sets the given mouse button to "pressed".
func (ms MouseState) Set(idx int) MouseState {
	ms.Buttons[idx] = true
	return ms
}

// Unset sets the given mouse button to "released".
func (ms MouseState) Unset(idx int) MouseState {
	ms.Buttons[idx] = false
	return ms
}

// KeyEvent changes the state of a single keyboard key.
type KeyEventMessage struct {
	Key     int
	Pressed bool
}

type rawKeyEventMessage struct {
	MessageType MessageType
	DownFlag    uint8
	Padding     [2]byte
	Key         uint32
}

func (kem KeyEventMessage) WriteTo(c Client) error {
	raw := &rawKeyEventMessage{
		MessageType: ClientMessageTypeKeyEvent,
		DownFlag:    0,
		Key:         uint32(kem.Key),
	}
	if kem.Pressed {
		raw.DownFlag = 1
	}
	return binary.Write(c, binary.BigEndian, raw)
}

func (kem *KeyEventMessage) ReadFrom(c Client) error {
	panic("Not implemented")
}

func (kem KeyEventMessage) String() string {
	return fmt.Sprintf("%#v", kem)
}
