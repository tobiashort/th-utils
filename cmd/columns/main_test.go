package main

import (
	"strings"
	"testing"
)

var trim = strings.TrimSpace

func detend(in string) string {
	lines := strings.Split(in, "\n")
	for idx := range lines {
		lines[idx] = trim(lines[idx])
	}
	return strings.Join(lines, "\n")
}

func visibleWhitespace(in string) string {
	out := strings.ReplaceAll(in, " ", ".")
	out = strings.ReplaceAll(out, "\n", "↵\n")
	return out
}

func testTextToCols(t *testing.T, in string, nCols int, expected string) {
	actual := textToCols(in, nCols)

	if actual != expected {
		t.Fatalf(
			"\nExected:\n%s\nActual:\n%s", 
			visibleWhitespace(expected),
			visibleWhitespace(actual))
	}
}

func TestTextToCols1(t *testing.T) {
	in := detend(`1
                  2
                  3
                  4
                  5
                  6
                  7
                  8
                  9
                  10
                  11
                  12
				  `)

	expected := detend(`1    4    7    10
                        2    5    8    11
                        3    6    9    12
						`)

	testTextToCols(t, in, 4, expected)
}
