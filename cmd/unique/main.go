package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"

	"github.com/tobiashort/clap-go"
	"github.com/tobiashort/orderedmap-go"
)

type Args struct {
	Count bool `clap:"description='The count of the number of times the line occurred'"`
	Plot  bool `clap:"description='Plot as horizontal ascii bar chart'"`
}

func main() {
	args := Args{}
	clap.Parse(&args)

	keywordCounts := orderedmap.NewOrderedMap[string, int]()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		count, _ := keywordCounts.Get(text)
		keywordCounts.Put(text, count+1)
	}

	if args.Count {
		maxCount := slices.Max(keywordCounts.Values())
		maxCountWidth := len(fmt.Sprintf("%d", maxCount))
		for value, count := range keywordCounts.Iterate() {
			fmt.Printf("%*d %s\n", maxCountWidth, count, value)
		}
		return
	}

	if args.Plot {
		maxCount := slices.Max(keywordCounts.Values())
		maxCountWidth := len(fmt.Sprintf("%d", maxCount))

		labelWidths := make([]int, 0)
		for value, _ := range keywordCounts.Iterate() {
			labelWidths = append(labelWidths, len(value))
		}

		maxLabelWidth := slices.Max(labelWidths)

		maxBarWidth := 80 - maxCountWidth - maxLabelWidth - 2

		for value, count := range keywordCounts.Iterate() {
		}

		return
	}

	for value := range keywordCounts.Keys() {
		fmt.Println(value)
	}
}
