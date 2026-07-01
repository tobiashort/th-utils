package main

import (
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/tobiashort/th-utils/lib/clap"
	"github.com/tobiashort/th-utils/lib/must"
)

type Args struct {
	Head  int    `clap:"short=H,desc='Number of beginning lines to take'"`
	Tail  int    `clap:"desc='Number of ending lines to take'"`
	Lines string `clap:"desc='Specific lines to take, e.g. 2,3,5-10,24'"`
	File  string `clap:"positional,desc='The file to read. Reads from Stdin if not specified.'"`
}

func linesToTake(head int, lines string, tail int, maxLines int) []int {
	linesToSkip := []int{}

	if head > 0 {
		for i := range head {
			linesToSkip = append(linesToSkip, i+1)
		}
	}

	type Range struct {
		From int
		To   int
	}

	ranges := []Range{}
	if lines != "" {
		split1 := strings.SplitSeq(lines, ",")
		for s := range split1 {
			split2 := strings.Split(s, "-")
			if len(split2) == 1 {
				from := must.Do2(strconv.Atoi(split2[0]))
				ranges = append(ranges, Range{From: from, To: from})
			} else if len(split2) == 2 {
				from := must.Do2(strconv.Atoi(split2[0]))
				to := must.Do2(strconv.Atoi(split2[1]))
				ranges = append(ranges, Range{From: from, To: to})

			} else {
				panic("cannot parse lines argument")
			}
		}
	}

	for _, r := range ranges {
		for i := r.From; i <= r.To; i++ {
			linesToSkip = append(linesToSkip, i)
		}
	}

	if tail > 0 {
		for i := maxLines; i > maxLines-tail; i-- {
			linesToSkip = append(linesToSkip, i)
		}
	}

	slices.Sort(linesToSkip)
	linesToSkip = slices.Compact(linesToSkip)
	return linesToSkip
}

func run(args Args) {
	file := os.Stdin
	if args.File != "" {
		file = must.Do2(os.Open(args.File))
	}

	content := string(must.Do2(io.ReadAll(file)))
	content = strings.TrimSuffix(content, "\n")
	lines := strings.Split(content, "\n")
	linesToTake := linesToTake(args.Head, args.Lines, args.Tail, len(lines))

	for _, l := range linesToTake {
		fmt.Println(lines[l-1])
	}
}

func main() {
	args := Args{}
	clap.Parse(&args)
	run(args)
}
