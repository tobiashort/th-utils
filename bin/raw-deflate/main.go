package main

import (
	"compress/flate"
	"compress/zlib"
	"io"
	"os"

	"github.com/tobiashort/clap-go"
	"github.com/tobiashort/utils-go/must"
)

type Args struct {
	File string `clap:"positional,desc='The file to raw-compress. Reads from Stdin if not specified.'"`
}

func main() {
	args := Args{}
	clap.Description("Raw-compresses a given file.")
	clap.Parse(&args)

	var file *os.File

	if args.File != "" {
		file = must.Do2(os.OpenFile(args.File, os.O_RDONLY, 0))
	} else {
		file = os.Stdin
	}

	writer := must.Do2(flate.NewWriter(os.Stdout, zlib.DefaultCompression))
	defer writer.Close()
	io.Copy(writer, file)
}
