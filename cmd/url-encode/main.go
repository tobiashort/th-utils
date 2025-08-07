package main

import (
	"fmt"
	"io"
	"net/url"
	"os"

	"github.com/tobiashort/clap-go"
)

type Args struct {
	PathEscape  bool   `clap:"long=path,conflicts-with='QueryEscape,EscapeAll',description='Escapes the string so it can be safely placed inside a URL path segment'"`
	QueryEscape bool   `clap:"long=query,conflicts-with='PathEscape,EscapeAll',description='Escapes the string so it can be safely placed inside a URL query'"`
	EscapeAll   bool   `clap:"short=a,long=all,conflicts-with='PathEscape,QueryEscape',description='Escape all characters'"`
	String      string `clap:"positional,description='The string to encode. Reads from Stdin if not specified.'"`
}

func main() {
	args := Args{}
	clap.Parse(&args)

	var input string
	if args.String != "" {
		input = args.String
	} else {
		read, err := io.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
		input = string(read)
	}

	var escaped string
	if args.EscapeAll {
		for i := range len(input) {
			escaped += fmt.Sprintf("%%%02X", input[i])
		}
	} else if args.PathEscape {
		escaped = url.PathEscape(input)
	} else if args.QueryEscape {
		escaped = url.QueryEscape(input)
	} else {
		fmt.Fprintln(os.Stderr, "You need to specifiey how to encode: path, query or all")
		os.Exit(1)
	}

	fmt.Print(escaped)
}
