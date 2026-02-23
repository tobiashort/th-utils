package main

import (
	"fmt"
	"io"
	"os"

	"github.com/tobiashort/clap-go"
	"github.com/tobiashort/utils-go/must"
)

type Args struct {
	String string `clap:"positional,desc='The input string. Reads from Stdin if not specified.'"`
	Start  int    `clap:"default=0,desc='Start of the substring.'"`
	End    int    `clap:"mandatory,desc='End of the substring (exclusive).'"`
}

func main() {
	args := Args{}
	clap.Parse(&args)
	s := args.String
	if s == "" {
		s = string(must.Do2(io.ReadAll(os.Stdin)))
	}
	runes := []rune(s)
	fmt.Print(string(runes[args.Start:args.End]))
}
