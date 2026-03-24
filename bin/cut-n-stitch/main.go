package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/tobiashort/clap-go"
	"github.com/tobiashort/utils-go/assert"
	"github.com/tobiashort/utils-go/must"
)

type Args struct {
	Delimiter string `clap:"positional,mandatory,desc='The delimiter where a given line from Stdin shall be cut.'"`
	Format    string `clap:"positional,mandatory,desc='The format how the cut line shall be stitched together'"`
}

func main() {
	args := Args{}
	clap.Example(`$ echo "left-middle-right" | cut-n-stitch -- "-"" "{{ index . 0 }}-{{ index . 2}}"
left-right`)
	clap.Parse(&args)

	delimiter := args.Delimiter
	format := must.Do2(template.New("").Parse(fmt.Sprintf("%s\n", args.Format)))
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		cut := strings.Split(line, delimiter)
		must.Do(format.Execute(os.Stdout, cut))
	}
	assert.Nil(scanner.Err(), "scanner error: %w", scanner.Err())
}
