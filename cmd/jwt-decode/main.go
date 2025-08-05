package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/tobiashort/th-utils/jwt-decode/jwt"
)

func printUsage() {
	fmt.Fprintf(os.Stderr, `Usage: jwt-decode [JWT]
Reads from STDIN if JWT is not provided as an argument.

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
	decoded, err := jwt.Decode(input)
	if err != nil {
		panic(err)
	}
	fmt.Println(decoded)
}
