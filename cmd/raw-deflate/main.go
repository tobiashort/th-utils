package main

import (
	"compress/flate"
	"compress/zlib"
	"io"
	"os"

	"github.com/tobiashort/clap-go"
)

type Args struct {
	File string `clap:"positional,description='The file to raw-compress. Reads from Stdin if not specified.'"`
}

func assertNoErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	args := Args{}
	clap.Description("Raw-compresses a given file.")
	clap.Parse(&args)

	var file *os.File

	if args.File != "" {
		var err error
		file, err = os.OpenFile(args.File, os.O_RDONLY, 0)
		assertNoErr(err)
	} else {
		file = os.Stdin
	}

	writer, err := flate.NewWriter(os.Stdout, zlib.DefaultCompression)
	assertNoErr(err)
	defer writer.Close()
	io.Copy(writer, file)
}
