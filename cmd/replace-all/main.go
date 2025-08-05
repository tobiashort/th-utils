package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func usage() {
	fmt.Print(`Usage: tr STR1 STR2

Reads from STDIN and transforms the string by replacing all
occurrences of STR1 with STR2.

`)
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 2 {
		usage()
		os.Exit(1)
	}

	s1 := flag.Arg(0)
	s2 := flag.Arg(1)

	bIn, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	sIn := string(bIn)
	sOut := strings.ReplaceAll(sIn, s1, s2)

	fmt.Print(sOut)
}
