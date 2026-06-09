package main

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/tobiashort/th-utils/lib/ansi"
)

type TokenType string

const (
	TokenRune TokenType = "RUNE"
	TokenAnsi TokenType = "ANSI"
)

type Token struct {
	Type    TokenType
	Literal string
}

func Parse(text string) [][]Token {
	text = strings.ReplaceAll(text, "\r", "")
	text = strings.TrimSuffix(text, "\n")
	text = strings.ReplaceAll(text, "\t", "    ")
	textLines := strings.Split(text, "\n")
	height := len(textLines)
	tokens := make([][]Token, height)
	rg := regexp.MustCompile("^" + ansi.Regexp.String())
	for i, textLine := range textLines {
		tokens[i] = make([]Token, 0)
		for j := 0; j < len(textLine); {
			a := rg.FindString(textLine[j:])
			if a != "" {
				tokens[i] = append(tokens[i], Token{Type: TokenAnsi, Literal: a})
				j += len(a)
			} else {
				r, rs := utf8.DecodeRuneInString(textLine[j:])
				tokens[i] = append(tokens[i], Token{Type: TokenRune, Literal: string(r)})
				j += rs
			}
		}
	}
	return tokens
}
