package main

import (
	"fmt"
	"io"
	"os"

	"github.com/tobiashort/clap-go"
	"github.com/tobiashort/utils-go/must"
)

type Args struct {
	String string `clap:"positional,description='The input string. Reads from Stdin if not specified.'"`
	Start  int    `clap:"default-value=0,description='Start of the substring.'"`
	End    int    `clap:"mandatory,description='End of the substring (exclusive).'"`
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
