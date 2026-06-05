package ellipsis

import (
	"unicode/utf8"

	"github.com/tobiashort/th-utils/lib/ansi"
)

func Ellipsis(text string, length int) string {
	return EllipsisSuffix(text, length, "...")
}

func EllipsisSuffix(text string, length int, suffix string) string {
	stripped := ansi.Strip(text)
	if utf8.RuneCountInString(stripped) <= length {
		return text
	}
	lenDiff := utf8.RuneCountInString(text) - utf8.RuneCountInString(stripped)
	return string([]rune(text)[:length+lenDiff-utf8.RuneCountInString(ansi.Strip(suffix))]) + suffix
}
