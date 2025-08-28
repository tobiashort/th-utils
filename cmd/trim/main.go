package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/tobiashort/clap-go"
	"github.com/tobiashort/th-utils/pkg/unescape"
)

type Args struct {
	Left   bool   `clap:"conflicts-with='Right,Prefix,Suffix',description='Only trim leading whitespace.'"`
	Right  bool   `clap:"conflicts-with='Left,Prefix,Suffix',description='Only trim trailing whitespace.'"`
	Prefix string `clap:"conflicts-with='Right,Left,Suffix',description='Trims the specified prefix.'"`
	Suffix string `clap:"conflicts-with='Right,Left,Prefix',description='Trims the specified suffix.'"`
}

func main() {
	args := Args{}
	clap.Description("Reads from Stdin and removes leading and/or trailing whitespaces for each line.")
	clap.Parse(&args)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if args.Left {
			text = strings.TrimLeftFunc(text, unicode.IsSpace)
		} else if args.Right {
			text = strings.TrimRightFunc(text, unicode.IsSpace)
		} else if args.Prefix != "" {
			text = strings.TrimPrefix(text, unescape.Unescape(args.Prefix))
		} else if args.Suffix != "" {
			text = strings.TrimSuffix(text, unescape.Unescape(args.Suffix))
		} else {
			text = strings.TrimLeftFunc(text, unicode.IsSpace)
			text = strings.TrimRightFunc(text, unicode.IsSpace)
		}
		fmt.Println(text)
	}
}
