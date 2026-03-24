package groupby_test

import (
	"reflect"
	"testing"

	"github.com/tobiashort/th-utils/lib/groupby"
)

func TestGroupBy(t *testing.T) {
	ints := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	grouped := groupby.GroupBy(ints, func(a, b int) bool { return a%2 == b%2 })
	if !reflect.DeepEqual(grouped, [][]int{{0, 2, 4, 6, 8, 10, 12, 14, 16}, {1, 3, 5, 7, 9, 11, 13, 15}}) {
		t.Error(grouped)
	}
}
