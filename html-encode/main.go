package main

import (
	"flag"
	"fmt"
	"html"
	"io"
	"os"
)

func printUsage() {
	fmt.Fprintf(os.Stderr, `Usage: html-encode [STRING]
Will read from STDIN if STRING is not defined as a parameter.
		
Flags:
`)
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	help := flag.Bool("help", false, "print help")
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
		input = string(data)
	}
	fmt.Print(html.EscapeString(input))
}
