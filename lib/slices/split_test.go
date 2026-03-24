package slices_test

import (
	"reflect"
	"testing"

	"github.com/tobiashort/th-utils/lib/slices"
	"github.com/tobiashort/th-utils/lib/strings"
)

func TestSplit(t *testing.T) {
	type Int struct {
		V int
	}

	ints := []Int{
		{V: 1},
		{V: 0},
		{V: 0},
		{V: 1},
		{V: 0},
		{V: 1},
		{V: 0},
		{V: 0},
		{V: 0},
		{V: 1},
	}

	actual := slices.Split(ints, func(i Int) bool { return i.V == 1 })

	expected := [][]Int{
		{},
		{
			{V: 0},
			{V: 0},
		},
		{
			{V: 0},
		},
		{
			{V: 0},
			{V: 0},
			{V: 0},
		},
		{},
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf(strings.Dedent(`expected %v
		                        |got %q`), expected, actual)

	}
}

func TestSplit2(t *testing.T) {
	type Int struct {
		V int
	}

	ints := []Int{}

	actual := slices.Split(ints, func(i Int) bool { return i.V == 1 })

	expected := [][]Int{}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected: %q, got %q", expected, actual)

	}
}
