package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

func printUsage() {
	fmt.Fprint(os.Stderr, `Usage: unixtime2rfc3339 [SECS]
Reads from STDIN if SECS is not provided as an argument.

Flags:
`)
	flag.PrintDefaults()
	os.Exit(1)
}

func invalidInput(input string) {
	fmt.Fprintf(os.Stderr, "The provided input is not an integer: %s\n", input)
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
	secs, err := strconv.Atoi(input)
	if err != nil {
		invalidInput(input)
		return
	}
	fmt.Println(time.Unix(int64(secs), 0).Format(time.RFC3339))
}
