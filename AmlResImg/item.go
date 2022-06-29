package AmlResImg

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	Item_Magic = uint32(0x27051956)
)

type Item_v2 struct {
	Magic          uint32
	Header_CRC     uint32
	Size           uint32
	DataOffset     uint32
	Entry          uint32
	NextItemOffset uint32
	Data_CRC       uint32
	Index          uint8
	Type1          uint8
	Type2          uint8
	Type3          uint8
	Name           [32]byte
}

type Item struct {
	Magic          uint32
	Header_CRC     uint32
	Size           uint32
	DataOffset     uint32
	Entry          uint32
	NextItemOffset uint32
	Data_CRC       uint32
	Index          uint8
	Type           uint32
	Name           string
}

func Item_Unpack(reader io.Reader, version uint32) (*Item, error) {
	if version == 2 {
		d := Item_v2{}

		buf := make([]byte, binary.Size(&d))
		_, err := io.ReadFull(reader, buf)
		if err != nil {
			return nil, err
		}

		err = binary.Read(bytes.NewReader(buf), binary.LittleEndian, &d)
		if err != nil {
			return nil, err
		}

		// crc := uint32(0xffffffff)
		// crc = AmlImg.AmlCRC(crc, buf[:4])
		// crc = AmlImg.AmlCRC(crc, buf[8:])
		// if d.Header_CRC != crc {
		// 	return nil, fmt.Errorf("incorrect crc: should %08X but is %08X", d.Header_CRC, crc)
		// }

		return &Item{
			d.Magic,
			d.Header_CRC,
			d.Size,
			d.DataOffset,
			d.Entry,
			d.NextItemOffset,
			d.Data_CRC,
			d.Index,
			uint32(d.Type1)<<16 | uint32(d.Type2)<<8 | uint32(d.Type3)<<0,
			string(bytes.TrimRight(d.Name[:], "\x00")),
		}, nil
	} else {
		return nil, fmt.Errorf("unsupport version: %d", version)
	}
}

func (item *Item) Pack(writer io.Writer, version uint32) error {
	if version == 2 {
		d := Item_v2{
			item.Magic,
			item.Header_CRC,
			item.Size,
			item.DataOffset,
			item.Entry,
			item.NextItemOffset,
			item.Data_CRC,
			item.Index,
			uint8((item.Type >> 16) & 0xFF),
			uint8((item.Type >> 8) & 0xFF),
			uint8((item.Type >> 0) & 0xFF),
			[32]byte{},
		}
		copy(d.Name[:], []byte(item.Name))
		return binary.Write(writer, binary.LittleEndian, d)
	} else {
		return fmt.Errorf("unsupport version: %d", version)
	}
}
