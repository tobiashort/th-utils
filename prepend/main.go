package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func usage() {
	fmt.Print(`Usage: prepend STRING

Reads from Stdin and prepends each line with the given STRING

EXAMPLE:
	$ echo "foobar" | prepend "prefix-"
	prefix-foobar
`)
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() != 1 {
		usage()
		os.Exit(1)
	}

	prefix := flag.Arg(0)
	prefixUnescaped := strings.Builder{}
	i := 0
	for i < len(prefix) {
		curr := prefix[i]
		if curr == '\\' && i+1 < len(prefix) {
			next := prefix[i+1]
			switch next {
			case 'n':
				prefixUnescaped.WriteByte('\n')
				i += 2
			case 'r':
				prefixUnescaped.WriteByte('\r')
				i += 2
			case 't':
				prefixUnescaped.WriteByte('\t')
				i += 2
			default:
				prefixUnescaped.WriteByte(curr)
				prefixUnescaped.WriteByte(next)
				i += 2
			}
		} else {
			prefixUnescaped.WriteByte(curr)
			i++
		}
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Printf("%s%s\n", prefixUnescaped.String(), text)
	}
}
