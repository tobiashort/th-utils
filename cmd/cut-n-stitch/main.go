package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/tobiashort/clap-go"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func must2[T any](val T, err error) T {
	must(err)
	return val
}

type Args struct {
	Delimiter string `clap:"positional,mandatory,description='The delimiter where a given line from Stdin shall be cut.'"`
	Format    string `clap:"positional,mandatory,description='The format how the cut line shall be stitched together'"`
}

func main() {
	args := Args{}
	clap.Example(`$ echo "left-middle-right" | cut-n-stitch -- "-"" "{{ index . 0 }}-{{ index . 2}}"
left-right`)
	clap.Parse(&args)

	delimiter := args.Delimiter
	format := must2(template.New("").Parse(fmt.Sprintf("%s\n", args.Format)))
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		cut := strings.Split(line, delimiter)
		must(format.Execute(os.Stdout, cut))
	}
}
