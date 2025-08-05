package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/tobiashort/th-utils/pkg/jwt"
)

func printUsage() {
	fmt.Fprintf(os.Stderr, `Usage: jwt-encode [HEADER\n\nPAYLOAD\n\nSIGNATURE]
Reads from STDIN if HEADER\n\nPAYLOAD\n\nSIGNATURE is not defined as an argument.

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
	input := ""
	if flag.NArg() == 1 {
		input = flag.Arg(0)
	} else {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
		input = strings.TrimSpace(string(data))
	}
	encoded, err := jwt.Encode(input)
	if err != nil {
		panic(err)
	}
	fmt.Println(encoded)
}
