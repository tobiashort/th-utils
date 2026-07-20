package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/tobiashort/th-utils/lib/assert"
	"github.com/tobiashort/th-utils/lib/clap"
	"github.com/tobiashort/th-utils/lib/must"
	strings2 "github.com/tobiashort/th-utils/lib/strings"
)

type Args struct {
	File string `clap:"positional,desc='The file to format. Reads from Stdin if not specified.'"`
}

func Fmt(src string) {
	var replacements [][]string
	var replacement []string
	var enabled bool

	srcLines := strings.Split(src, "\n")
	srcLines2 := make([]string, 0)
	for i := range srcLines {
		srcLine := srcLines[i]
		srcLine = strings2.TrimRightSpace(srcLine)
		if strings.HasSuffix(srcLine, "//nofmt") {
			srcLines2 = append(srcLines2, "//nofmt:enable")
			srcLines2 = append(srcLines2, strings.TrimSuffix(srcLine, " //nofmt"))
			srcLines2 = append(srcLines2, "//nofmt:disable")
		} else {
			srcLines2 = append(srcLines2, srcLine)
		}
	}

	for _, srcLine := range srcLines2 {
		if strings.TrimSpace(srcLine) == "//nofmt:enable" {
			enabled = true
		} else if strings.TrimSpace(srcLine) == "//nofmt:disable" {
			enabled = false
			replacements = append(replacements, replacement)
			replacement = make([]string, 0)
		} else {
			if enabled {
				replacement = append(replacement, srcLine)
			}
		}
	}
	replacements = append(replacements, replacement)
	replacement = make([]string, 0)

	src = strings.Join(srcLines2, "\n")

	cmd := exec.Command("goimports")
	cmd.Stdin = strings.NewReader(src)
	goimportsOut := string(must.Do2(cmd.CombinedOutput()))

	enabled = false
	replacementIndex := 0
	srcLines3 := make([]string, 0)
	scanner := bufio.NewScanner(strings.NewReader(goimportsOut))
	for scanner.Scan() {
		srcLine := scanner.Text()
		if strings.TrimSpace(srcLine) == "//nofmt:enable" {
			srcLines3 = append(srcLines3, srcLine)
			enabled = true
			replacement := replacements[replacementIndex]
			for _, replacementLine := range replacement {
				srcLines3 = append(srcLines3, replacementLine)
			}
			replacementIndex++
		} else if strings.TrimSpace(srcLine) == "//nofmt:disable" {
			srcLines3 = append(srcLines3, srcLine)
			enabled = false
		} else {
			if enabled {
				continue
			} else {
				srcLines3 = append(srcLines3, srcLine)
			}
		}
	}
	assert.Nil(scanner.Err(), "scanner error")

	buffer := make([]string, 0)
	enabled = false
	for _, srcLine := range srcLines3 {
		if strings.TrimSpace(srcLine) == "//nofmt:enable" {
			enabled = true
			buffer = append(buffer, srcLine)
		} else if strings.TrimSpace(srcLine) == "//nofmt:disable" {
			enabled = false
			buffer = append(buffer, srcLine)
			if len(buffer) == 3 {
				fmt.Println(buffer[1], "//nofmt")
			} else {
				for _, b := range buffer {
					fmt.Println(b)
				}
			}
			buffer = make([]string, 0)
		} else {
			if enabled {
				buffer = append(buffer, srcLine)
			} else {
				fmt.Println(srcLine)
			}
		}
	}
}

func main() {
	args := Args{}
	clap.Example(strings2.Dedent(`//nofmt:enable
                                 |[custom formatted code]
                                 |//nofmt:disable`))
	clap.Parse(&args)

	var src string
	if args.File != "" {
		src = string(must.Do2(os.ReadFile(args.File)))
	} else {
		src = string(must.Do2(io.ReadAll(os.Stdin)))
	}

	Fmt(src)
}
