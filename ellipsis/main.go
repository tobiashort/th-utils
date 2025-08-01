package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func must2[T any](val T, err error) T {
	must(err)
	return val
}

func printUsage() {
	fmt.Fprint(os.Stderr, `Usage: ellipsis [-l length]

This program reads string from stdin and cuts ist at the specified length minus
three and adds three dots.

`)
	flag.PrintDefaults()
	os.Exit(1)
}

var length int

func main() {
	flag.IntVar(&length, "l", 20, "Max length of the string including the three dots")
	flag.Parse()

	if length < 0 {
		printUsage()
	}

	bytesRead := must2(io.ReadAll(os.Stdin))
	text := string(bytesRead)
	text = strings.TrimSpace(text)
	if len(text) <= length {
		fmt.Print(text)
	} else {
		fmt.Print(text[:length-3] + "...")
	}
}
