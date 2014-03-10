package rfb

import (
	"encoding/binary"
	"fmt"
	"io"
)

type ProtocolVersionMessage struct {
	Major, Minor int
}

func (pvm ProtocolVersionMessage) WriteTo(c Client) error {
	_, err := fmt.Fprintf(c, "RFB %03d.%03d\n", pvm.Major, pvm.Minor)
	return err
}

func (pvm *ProtocolVersionMessage) ReadFrom(c Client) error {
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

func (sstm SupportedSecurityTypesMessage) WriteTo(c Client) error {
	panic("Not implemented")
}

func (sstm *SupportedSecurityTypesMessage) ReadFrom(c Client) error {
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

func (em ErrorMessage) WriteTo(c Client) error {
	if err := binary.Write(c, binary.BigEndian, int32(len(em.Message))); err != nil {
		return err
	}

	_, err := io.WriteString(c, em.Message)
	return err
}

func (em *ErrorMessage) ReadFrom(c Client) error {
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

func (cstm ChooseSecurityTypeMessage) WriteTo(c Client) error {
	return binary.Write(c, binary.BigEndian, cstm.SecurityType)
}

func (cstm *ChooseSecurityTypeMessage) ReadFrom(c Client) error {
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

func (srm SecurityResultMessage) WriteTo(c Client) error {
	panic("Not implemented")
}

func (srm *SecurityResultMessage) ReadFrom(c Client) error {
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

func (cim ClientInitMessage) WriteTo(c Client) error {
	i := uint8(0)
	if cim.Share {
		i = 1
	}
	return binary.Write(c, binary.BigEndian, i)
}

func (cim *ClientInitMessage) ReadFrom(c Client) error {
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

func (sim ServerInitMessage) WriteTo(c Client) error {
	panic("Not implemented")
}

func (sim *ServerInitMessage) ReadFrom(c Client) error {
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
