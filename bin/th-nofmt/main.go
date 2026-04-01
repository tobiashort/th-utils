package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/tobiashort/th-utils/lib/clap"
	"github.com/tobiashort/th-utils/lib/must"
	strings2 "github.com/tobiashort/th-utils/lib/strings"
)

type Args struct {
	File string `clap:"positional,desc='The file to format. Reads from Stdin if not specified.'"`
}

func main() {
	args := Args{}
	clap.Example(strings2.Dedent(`//nofmt:enable
                                 |[custom formatted code]
                                 |//nofmt:disable`))
	clap.Parse(&args)

	var replacements [][]string
	var replacement []string
	var enabled bool

	var src string
	if args.File != "" {
		src = string(must.Do2(os.ReadFile(args.File)))
	} else {
		src = string(must.Do2(io.ReadAll(os.Stdin)))
	}

	scanner := bufio.NewScanner(strings.NewReader(src))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "//nofmt:enable" {
			enabled = true
		} else if strings.TrimSpace(line) == "//nofmt:disable" {
			enabled = false
			replacements = append(replacements, replacement)
			replacement = make([]string, 0)
		} else {
			if enabled {
				replacement = append(replacement, line)
			}
		}
	}
	replacements = append(replacements, replacement)
	replacement = make([]string, 0)

	cmd := exec.Command("goimports")
	cmd.Stdin = strings.NewReader(src)
	goimportsOut := string(must.Do2(cmd.CombinedOutput()))

	enabled = false
	replacementIndex := 0
	scanner = bufio.NewScanner(strings.NewReader(goimportsOut))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "//nofmt:enable" {
			fmt.Println(line)
			enabled = true
			replacement := replacements[replacementIndex]
			for _, replacementLine := range replacement {
				fmt.Println(replacementLine)
			}
			replacementIndex++
		} else if strings.TrimSpace(line) == "//nofmt:disable" {
			fmt.Println(line)
			enabled = false
		} else {
			if enabled {
				continue
			} else {
				fmt.Println(line)
			}
		}
	}
}
