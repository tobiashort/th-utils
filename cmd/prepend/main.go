package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/tobiashort/clap-go"
)

type Args struct {
	Prefix string `clap:"positional,mandatory,description='The prefix to append'"`
}

func main() {
	args := Args{}
	clap.Description("Reads from Stdin and prepends each line with the given prefix.")
	clap.Parse(&args)

	prefix := args.Prefix
	prefixUnescaped := strings.Builder{}
	for i := 0; i < len(prefix); i++ {
		curr := prefix[i]
		if curr == '\\' && i+1 < len(prefix) {
			next := prefix[i+1]
			switch next {
			case 'n':
				prefixUnescaped.WriteByte('\n')
				i++
			case 'r':
				prefixUnescaped.WriteByte('\r')
				i++
			case 't':
				prefixUnescaped.WriteByte('\t')
				i++
			default:
				prefixUnescaped.WriteByte(curr)
				prefixUnescaped.WriteByte(next)
				i++
			}
		} else {
			prefixUnescaped.WriteByte(curr)
		}
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Printf("%s%s\n", prefixUnescaped.String(), text)
	}
}
