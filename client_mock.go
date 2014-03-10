package rfb

import (
	"image"
)

// ClientMock mocks each function of the Client interface.
// If a function is nil, a sane default value will be returned,
// otherwise the given function will be called.
type ClientMock struct {
	ReadFunc  func([]byte) (int, error)
	WriteFunc func([]byte) (int, error)
	CloseFunc func() error

	SendMessageFunc       func(Message) error
	MessageChannelFunc    func() <-chan Message
	InitFunc              func() error
	PixelFormatFunc       func() PixelFormat
	FramebufferSizeFunc   func() image.Rectangle
	LastMousePositionFunc func() image.Point
	EncodingFunc          func(EncodingType) Encoding
	MessageFunc           func(MessageType) MessageFactory

	RegisterEncodingFunc    func(EncodingType, Encoding)
	RegisterMessageTypeFunc func(MessageType, MessageFactory)
}

func (cm *ClientMock) Read(data []byte) (int, error) {
	if cm.ReadFunc == nil {
		return 0, nil
	}
	return cm.ReadFunc(data)
}

func (cm *ClientMock) Write(data []byte) (int, error) {
	if cm.WriteFunc == nil {
		return 0, nil
	}
	return cm.WriteFunc(data)
}

func (cm *ClientMock) Close() error {
	if cm.CloseFunc == nil {
		return nil
	}
	return cm.CloseFunc()
}

func (cm *ClientMock) SendMessage(msg Message) error {
	if cm.SendMessageFunc == nil {
		return nil
	}
	return cm.SendMessageFunc(msg)
}

func (cm *ClientMock) MessageChannel() <-chan Message {
	if cm.MessageChannelFunc == nil {
		return make(chan Message)
	}
	return cm.MessageChannelFunc()
}

func (cm *ClientMock) Init() error {
	if cm.InitFunc == nil {
		return nil
	}
	return cm.InitFunc()
}

func (cm *ClientMock) PixelFormat() PixelFormat {
	if cm.PixelFormat == nil {
		return PixelFormat{}
	}
	return cm.PixelFormatFunc()
}

func (cm *ClientMock) FramebufferSize() image.Rectangle {
	if cm.FramebufferSizeFunc == nil {
		return image.Rect(0, 0, 0, 0)
	}
	return cm.FramebufferSizeFunc()
}

func (cm *ClientMock) LastMousePosition() image.Point {
	if cm.LastMousePositionFunc == nil {
		return image.Point{0, 0}
	}
	return cm.LastMousePositionFunc()
}

func (cm *ClientMock) Encoding(et EncodingType) Encoding {
	if cm.EncodingFunc == nil {
		return nil
	}
	return cm.EncodingFunc(et)
}

func (cm *ClientMock) Message(mt MessageType) MessageFactory {
	if cm.MessageFunc == nil {
		return nil
	}
	return cm.MessageFunc(mt)
}

func (cm *ClientMock) RegisterEncoding(et EncodingType, enc Encoding) {
	if cm.RegisterEncodingFunc == nil {
		return
	}
	cm.RegisterEncodingFunc(et, enc)
}

func (cm *ClientMock) RegisterMessageType(mt MessageType, fac MessageFactory) {
	if cm.RegisterMessageType == nil {
		return
	}
	cm.RegisterMessageTypeFunc(mt, fac)
}
