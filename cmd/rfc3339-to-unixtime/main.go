package main

import (
	"fmt"
	"os"
	"time"

	"github.com/tobiashort/clap-go"
)

type Args struct {
	RFC3339 string `clap:"positional,mandatory,description='The datetime in RFC3339 format, e.g. 1970-01-01T00:00:00Z'"`
}

func main() {
	args := Args{}
	clap.Parse(&args)

	input := args.RFC3339
	rfc3339, err := time.Parse(time.RFC3339, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "the provided input is not in RFC3339 format, e.g. 1970-01-01T00:00:00Z: %s", input)
		os.Exit(1)
	}

	fmt.Println(rfc3339.Unix())
}
