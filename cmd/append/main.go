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
	Suffix string `clap:"positional,mandatory,description='the string that should be appended'"`
}

func main() {
	args := Args{}
	clap.Description("Reads from Stdin and appends each line with the given suffix.")
	clap.Parse(&args)

	suffix := args.Suffix
	suffixUnescaped := unescape.Unescape(suffix)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Printf("%s%s\n", text, suffixUnescaped)
	}
	assert.Nil(scanner.Err(), "scanner error: %w", scanner.Err())
}
