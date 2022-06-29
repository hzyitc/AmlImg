package main

import (
	"os"
)

var version = "v0.0.0"

func usage() {
	print(os.Args[0] + " (" + version + ")\n")
	print("Usage:\n")
	print("  " + os.Args[0] + " unpack <img path> <extract dir path>\n")
	print("  " + os.Args[0] + " pack <img path> <dir path>\n")
	print("  " + os.Args[0] + " res_unpack <img path> <extract dir path>\n")
	print("  " + os.Args[0] + " res_pack <img path> <dir path>\n")
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
	case "res_unpack":
		os.MkdirAll(os.Args[3], 0755)

		err := res_unpack(os.Args[2], os.Args[3])
		if err != nil {
			println(err.Error())
			return
		}

	case "res_pack":
		err := res_pack(os.Args[2], os.Args[3])
		if err != nil {
			println(err.Error())
			return
		}
	}
}
