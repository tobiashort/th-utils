package main

import (
	"compress/flate"
	"io"
	"os"

	"github.com/tobiashort/clap-go"
	. "github.com/tobiashort/utils-go/must"
)

type Args struct {
	File string `clap:"positional,description='The file to raw-decompress. Reads from Stdin if not specified.'"`
}

func main() {
	args := Args{}
	clap.Description("Raw-decompresses a given file.")
	clap.Parse(&args)

	var file *os.File
	if args.File != "" {
		file = Must2(os.OpenFile(args.File, os.O_RDONLY, 0))
	} else {
		file = os.Stdin
	}

	reader := flate.NewReader(file)
	defer reader.Close()
	io.Copy(os.Stdout, reader)
}
