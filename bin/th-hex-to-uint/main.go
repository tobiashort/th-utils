package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/tobiashort/th-utils/lib/clap"
	"github.com/tobiashort/th-utils/lib/must"
)

type Args struct {
	Hex string `clap:"positional,desc='The hex value to convert. Reads from Stdin if not specified.'"`
}

func main() {
	args := Args{}
	clap.Parse(&args)

	hex := args.Hex
	if hex == "" {
		hex = string(must.Do2(io.ReadAll(os.Stdin)))
	}

	hex = strings.TrimPrefix(hex, "0x")

	i := must.Do2(strconv.ParseUint(hex, 16, 64))
	fmt.Print(i)
}
