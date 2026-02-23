package main

import (
	"fmt"
	"io"
	"os"

	"github.com/tobiashort/th-utils/pkg/json"

	"github.com/tobiashort/clap-go"
	"github.com/tobiashort/utils-go/assert"
)

type Args struct {
	File string `clap:"positional,desc='JSON file to convert. Reads from Stdin if not specified'"`
}

func main() {
	args := Args{}
	clap.Parse(&args)

	assert.PanicFunc = func(err error) {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	var file *os.File
	if args.File == "" {
		file = os.Stdin
	} else {
		var err error
		file, err = os.Open(args.File)
		assert.Nil(err, "open file error: %w", err)
	}

	b, err := io.ReadAll(file)
	assert.Nil(err, "read file error: %w", err)

	s, err := json.ToGoStruct(b)
	assert.Nil(err, "conversion error: %w", err)

	fmt.Println(s)
}
