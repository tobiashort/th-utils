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
	Match         string   `clap:"positional,mandatory,description='The suffix to match'"`
	OrMatch       []string `clap:"description='adds an additional pattern to match'"`
	Invert        bool     `clap:"description='inverts the logic'"`
	CaseSensitive bool     `clap:"description='case sensitive match'"`
}

func main() {
	args := Args{}
	clap.Description("Reads from stdin and prints all lines that ends with the given suffix.")
	clap.Parse(&args)

	suffixes := make([]string, 0)
	suffixes = append(suffixes, args.Match)
	suffixes = append(suffixes, args.OrMatch...)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSuffix(line, "\n")
		matches := false
		for _, suffix := range suffixes {
			if args.CaseSensitive {
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
