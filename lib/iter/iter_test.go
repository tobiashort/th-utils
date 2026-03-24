package iter_test

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/tobiashort/th-utils/lib/iter"
	"github.com/tobiashort/th-utils/lib/must"
)

func assertEqual(actual, expected any, t *testing.T) {
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected %+v, got %+v\n", expected, actual)
	}
}

func Test(t *testing.T) {
	slice := []string{"Numbers", "1", "2", "3", "", "", "", "4", "5"}
	iterator := iter.From(slice)
	iterator = iter.Skip(iterator, 1)
	assertEqual(iter.Collect(iterator), []string{"1", "2", "3", "", "", "", "4", "5"}, t)
	iterator = iter.Filter(iterator, func(item string) bool { return item != "" })
	assertEqual(iter.Collect(iterator), []string{"1", "2", "3", "4", "5"}, t)
	iterator2 := iter.Map(iterator, func(item string) int { return must.Do2(strconv.Atoi(item)) })
	assertEqual(iter.Collect(iterator2), []int{1, 2, 3, 4, 5}, t)
	sum := iter.Reduce(iterator2, 0, func(a, b int) int { return a + b })
	assertEqual(sum, 15, t)
}
