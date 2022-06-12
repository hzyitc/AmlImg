package AmlImg

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
)

type ImageWriter struct {
	items []*imageWriter_Item
}

type imageWriter_Item struct {
	Type     string
	Name     string
	imgType  uint32
	callback func(w io.Writer) error
}

func NewWriter() (*ImageWriter, error) {
	return &ImageWriter{
		make([]*imageWriter_Item, 0),
	}, nil
}

func (w *ImageWriter) Add(Type string, Name string, imgType uint32, callback func(w io.Writer) error) {
	w.items = append(w.items, &imageWriter_Item{
		Type,
		Name,
		imgType,
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
	if version == 1 {
		item_size = binary.Size(Item_v1{})
	} else if version == 2 {
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

		if current%4 != 0 {
			n, err := file.Write(make([]byte, 4-(current%4)))
			if err != nil {
				return err
			}
			current += int64(n)
		}

		_, err = file.Seek(items_current, io.SeekStart)
		if err != nil {
			return err
		}

		err = (&Item{
			Id:            uint32(i),
			ImgType:       item.imgType,
			OffsetOfItem:  0,
			OffsetOfImage: uint64(data_current),
			Size:          uint64(size),
			Type:          item.Type,
			Name:          item.Name,
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
		Magic:     Magic,
		Size:      uint64(data_current),
		AlignSize: 4,
		ItemCount: uint32(len(w.items)),
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
		crc = AmlCRC(crc, buf[:n])
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
