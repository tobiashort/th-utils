package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/tobiashort/clap-go"
)

type Args struct {
	JSON string `clap:"positional,description='The JSON string. Reads from Stdin if not specified.'"`
}

func main() {
	args := Args{}
	clap.Parse(&args)

	var input []byte
	if args.JSON == "" {
		var err error
		input, err = io.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
	} else {
		input = []byte(args.JSON)
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
