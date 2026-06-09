package ellipsis

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/tobiashort/th-utils/lib/ansi"
	"github.com/tobiashort/th-utils/lib/slices"
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
}

func tokenize(s string) []token {
	tokens := make([]token, 0)
	rg := regexp.MustCompile("^" + ansi.Regexp.String())
	for i := 0; i < len(s); {
		a := rg.FindString(s[i:])
		if a != "" {
			tokens = append(tokens, token{kind: tokenAnsi, literal: a})
			i += len(a)
		} else {
			r, rs := utf8.DecodeRuneInString(s[i:])
			tokens = append(tokens, token{kind: tokenRune, literal: string(r)})
			i += rs
		}
	}
	return tokens
}

func Ellipsis(text string, length int, cutString string, pos int) string {
	textTokens := tokenize(text)
	visualTextLen := slices.Count(textTokens, func(t token) bool { return t.kind == tokenRune })

	if visualTextLen <= length {
		return text
	}

	cutStringTokens := tokenize(cutString)
	visualCutStringLen := slices.Count(cutStringTokens, func(t token) bool { return t.kind == tokenRune })

	if length <= visualCutStringLen {
		return cutString
	}

	ellipsed := make([]token, 0)

	switch pos {
	case PosStart:
		index := len(textTokens) - 1
		count := 0
		for ; index > 0; index-- {
			t := textTokens[index]
			if t.kind == tokenAnsi {
				continue
			}
			count++
			if count == (length - visualCutStringLen) {
				break
			}
		}
		left := textTokens[:index]
		left = slices.Filter(left, func(t token) bool { return t.kind == tokenAnsi })
		right := textTokens[index:]
		ellipsed = append(ellipsed, cutStringTokens...)
		ellipsed = append(ellipsed, left...)
		ellipsed = append(ellipsed, right...)
	case PosCenter:
		leftLen := length / 2
		leftCutString := cutString[:len(cutString)/2]
		rightLen := (length / 2) + (length % 2)
		rightCutString := cutString[len(cutString)/2:]
		left := Ellipsis(text, leftLen, leftCutString, PosEnd)
		right := Ellipsis(text, rightLen, rightCutString, PosStart)
		return left + right
	case PosEnd:
		index := 0
		count := 0
		for ; index < len(textTokens); index++ {
			if count == (length - visualCutStringLen) {
				break
			}
			t := textTokens[index]
			if t.kind == tokenAnsi {
				continue
			}
			count++
		}
		left := textTokens[:index]
		right := textTokens[index:]
		right = slices.Filter(right, func(t token) bool { return t.kind == tokenAnsi })
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
