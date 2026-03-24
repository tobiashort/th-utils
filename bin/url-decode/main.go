package main

import (
	"fmt"
	"io"
	"net/url"
	"os"

	"github.com/tobiashort/clap-go"
)

type Args struct {
	PathUnescape  bool   `clap:"long=path,conflicts='QueryUnescape',desc='Unescapes the string from inside a URL path segment'"`
	QueryUnescape bool   `clap:"long=query,conflicts='PathUnescape',desc='Unescapes the string from inside a URL query'"`
	String        string `clap:"positional,desc='The string to dencode. Reads from Stdin if not specified.'"`
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

	var unescaped string
	var err error
	if args.PathUnescape {
		unescaped, err = url.PathUnescape(input)
	} else if args.QueryUnescape {
		unescaped, err = url.QueryUnescape(input)
	} else {
		fmt.Fprintln(os.Stderr, "You need to specifiy how the string shall be decoded: path or query")
		os.Exit(1)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	fmt.Print(unescaped)
}
