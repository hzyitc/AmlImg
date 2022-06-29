package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hzyitc/AmlImg"
)

var version = "v0.0.0"

func usage() {
	print(os.Args[0] + " (" + version + ")\n")
	print("Usage:\n")
	print("  " + os.Args[0] + " unpack <img path> <extract dir path>\n")
	print("  " + os.Args[0] + " pack <img path> <dir path>\n")
}

func main() {
	if len(os.Args) != 4 {
		usage()
		return
	}

	switch os.Args[1] {
	case "unpack":
		os.MkdirAll(os.Args[3], 0755)

		err := unpack(os.Args[2], os.Args[3])
		if err != nil {
			println(err.Error())
			return
		}

	case "pack":
		err := pack(os.Args[2], os.Args[3])
		if err != nil {
			println(err.Error())
			return
		}

	}
}

func unpack(filePath, extractPath string) error {
	img, err := AmlImg.NewReader(filePath, true)
	if err != nil {
		return errors.New("NewReader error: " + err.Error())
	}
	defer img.Close()

	cmdfile, err := os.Create(extractPath + "/commands.txt")
	if err != nil {
		return errors.New("Create error: " + err.Error())
	}
	defer cmdfile.Close()

	for i := 0; i < int(img.Header.ItemCount); i++ {
		item := img.Items[i]

		filename := fmt.Sprintf("%d.%s.%s", item.Id, item.Name, item.Type)
		if item.ImgType == AmlImg.ImgType_Sparse {
			filename += ".sparse"
		}

		println("Extracting ", extractPath+"/"+filename)

		imtType := "unknown"
		if item.ImgType == AmlImg.ImgType_Normal {
			imtType = "normal"
		} else if item.ImgType == AmlImg.ImgType_Sparse {
			imtType = "sparse"
		}
		fmt.Fprintf(cmdfile, "%s:%s:%s:%s\n", item.Type, item.Name, imtType, filename)

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

func pack(filePath, dirPath string) error {
	img, err := AmlImg.NewWriter()
	if err != nil {
		return errors.New("NewWriter error: " + err.Error())
	}

	cmdfile, err := os.Open(dirPath + "/commands.txt")
	if err != nil {
		return errors.New("Open error: " + err.Error())
	}
	defer cmdfile.Close()

	scanner := bufio.NewScanner(cmdfile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		txt := scanner.Text()
		if txt == "" {
			continue
		} else if strings.HasPrefix(txt, "#") {
			continue
		}

		c := strings.SplitN(txt, ":", 4)
		Type := c[0]
		Name := c[1]
		filename := c[3]

		imgType := AmlImg.ImgType_Normal
		switch c[2] {
		case "normal":
			imgType = AmlImg.ImgType_Normal
		case "sparse":
			imgType = AmlImg.ImgType_Sparse
		default:
			return errors.New("unknown imgType: " + c[2])
		}

		img.Add(Type, Name, imgType, func(w io.Writer) error {
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
