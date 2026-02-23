package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/tobiashort/clap-go"
	"github.com/tobiashort/utils-go/assert"
)

type Args struct {
	Match         string   `clap:"positional,mandatory,desc='the pattern to match'"`
	OrMatch       []string `clap:"desc='adds an additional pattern to match'"`
	Invert        bool     `clap:"desc='inverts the logic.'"`
	CaseSensitive bool     `clap:"desc=case sensitive match'"`
}

func main() {
	args := Args{}
	clap.Description("Reads from stdin and prints all lines that contains the given string.")
	clap.Parse(&args)

	matchStrings := make([]string, 0)
	matchStrings = append(matchStrings, args.Match)
	matchStrings = append(matchStrings, args.OrMatch...)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSuffix(line, "\n")
		matches := false
		for _, matchString := range matchStrings {
			if args.CaseSensitive {
				if strings.Contains(line, matchString) {
					matches = true
					break
				}
			} else {
				if strings.Contains(strings.ToLower(line), strings.ToLower(matchString)) {
					matches = true
					break
				}
			}
		}
		if args.Invert {
			if !matches {
				fmt.Println(line)
			}
		} else {
			if matches {
				fmt.Println(line)
			}
		}
	}

	assert.Nil(scanner.Err(), "scanner error: %w", scanner.Err())
}
