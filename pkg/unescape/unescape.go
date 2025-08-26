package unescape

import "strings"

// Unescape replaces the escape sequences "\n", "\r", and "\t"
// in a string with their corresponding characters
func Unescape(text string) string {
	unescaped := strings.Builder{}
	for i := 0; i < len(text); i++ {
		curr := text[i]
		if curr == '\\' && i+1 < len(text) {
			next := text[i+1]
			switch next {
			case 'n':
				unescaped.WriteByte('\n')
				i++
			case 'r':
				unescaped.WriteByte('\r')
				i++
			case 't':
				unescaped.WriteByte('\t')
				i++
			default:
				unescaped.WriteByte(curr)
				unescaped.WriteByte(next)
				i++
			}
		} else {
			unescaped.WriteByte(curr)
		}
	}
	return unescaped.String()
}
