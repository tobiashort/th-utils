package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/tobiashort/th-utils/lib/clap"
)

type Args struct {
	Delimiter string `clap:"default=' => ',desc='The delimiter'"`
}

func main() {
	args := Args{}
	clap.Description("Parses tree outputs (tree, cargo tree, gradle dependencies) and produces paths. Reads from Stdin.")
	clap.Parse(&args)

	var path [65535]string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		depth := 0
		text := strings.TrimRight(scanner.Text(), "\n")
	parseText:
		if after, ok := strings.CutPrefix(text, "├── "); ok {
			text = after
			depth++
			path[depth] = text
		} else if after, ok := strings.CutPrefix(text, "+--- "); ok {
			text = after
			depth++
			path[depth] = text
		} else if after, ok := strings.CutPrefix(text, "└── "); ok {
			text = after
			depth++
			path[depth] = text
		} else if after, ok := strings.CutPrefix(text, `\--- `); ok {
			text = after
			depth++
			path[depth] = text
		} else if after, ok := strings.CutPrefix(text, "│   "); ok {
			text = after
			depth++
			goto parseText
		} else if after, ok := strings.CutPrefix(text, "│   "); ok {
			text = after
			depth++
			goto parseText
		} else if after, ok := strings.CutPrefix(text, "|    "); ok {
			text = after
			depth++
			goto parseText
		} else if after, ok := strings.CutPrefix(text, "     "); ok {
			text = after
			depth++
			goto parseText
		} else if after, ok := strings.CutPrefix(text, "    "); ok {
			text = after
			depth++
			goto parseText
		} else {
			path[0] = text
		}
		fmt.Println(strings.Join(path[:depth+1], args.Delimiter))
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
