package strings_test

import (
	"reflect"
	"testing"

	"github.com/tobiashort/th-utils/lib/strings"
)

func TestAllIndexes(t *testing.T) {
	text := "data aa breach aa offline aa bearer"
	indexes := strings.AllIndexes(text, "aa")
	expected := []int{5, 15, 26}
	if !reflect.DeepEqual(indexes, expected) {
		t.Fatalf("expected: %v, got: %v", expected, indexes)
	}
}
