package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/tobiashort/clap-go"
	"github.com/tobiashort/orderedmap-go"
	. "github.com/tobiashort/utils-go/must"

	"github.com/tobiashort/th-utils/pkg/ellipsis"
)

type Args struct {
	Count           bool `clap:"description='The count of the number of times the line occurred'"`
	Plot            bool `clap:"conflicts-with=Count,description='Plot as horizontal ascii bar chart'"`
	PlotLabelWidth  int  `clap:"short=,default-value=10,description='The label width when plotting.'"`
	PlotMaxBarWidth int  `clap:"short=,default-value=80,description='The max bar width when plotting.'"`
}

func main() {
	args := Args{}
	clap.Parse(&args)

	if args.PlotLabelWidth < 3 {
		args.PlotLabelWidth = 3
	}

	if args.PlotMaxBarWidth < 0 {
		args.PlotMaxBarWidth = 0
	}

	keywordCounts := orderedmap.NewOrderedMap[string, int]()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		count, _ := keywordCounts.Get(text)
		keywordCounts.Put(text, count+1)
	}
	Must(scanner.Err())

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
		for value, count := range keywordCounts.Iterate() {
			barWidth := int(float64(args.PlotMaxBarWidth) * float64(count) / float64(maxCount))
			bar := strings.Repeat("\u28FF", barWidth/2)
			bar += strings.Repeat("\u2847", barWidth%2)
			fmt.Printf("%*s %s %d\n", args.PlotLabelWidth, ellipsis.Ellipsis(value, args.PlotLabelWidth), bar, count)
		}
		return
	}

	for _, value := range keywordCounts.Keys() {
		fmt.Println(value)
	}
}
