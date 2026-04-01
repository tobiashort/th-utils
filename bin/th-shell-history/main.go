package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/tobiashort/th-utils/lib/ansi"
	"github.com/tobiashort/th-utils/lib/cfmt"
	strings2 "github.com/tobiashort/th-utils/lib/strings"

	"github.com/tobiashort/th-utils/lib/assert"
	"github.com/tobiashort/th-utils/lib/choose"
	"github.com/tobiashort/th-utils/lib/clap"
	"github.com/tobiashort/th-utils/lib/orderedmap"
)

type Args struct {
	Fish        bool `clap:"desc='Integration for fish shell'"`
	Integration bool `clap:"desc='Print integration code'"`
}

func run() int {
	args := Args{}

	clap.Example(`Fish:
	th-shell-history --integration --fish | source`)

	clap.Parse(&args)

	history := orderedmap.NewOrderedMap[string, struct{}]()

	if args.Integration {
		if args.Fish {
			fmt.Print(strings2.Dedent(`function shell-history
									  |  set selected (th-shell-history --fish)
    								  |  if test -n "$selected"
        							  |    commandline --replace $selected
    								  |  end
          							  |end
          							  |
									  |bind \cr shell-history
									  |bind -M insert \cr shell-history
									  |`))
			return 0
		} else {
			clap.PrintHelp(&args, os.Stderr)
			return 1
		}
	} else if args.Fish {
		historyFilePath := os.ExpandEnv("${HOME}/.local/share/fish/fish_history")
		historyFile, err := os.Open(historyFilePath)
		if err != nil {
			if os.IsNotExist(err) {
				return 1
			}
			panic(err)
		}
		scanner := bufio.NewScanner(historyFile)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "- cmd: ") {
				cmd := line[7:]
				history.Del(cmd)
				history.Put(cmd, struct{}{})
			}
		}
		assert.Nil(scanner.Err(), "scanner error")
	} else {
		clap.PrintHelp(&args, os.Stderr)
		return 1
	}

	cmds := history.Keys()
	slices.Reverse(cmds)
	chooser := choose.Chooser{
		Writer:    os.Stderr,
		Formatter: cfmt.Formatter{ForceColors: true},
		SortFunc:  nil,
	}
	_, col := ansi.CursorGetCurrentPosition()
	fmt.Fprintln(chooser.Writer)
	option, ok := chooser.One("Search history:", choose.ToOptions(cmds))
	fmt.Fprint(chooser.Writer, ansi.EraseEntireLine)
	fmt.Fprint(chooser.Writer, ansi.CursorMoveUp(1))
	fmt.Fprint(chooser.Writer, ansi.CursorMoveToColumn(col))
	if ok {
		fmt.Fprint(os.Stdout, option.Value)
		return 0
	}

	return 1
}

func main() {
	os.Exit(run())
}
