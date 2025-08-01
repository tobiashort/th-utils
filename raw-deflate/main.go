package main

import (
	"compress/flate"
	"compress/zlib"
	"flag"
	"fmt"
	"io"
	"os"
)

func assertNoErr(err error) {
	if err != nil {
		panic(err)
	}
}

func printUsage() {
	fmt.Fprint(os.Stderr, `Usage: raw-deflate [FILE]

Raw-compresses a given FILE.
If no FILE is specified, this program reads from STDIN.
`)

	flag.PrintDefaults()
}

func main() {
	flag.Usage = printUsage
	flag.Parse()

	var file *os.File
	
	if flag.NArg() > 1 {
		printUsage()
		os.Exit(1)
	} else if flag.NArg() == 1 {
		var err error
		file, err = os.OpenFile(flag.Arg(0), os.O_RDONLY, 0)
		assertNoErr(err)
	} else {
		file = os.Stdin
	}

	writer, err := flate.NewWriter(os.Stdout, zlib.DefaultCompression)
	assertNoErr(err)
	defer writer.Close()
	io.Copy(writer, file)
}
