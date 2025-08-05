package main

import (
	"compress/flate"
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
	fmt.Fprint(os.Stderr, `Usage: raw-inflate [FILE]

Raw-decompresses a given FILE.
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

	reader := flate.NewReader(file)
	defer reader.Close()
	io.Copy(os.Stdout, reader)
}
