package AmlImg

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	ImgType_Normal = uint32(0x00)
	ImgType_Sparse = uint32(0xFE)
)

type Item_v1 struct {
	Id            uint32
	ImgType       uint32
	OffsetOfItem  uint64
	OffsetOfImage uint64
	Size          uint64
	Type          [32]byte
	Name          [32]byte
	Reserved      [32]byte
}

type Item_v2 struct {
	Id            uint32
	ImgType       uint32
	OffsetOfItem  uint64
	OffsetOfImage uint64
	Size          uint64
	Type          [256]byte
	Name          [256]byte
	Reserved      [32]byte
}

type Item struct {
	Id            uint32
	ImgType       uint32
	OffsetOfItem  uint64
	OffsetOfImage uint64
	Size          uint64
	Type          string
	Name          string
}

func Item_Unpack(reader io.Reader, version uint32) (*Item, error) {
	if version == 1 {
		d := Item_v1{}
		err := binary.Read(reader, binary.LittleEndian, &d)
		if err != nil {
			return nil, err
		}

		return &Item{
			d.Id,
			d.ImgType,
			d.OffsetOfItem,
			d.OffsetOfImage,
			d.Size,
			string(bytes.TrimRight(d.Type[:], "\x00")),
			string(bytes.TrimRight(d.Name[:], "\x00")),
		}, nil
	} else if version == 2 {
		d := Item_v2{}
		err := binary.Read(reader, binary.LittleEndian, &d)
		if err != nil {
			return nil, err
		}

		return &Item{
			d.Id,
			d.ImgType,
			d.OffsetOfItem,
			d.OffsetOfImage,
			d.Size,
			string(bytes.TrimRight(d.Type[:], "\x00")),
			string(bytes.TrimRight(d.Name[:], "\x00")),
		}, nil
	} else {
		return nil, fmt.Errorf("unsupport version: %d", version)
	}
}

func (item *Item) Pack(writer io.Writer, version uint32) error {
	if version == 1 {
		d := Item_v1{
			item.Id,
			item.ImgType,
			item.OffsetOfItem,
			item.OffsetOfImage,
			item.Size,
			[32]byte{},
			[32]byte{},
			[32]byte{},
		}
		copy(d.Type[:], []byte(item.Type))
		copy(d.Name[:], []byte(item.Name))
		if item.Type == "PARTITION" {
			d.Reserved[1] = 1 // Unknown why
		}
		return binary.Write(writer, binary.LittleEndian, d)
	} else if version == 2 {
		d := Item_v2{
			item.Id,
			item.ImgType,
			item.OffsetOfItem,
			item.OffsetOfImage,
			item.Size,
			[256]byte{},
			[256]byte{},
			[32]byte{},
		}
		copy(d.Type[:], []byte(item.Type))
		copy(d.Name[:], []byte(item.Name))
		if item.Type == "PARTITION" {
			d.Reserved[0] = 1 // Unknown why
		}
		return binary.Write(writer, binary.LittleEndian, d)
	} else {
		return fmt.Errorf("unsupport version: %d", version)
	}
}
