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
		max := slices.Max(keywordCounts.Values())
		width := len(fmt.Sprintf("%d", max))
		for value, count := range keywordCounts.Iterate() {
			fmt.Printf("%*d %s\n", width, count, value)
		}
		return
	}

	if args.Plot {
		panic("not implemented")
	}

	for value := range keywordCounts.Keys() {
		fmt.Println(value)
	}
}
