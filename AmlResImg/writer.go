package AmlResImg

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/hzyitc/AmlImg/AmlCRC"
)

type ImageWriter struct {
	items []*imageWriter_Item
}

type imageWriter_Item struct {
	Type     uint32
	Name     string
	callback func(w io.Writer) error
}

func NewWriter() (*ImageWriter, error) {
	return &ImageWriter{
		make([]*imageWriter_Item, 0),
	}, nil
}

func (w *ImageWriter) Add(Type uint32, Name string, callback func(w io.Writer) error) {
	w.items = append(w.items, &imageWriter_Item{
		Type,
		Name,
		callback,
	})
}

func (w *ImageWriter) Write(path string, version uint32) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	items_current := 0 + int64(binary.Size(Header{}))
	item_size := 0
	if version == 2 {
		item_size = binary.Size(Item_v2{})
	} else {
		return fmt.Errorf("unsupport version: %d", version)
	}
	data_current := items_current + int64(item_size*len(w.items))

	for i, item := range w.items {
		_, err := file.Seek(data_current, io.SeekStart)
		if err != nil {
			return err
		}

		err = item.callback(file)
		if err != nil {
			return err
		}

		current, err := file.Seek(0, io.SeekCurrent)
		if err != nil {
			return err
		}
		size := current - data_current

		if current%16 != 0 {
			n, err := file.Write(make([]byte, 16-(current%16)))
			if err != nil {
				return err
			}
			current += int64(n)
		}

		_, err = file.Seek(items_current, io.SeekStart)
		if err != nil {
			return err
		}

		next := uint32(items_current) + uint32(item_size)
		if i == len(w.items)-1 {
			next = 0
		}

		err = (&Item{
			Magic:          Item_Magic,
			Header_CRC:     0,
			Size:           uint32(size),
			DataOffset:     uint32(data_current),
			Entry:          0,
			NextItemOffset: next,
			Data_CRC:       0,
			Index:          uint8(i),
			Type:           item.Type,
			Name:           item.Name,
		}).Pack(file, version)
		if err != nil {
			return err
		}

		items_current += int64(item_size)
		data_current = current
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	header := &Header{
		CRC:       0,
		Version:   version,
		Magic:     Header_Magic,
		Size:      uint32(data_current),
		ItemCount: uint32(len(w.items)),
		AlignSize: 16,
	}
	err = header.Pack(file)
	if err != nil {
		return err
	}

	_, err = file.Seek(4, io.SeekStart)
	if err != nil {
		return err
	}

	crc := uint32(0xffffffff)
	var buf [4096]byte
	for {
		n, err := file.Read(buf[:])
		crc = AmlCRC.AmlCRC(crc, buf[:n])
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return err
		}
	}

	header.CRC = crc
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	err = header.Pack(file)
	if err != nil {
		return err
	}

	return nil
}
