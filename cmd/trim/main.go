package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/tobiashort/clap-go"
)

type Args struct {
	Left  bool `clap:"conflicts-with='Right',description='Only trim leading whitespace.'"`
	Right bool `clap:"conflicts-with='Left',description='Only trim trailing whitespace.'"`
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
		} else {
			text = strings.TrimLeftFunc(text, unicode.IsSpace)
			text = strings.TrimRightFunc(text, unicode.IsSpace)
		}
		fmt.Println(text)
	}
}
