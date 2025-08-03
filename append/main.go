package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/tobiashort/clap-go"
)

type Args struct {
	Suffix string `clap:"positional,mandatory,description='the string that should be appended'"`
}

func main() {
	args := Args{}
	clap.Description("Reads from Stdin and appends each line with the given suffix.")
	clap.Parse(&args)

	suffix := args.Suffix
	suffixUnescaped := strings.Builder{}
	for i := 0; i < len(suffix); i++ {
		curr := suffix[i]
		if curr == '\\' && i+1 < len(suffix) {
			next := suffix[i+1]
			switch next {
			case 'n':
				suffixUnescaped.WriteByte('\n')
				i++
			case 'r':
				suffixUnescaped.WriteByte('\r')
				i++
			case 't':
				suffixUnescaped.WriteByte('\t')
				i++
			default:
				suffixUnescaped.WriteByte(curr)
				suffixUnescaped.WriteByte(next)
				i++
			}
		} else {
			suffixUnescaped.WriteByte(curr)
		}
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Printf("%s%s\n", text, suffixUnescaped.String())
	}
}
