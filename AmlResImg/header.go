package AmlResImg

import (
	"encoding/binary"
	"io"
)

const (
	Header_Magic = uint64(0x215345525F4C4D41) // "AML_RES!"
)

type Header struct {
	CRC       uint32
	Version   uint32
	Magic     uint64
	Size      uint32
	ItemCount uint32
	AlignSize uint32
	Reserved  [36]byte
}

func Header_Unpack(reader io.Reader) (*Header, error) {
	header := Header{}
	err := binary.Read(reader, binary.LittleEndian, &header)
	return &header, err
}

func (header *Header) Pack(writer io.Writer) error {
	return binary.Write(writer, binary.LittleEndian, header)
}
