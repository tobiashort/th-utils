package strings

import (
	stdstrings "strings"
)

func AllIndexes(s string, substr string) []int {
	indexes := []int{}
	for pos := 0; ; {
		index := stdstrings.Index(s[pos:], substr)
		if index < 0 {
			break
		}
		index = index + pos
		indexes = append(indexes, index)
		pos = index + len(substr)
	}
	return indexes
}
