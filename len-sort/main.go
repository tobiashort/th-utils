package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
)

func main() {
	var reverse bool

	flag.BoolVar(&reverse, "r", false, "reverse")
	flag.Parse()

	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(data), "\n")
	slices.SortFunc(lines, func(a, b string) int {
		if reverse {
			return len(a) - len(b)
		}
		return len(b) - len(a)
	})

	fmt.Println(strings.Join(lines, "\n"))
}
