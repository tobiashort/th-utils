package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func printUsage() {
	fmt.Printf(`Usage: url-encode-all [STRING]
URL encodes each character, e.g. 'p' -> '%%70'
Reads from STDIN if STRING is not defined as a parameter.
			
Flags:
`)
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	help := flag.Bool("h", false, "print help")
	flag.Parse()
	if *help {
		printUsage()
		return
	}
	if flag.NArg() > 1 {
		printUsage()
		return
	}
	input := []byte{}
	if flag.NArg() == 1 {
		input = []byte(flag.Arg(0))
	} else {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
		input = data
	}
	out := ""
	for _, b := range input {
		out += fmt.Sprintf("%%%02x", b)
	}
	fmt.Print(out)
}
