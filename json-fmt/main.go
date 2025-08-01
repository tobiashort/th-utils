package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
)

func PrintUsage() {
	fmt.Printf(`json-fmt [JSON]

Arguments:
	json-fmt takes a JSON string or reads from STDIN

`)
	flag.PrintDefaults()
}

func main() {
	flag.Parse()
	flag.Usage = PrintUsage
	if flag.NArg() > 1 {
		PrintUsage()
		os.Exit(1)
	}
	var input []byte
	if flag.NArg() == 1 {
		input = []byte(flag.Arg(0))
	} else {
		var err error
		input, err = io.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
	}

	var unmarshalled any
	err := json.Unmarshal(input, &unmarshalled)
	if err != nil {
		panic(err)
	}

	output, err := json.MarshalIndent(unmarshalled, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Print(string(output))
}
