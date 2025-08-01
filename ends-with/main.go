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
		fmt.Print(`Usage: ends-with

Reads from stdin and prints all lines that ends
with the given suffix.

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

	suffixes := make([]string, 0)
	suffixes = append(suffixes, matchFlag)
	for _, suffix := range orMatchFlags {
		suffixes = append(suffixes, suffix)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSuffix(line, "\n")
		matches := false
		for _, suffix := range suffixes {
			if caseSensitiveFlag {
				if strings.HasSuffix(line, suffix) {
					matches = true
					break
				}
			} else {
				if strings.HasSuffix(strings.ToLower(line), strings.ToLower(suffix)) {
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
