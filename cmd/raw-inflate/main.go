package main

import (
	"compress/flate"
	"io"
	"os"

	"github.com/tobiashort/clap-go"
)

type Args struct {
	File string `clap:"positional,description='The file to raw-decompress. Reads from Stdin if not specified.'"`
}

func assertNoErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	args := Args{}
	clap.Description("Raw-decompresses a given file.")
	clap.Parse(&args)

	var file *os.File
	if args.File != "" {
		var err error
		file, err = os.OpenFile(args.File, os.O_RDONLY, 0)
		assertNoErr(err)
	} else {
		file = os.Stdin
	}

	reader := flate.NewReader(file)
	defer reader.Close()
	io.Copy(os.Stdout, reader)
}
