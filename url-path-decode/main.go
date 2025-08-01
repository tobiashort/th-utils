package main

import (
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
)

func printUsage() {
	fmt.Printf(`Usage: url-path-decode [STRING]
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
	out, err := url.PathUnescape(input)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}
	fmt.Print(out)
}
