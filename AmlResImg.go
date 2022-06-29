package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/hzyitc/AmlImg/AmlResImg"
)

func res_unpack(filePath, extractPath string) error {
	img, err := AmlResImg.NewReader(filePath, true)
	if err != nil {
		return errors.New("NewReader error: " + err.Error())
	}
	defer img.Close()

	listfile, err := os.Create(extractPath + "/list.txt")
	if err != nil {
		return errors.New("Create error: " + err.Error())
	}
	defer listfile.Close()

	for i := 0; i < int(img.Header.ItemCount); i++ {
		item := img.Items[i]

		filename := fmt.Sprintf("%d.%s", item.Index, item.Name)
		if item.Type == 0x090000 {
			filename += ".bmp"
		}

		println("Extracting ", extractPath+"/"+filename)

		fmt.Fprintf(listfile, "%08X:%s:%s\n", item.Type, item.Name, filename)

		file, err := os.Create(extractPath + "/" + filename)
		if err != nil {
			return errors.New("Create error:" + err.Error())
		}

		err = img.Seek(uint32(i), 0)
		if err != nil {
			file.Close()
			return errors.New("Seek error:" + err.Error())
		}

		_, err = io.Copy(file, img)
		if err != nil {
			file.Close()
			return errors.New("Copy error:" + err.Error())
		}

		file.Close()
	}

	return nil
}

func res_pack(filePath, dirPath string) error {
	img, err := AmlResImg.NewWriter()
	if err != nil {
		return errors.New("NewWriter error: " + err.Error())
	}

	listfile, err := os.Open(dirPath + "/list.txt")
	if err != nil {
		return errors.New("Open error: " + err.Error())
	}
	defer listfile.Close()

	scanner := bufio.NewScanner(listfile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		txt := scanner.Text()
		if txt == "" {
			continue
		} else if strings.HasPrefix(txt, "#") {
			continue
		}

		c := strings.SplitN(txt, ":", 3)
		Name := c[1]
		filename := c[2]

		Type, err := strconv.ParseInt(c[0], 16, 24)
		if err != nil {
			return errors.New("ParseInt error: " + err.Error())
		}

		img.Add(uint32(Type), Name, func(w io.Writer) error {
			println("Packing ", filename)

			file, err := os.Open(dirPath + "/" + filename)
			if err != nil {
				return errors.New("Open error: " + err.Error())
			}

			_, err = io.Copy(w, file)
			return err
		})
	}

	err = scanner.Err()
	if err != nil {
		return errors.New("scanner error: " + err.Error())
	}

	err = img.Write(filePath, 2)
	if err != nil {
		return errors.New("Write error: " + err.Error())
	}

	return nil
}
