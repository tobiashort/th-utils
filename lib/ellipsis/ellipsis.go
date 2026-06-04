package ellipsis

import (
	"unicode/utf8"

	"github.com/tobiashort/th-utils/lib/ansi"
)

func Ellipsis(text string, length int) string {
	return EllipsisSuffix(text, length, "...")
}

func EllipsisSuffix(text string, length int, suffix string) string {
	if utf8.RuneCountInString(text) <= length {
		return text
	}
	return string([]rune(text)[:length-utf8.RuneCountInString(ansi.Strip(suffix))]) + suffix
}
