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
	Reverse bool `clap:"desc='Reverses the sort order'"`
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
			return len([]rune(a)) - len([]rune(b))
		}
		return len([]rune(b)) - len([]rune(a))
	})

	fmt.Println(strings.Join(lines, "\n"))
}
