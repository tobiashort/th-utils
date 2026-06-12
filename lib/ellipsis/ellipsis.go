package ellipsis

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/tobiashort/th-utils/lib/ansi"
	"github.com/tobiashort/th-utils/lib/slices"
	"github.com/tobiashort/th-utils/lib/unicode"
)

const (
	PosStart int = iota
	PosCenter
	PosEnd
)

const (
	tokenRune = "RUNE"
	tokenAnsi = "ANSI"
)

type token struct {
	kind    string
	literal string
	width   int
}

func tokenize(s string) []token {
	tokens := make([]token, 0)
	rg := regexp.MustCompile("^" + ansi.Regexp.String())
	for i := 0; i < len(s); {
		a := rg.FindString(s[i:])
		if a != "" {
			tokens = append(tokens, token{kind: tokenAnsi, literal: a, width: unicode.WidthString(a)})
			i += len(a)
		} else {
			r, rs := utf8.DecodeRuneInString(s[i:])
			tokens = append(tokens, token{kind: tokenRune, literal: string(r), width: unicode.Width(r)})
			i += rs
		}
	}
	return tokens
}

func Ellipsis(text string, width int, cutString string, pos int) string {
	textTokens := tokenize(text)
	visualTextLen := 0
	for _, t := range textTokens {
		if t.kind == tokenRune {
			visualTextLen += t.width
		}
	}

	if visualTextLen <= width {
		return text
	}

	cutStringTokens := tokenize(cutString)
	visualCutStringLen := 0
	for _, t := range cutStringTokens {
		if t.kind == tokenRune {
			visualCutStringLen += t.width
		}
	}

	if width <= visualCutStringLen {
		return cutString
	}

	ellipsed := make([]token, 0)

	switch pos {
	case PosStart:
		index := len(textTokens) - 1
		count := 0
		fill := 0
		for ; index > 0; index-- {
			t := textTokens[index]
			if t.kind == tokenAnsi {
				continue
			}
			count += t.width
			if count == (width - visualCutStringLen) {
				break
			}
			if count > (width - visualCutStringLen) {
				index++
				fill = count - (width - visualCutStringLen)
				break
			}
		}
		left := textTokens[:index]
		left = slices.Filter(left, func(t token) bool { return t.kind == tokenAnsi })
		for range fill {
			left = append([]token{{kind: tokenRune, literal: " ", width: 1}}, left...)
		}
		right := textTokens[index:]
		ellipsed = append(ellipsed, cutStringTokens...)
		ellipsed = append(ellipsed, left...)
		ellipsed = append(ellipsed, right...)
	case PosCenter:
		leftLen := width / 2
		leftCutString := cutString[:len(cutString)/2]
		rightLen := (width / 2) + (width % 2)
		rightCutString := cutString[len(cutString)/2:]
		left := Ellipsis(text, leftLen, leftCutString, PosEnd)
		right := Ellipsis(text, rightLen, rightCutString, PosStart)
		return left + right
	case PosEnd:
		index := 0
		count := 0
		fill := 0
		for ; index < len(textTokens); index++ {
			if count == (width - visualCutStringLen) {
				break
			}
			if count > (width - visualCutStringLen) {
				index--
				fill = count - (width - visualCutStringLen)
				break
			}
			t := textTokens[index]
			if t.kind == tokenAnsi {
				continue
			}
			count += t.width
		}
		left := textTokens[:index]
		right := textTokens[index:]
		right = slices.Filter(right, func(t token) bool { return t.kind == tokenAnsi })
		for range fill {
			right = append(right, token{kind: tokenRune, literal: " ", width: 1})
		}
		ellipsed = append(ellipsed, left...)
		ellipsed = append(ellipsed, right...)
		ellipsed = append(ellipsed, cutStringTokens...)
	}

	sb := strings.Builder{}
	for _, t := range ellipsed {
		sb.WriteString(t.literal)
	}
	return sb.String()
}
