package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/tobiashort/clap-go"
	"github.com/tobiashort/utils-go/assert"

	"github.com/tobiashort/th-utils/pkg/unescape"
)

type Args struct {
	Prefix string `clap:"positional,mandatory,description='The prefix to append'"`
}

func main() {
	args := Args{}
	clap.Description("Reads from Stdin and prepends each line with the given prefix.")
	clap.Parse(&args)

	prefix := args.Prefix
	prefixUnescaped := unescape.Unescape(prefix)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Printf("%s%s\n", prefixUnescaped, text)
	}
	assert.Nilf(scanner.Err(), "scanner error: %w", scanner.Err())
}
