package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/tobiashort/clap-go"
	"github.com/tobiashort/th-utils/pkg/unescape"
	"github.com/tobiashort/utils-go/must"
)

type Args struct {
	OldString string `clap:"positional,mandatory,description='The string to replace'"`
	NewString string `clap:"positional,mandatory,description='The strint that replaces OldString'"`
}

func main() {
	args := Args{}
	clap.Description("Reads from Stdin and transforms the string by replacing all occurrences of OldString with NewString.")
	clap.Parse(&args)

	oldText := string(must.Do2(io.ReadAll(os.Stdin)))
	newText := strings.ReplaceAll(oldText, unescape.Unescape(args.OldString), unescape.Unescape(args.NewString))
	fmt.Println(newText)
}
