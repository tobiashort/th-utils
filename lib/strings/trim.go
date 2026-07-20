package strings

import (
	"strings"
	"unicode"
)

func TrimLeftSpace(s string) string {
	return strings.TrimLeftFunc(s, func(r rune) bool { return unicode.IsSpace(r) })
}

func TrimRightSpace(s string) string {
	return strings.TrimRightFunc(s, func(r rune) bool { return unicode.IsSpace(r) })
}
