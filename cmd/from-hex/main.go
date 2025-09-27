package main

import (
	"encoding/hex"
	"io"
	"os"
	"strings"

	"github.com/tobiashort/clap-go"
	"github.com/tobiashort/utils-go/must"
)

type Args struct {
	Hex string `clap:"positional,description='The hex code to decode. Reads from Stdin if not specified.'"`
}

func main() {
	args := Args{}
	clap.Parse(&args)

	var input io.Reader

	if args.Hex == "" {
		input = os.Stdin
	} else {
		input = strings.NewReader(args.Hex)
	}

	decoder := hex.NewDecoder(input)
	must.Do2(io.Copy(os.Stdout, decoder))
}
