package AmlImg

import (
	"encoding/binary"
	"io"
)

const (
	Magic = uint32(0x27B51956)
)

type Header struct {
	CRC       uint32
	Version   uint32
	Magic     uint32
	Size      uint64
	AlignSize uint32
	ItemCount uint32
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
