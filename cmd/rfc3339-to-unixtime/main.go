package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/tobiashort/clap-go"
)

func invalidInput(input string) {
	fmt.Fprintf(os.Stderr, "the provided input is not in RFC3339 format, e.g. 1970-01-01T00:00:00Z: %s", input)
	os.Exit(1)
}

type Args struct {
	RFC3339 string `clap:"positional,description='The datetime in RFC3339 format, e.g. 1970-01-01T00:00:00Z. Reads from Stdin if not specified.'"`
}

func main() {
	args := Args{}
	clap.Parse(&args)

	var input string
	if args.RFC3339 != "" {
		input = args.RFC3339
	} else {
		read, err := io.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
		input = string(read)
	}

	rfc3339, err := time.Parse(time.RFC3339, input)
	if err != nil {
		invalidInput(input)
		return
	}

	fmt.Println(rfc3339.Unix())
}
