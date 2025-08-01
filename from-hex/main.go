package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func usage() {
	fmt.Printf(`Usage: from-hex [STRING]

If STRING is not specified it reads from STDIN.

`)
	flag.PrintDefaults()
}

func assertNil(val any) {
	if val != nil {
		panic(val)
	}
}

func main() {
	flag.Usage = usage
	flag.Parse()

	var input io.Reader

	if flag.NArg() == 0 {
		input = os.Stdin
	} else if flag.NArg() == 1 {
		input = strings.NewReader(flag.Arg(0))
	} else {
		usage()
		os.Exit(1)
	}

	decoder := hex.NewDecoder(input)
	_, err := io.Copy(os.Stdout, decoder)
	assertNil(err)
}
