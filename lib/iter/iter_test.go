package iter_test

import (
	"fmt"
	"strconv"

	"github.com/tobiashort/th-utils/lib/iter"
	"github.com/tobiashort/th-utils/lib/must"
)

func Example() {
	fmt.Println(
		iter.From([]string{"Numbers", "1", "2", "3", "", "", "", "4", "5"}).
			Skip(1).
			Filter(func(item string) bool { return item != "" }).
			Map(func(item string) int { return must.Do2(strconv.Atoi(item)) }).
			Reduce(0, func(a, b int) int { return a + b }))
	// Output: 15
}
