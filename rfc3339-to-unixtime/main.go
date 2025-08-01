package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

func printUsage() {
	fmt.Fprint(os.Stderr, `Usage: rfc33392unixtime [RFC3339]
Reads from STDIN if RFC3339 is not provided as an argument.

Flags:
`)
	flag.PrintDefaults()
	os.Exit(1)
}

func invalidInput(input string) {
	fmt.Fprintf(os.Stderr, `The provided input is not in RFC3339 format: %s
Example: 1970-01-01T00:00:00Z
`, input)
	os.Exit(1)
}

func main() {
	help := flag.Bool("h", false, "help")
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
	} else if flag.NArg() == 0 {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
		input = string(data)
		input = strings.TrimSpace(input)
	}
	rfc3339, err := time.Parse(time.RFC3339, input)
	if err != nil {
		invalidInput(input)
		return
	}
	fmt.Println(rfc3339.Unix())
}
