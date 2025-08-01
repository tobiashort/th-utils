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
		fmt.Print(`Usage: starts-with

Reads from stdin and prints all lines that start
with the given prefix.

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

	prefixes := make([]string, 0)
	prefixes = append(prefixes, matchFlag)
	for _, prefix := range orMatchFlags {
		prefixes = append(prefixes, prefix)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSuffix(line, "\n")
		matches := false
		for _, prefix := range prefixes {
			if caseSensitiveFlag {
				if strings.HasPrefix(line, prefix) {
					matches = true
					break
				}
			} else {
				if strings.HasPrefix(strings.ToLower(line), strings.ToLower(prefix)) {
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
