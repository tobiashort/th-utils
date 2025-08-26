package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/tobiashort/clap-go"
	"github.com/tobiashort/th-utils/pkg/unescape"
)

type Args struct {
	OldString string `clap:"positional,mandatory,description='The string to replace'"`
	NewString string `clap:"positional,mandatory,description='The strint that replaces OldString'"`
}

func main() {
	args := Args{}
	clap.Description("Reads from Stdin and transforms the string by replacing all occurrences of OldString with NewString.")
	clap.Parse(&args)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		oldText := scanner.Text()
		newText := strings.ReplaceAll(oldText, unescape.Unescape(args.OldString), unescape.Unescape(args.NewString))
		fmt.Println(newText)
	}
}
