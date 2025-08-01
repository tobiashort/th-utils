package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

type stringList []string

func (l *stringList) String() string {
	return fmt.Sprintf("%v", *l)
}

func (l *stringList) Set(value string) error {
	*l = append(*l, value)
	return nil
}

var (
	matchFlag         string
	orMatchFlags      stringList
	invertFlag        bool
	caseSensitiveFlag bool
)

func main() {
	flag.Usage = func() {
		fmt.Print(`Usage: contains

Reads from stdin and prints all lines that contains the given string.

`)
		flag.PrintDefaults()
	}

	flag.StringVar(&matchFlag, "match", "", "the pattern to match")
	flag.StringVar(&matchFlag, "m", "", "alias for -match")
	flag.Var(&orMatchFlags, "or-match", "adds an additional pattern to match")
	flag.Var(&orMatchFlags, "or", "alias for -or-match")
	flag.BoolVar(&invertFlag, "invert", false, "inverts the logic.")
	flag.BoolVar(&invertFlag, "i", false, "alias for -invert")
	flag.BoolVar(&caseSensitiveFlag, "case-sensitive", false, "case sensitive match")
	flag.BoolVar(&caseSensitiveFlag, "case", false, "alias for -case-sensitive")
	flag.Parse()

	if matchFlag == "" {
		flag.Usage()
		os.Exit(1)
	}

	matchStrings := make([]string, 0)
	matchStrings = append(matchStrings, matchFlag)
	for _, matchString := range orMatchFlags {
		matchStrings = append(matchStrings, matchString)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSuffix(line, "\n")
		matches := false
		for _, matchString := range matchStrings {
			if caseSensitiveFlag {
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
		if invertFlag {
			if !matches {
				fmt.Println(line)
			}
		} else {
			if matches {
				fmt.Println(line)
			}
		}
	}
}
