package main

import (
	"fmt"
	"html"
	"io"
	"os"

	"github.com/tobiashort/clap-go"
)

type Args struct {
	String string `clap:"positional,description='The string to decode. Reads from Stdin if not specified.'"`
}

func main() {
	args := Args{}
	clap.Parse(&args)

	var input string
	if args.String != "" {
		input = args.String
	} else {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
		input = string(data)
	}

	fmt.Print(html.UnescapeString(input))
}
