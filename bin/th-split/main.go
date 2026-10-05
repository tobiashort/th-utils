package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/tobiashort/th-utils/lib/clap"
	"github.com/tobiashort/th-utils/lib/iter"
	"github.com/tobiashort/th-utils/lib/must"
	"github.com/tobiashort/th-utils/lib/slices"
)

type Args struct {
	Delimiter string `clap:"desc='The delimiter that defines where to split.',conflicts='Chars'"`
	Chars     int    `clap:"short='n',desc='The number of chars after it should split'"`
}

func run(args Args) int {

	bs := must.Do2(io.ReadAll(os.Stdin))

	if args.Delimiter != "" {
		split := strings.Split(string(bs), args.Delimiter)
		fmt.Print(strings.Join(split, "\n"))
		return 0
	}

	if args.Chars > 0 {
		fmt.Print(
			iter.
				From(slices.Chunks(bs, args.Chars)).
				Map[string](func(a []byte) string { return string(a) }).
				Reduce("", func(a, b string) string { return a + "\n" + b }))
	}

	fmt.Fprintln(os.Stderr, "Don't know where to split")
	return 1
}

func main() {
	args := Args{}
	clap.Parse(&args)
	os.Exit(run(args))
}
