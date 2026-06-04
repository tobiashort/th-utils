package main

import (
	"fmt"
	"io"
	"os"

	"github.com/tobiashort/th-utils/lib/ansi"
	"github.com/tobiashort/th-utils/lib/clap"
	"github.com/tobiashort/th-utils/lib/must"
)

type Args struct{}

func main() {
	args := Args{}
	clap.Parse(&args)

	in := string(must.Do2(io.ReadAll(os.Stdin)))
	stripped := ansi.Strip(in)
	fmt.Print(stripped)
}
