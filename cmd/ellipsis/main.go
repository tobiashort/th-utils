package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/tobiashort/clap-go"
	"github.com/tobiashort/utils-go/must"

	"github.com/tobiashort/th-utils/pkg/ellipsis"
)

type Args struct {
	Length int `clap:"default-value=20,description='Max length of the string including the three dots'"`
}

func main() {
	args := Args{}
	clap.Description("This program reads string from stdin and cuts ist at the specified length minus three and adds three dots.")
	clap.Parse(&args)

	if args.Length < 0 {
		fmt.Fprintln(os.Stderr, "length must be greather than 0")
		os.Exit(1)
	}

	bytesRead := must.Do2(io.ReadAll(os.Stdin))
	text := string(bytesRead)
	text = strings.TrimSpace(text)
	text = ellipsis.Ellipsis(text, args.Length)
	fmt.Print(text)
}
