package main

import (
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/tobiashort/clap-go"
)

type Args struct {
	Reverse bool `clap:"description='Reverses the sort order'"`
}

func main() {
	args := Args{}
	clap.Parse(&args)

	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(data), "\n")
	slices.SortFunc(lines, func(a, b string) int {
		if args.Reverse {
			return len(a) - len(b)
		}
		return len(b) - len(a)
	})

	fmt.Println(strings.Join(lines, "\n"))
}
