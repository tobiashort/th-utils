package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/tobiashort/th-utils/lib/clap"
	"github.com/tobiashort/th-utils/lib/must"
	"github.com/tobiashort/th-utils/lib/unescape"
)

type Args struct {
	Regex     bool   `clap:"desc='Use regex"`
	OldString string `clap:"positional,mandatory,desc='The string to replace'"`
	NewString string `clap:"positional,mandatory,desc='The strint that replaces OldString'"`
}

func main() {
	args := Args{}
	clap.Description("Reads from Stdin and transforms the string by replacing all occurrences of OldString with NewString.")
	clap.Parse(&args)

	if args.Regex {
		oldText := regexp.MustCompile(args.OldString)
		newText := string(oldText.ReplaceAll(must.Do2(io.ReadAll(os.Stdin)), []byte(args.NewString)))
		fmt.Println(newText)
	} else {
		oldText := string(must.Do2(io.ReadAll(os.Stdin)))
		newText := strings.ReplaceAll(oldText, unescape.Unescape(args.OldString), unescape.Unescape(args.NewString))
		fmt.Println(newText)
	}
}
