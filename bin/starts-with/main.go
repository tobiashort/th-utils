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
	Match         string   `clap:"positional,mandatory,desc='The prefix to match'"`
	OrMatch       []string `clap:"desc='adds an additional pattern to match'"`
	Invert        bool     `clap:"desc='inverts the logic'"`
	CaseSensitive bool     `clap:"desc='case sensitive match'"`
}

func main() {
	args := Args{}
	clap.Description("Reads from stdin and prints all lines that starts with the given prefix.")
	clap.Parse(&args)

	prefixes := make([]string, 0)
	prefixes = append(prefixes, args.Match)
	prefixes = append(prefixes, args.OrMatch...)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSuffix(line, "\n")
		matches := false
		for _, prefix := range prefixes {
			if args.CaseSensitive {
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
