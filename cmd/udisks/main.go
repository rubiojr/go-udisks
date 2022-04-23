package main

import (
	"encoding/json"
	"fmt"

	"os"

	"github.com/rubiojr/go-udisks"
)

func main() {
	client, err := udisks.NewClient()
	if err != nil {
		panic(err)
	}

	cmd := "blkdevs"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	} else {
		usage()
	}

	switch cmd {
	case "blkdevs":
		devs, err := client.BlockDevices()
		if err != nil {
			panic(err)
		}
		pretty(devs)
	case "drives":
		drives, err := client.Drives()
		if err != nil {
			panic(err)
		}
		pretty(drives)
	default:
		usage()
	}
}

func pretty(dev interface{}) {
	prettyString, _ := json.MarshalIndent(dev, "", "  ")
	fmt.Println(string(prettyString))
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: udisks <blkdevs|drives>\n")
	os.Exit(2)
}
