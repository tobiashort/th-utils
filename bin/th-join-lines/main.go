package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/tobiashort/th-utils/lib/clap"
	"github.com/tobiashort/th-utils/lib/must"
	"github.com/tobiashort/th-utils/lib/slices"
	"github.com/tobiashort/th-utils/lib/zip"
)

type Args struct {
	Files     []string `clap:"mandatory,positional,desc='The files which lines shall be joined.'"`
	Separator string   `clap:"default=' ',desc='The separator how lines shall be joined.'"`
}

func join(args Args) {
	lines := [][]string{}

	for _, f := range args.Files {
		content := string(must.Do2(os.ReadFile(f)))
		content = strings.ReplaceAll(content, "\r", "")
		lines = append(lines, strings.Split(content, "\n"))
	}

	if len(lines) == 0 {
		return
	}

	if len(lines) == 1 {
		fmt.Print(strings.Join(lines[0], args.Separator))
		return
	}

	joined := slices.Reduce(lines[1:], lines[0], func(a, b []string) []string {
		res := zip.Zip(struct{ A, B string }{}, a, b)
		return slices.Map(res, func(s struct{ A, B string }) string { return s.A + args.Separator + s.B })
	})

	for _, line := range joined {
		fmt.Println(line)
	}
}

func main() {
	args := Args{}
	clap.Description(`This tool combines the contents of one or more input files by joining their
lines in a structured, line-by-line manner. When multiple files are provided,
the tool merges corresponding lines from each file together, creating a new
output where each line consists of the combined content from the same line
position across all input files. When a single file is provided, the tool
combines the lines within that file into a single continuous joined sequence`)
	clap.Parse(&args)
	join(args)
}
