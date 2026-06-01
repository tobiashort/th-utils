package main

import (
	"fmt"
	"io"
	"os"

	"github.com/tobiashort/th-utils/lib/ansi"
	"github.com/tobiashort/th-utils/lib/clap"
	"github.com/tobiashort/th-utils/lib/clog"
	"github.com/tobiashort/th-utils/lib/must"
	"github.com/tobiashort/th-utils/lib/term"
)

type Args struct {
	File string `clap:"desc='The file to open. Reads from Stdin if not specified'"`
}

func main() {
	args := Args{}
	clap.Parse(&args)

	var reader io.Reader
	if args.File != "" {
		reader = must.Do2(os.Open(args.File))
	} else {
		reader = os.Stdin
	}

	text := string(must.Do2(io.ReadAll(reader)))
	clog.Infof("read %d runes", len(text))

	defer fmt.Print(ansi.ScreenAlternativeLeave)
	fmt.Print(ansi.ScreenAlternativeEnter)

	defer term.Restore()
	must.Do(term.MakeRaw())

	cols, lines := must.Do3(term.Size())
	fmt.Print(ansi.CursorMoveToHomePosition)
	fmt.Print("┌")
	for range cols - 2 {
		fmt.Print("─")
	}
	fmt.Print("┐")
	fmt.Print(ansi.CursorMoveDown(1))
	fmt.Print(ansi.CursorMoveToColumn(0))
	for range lines - 3 {
		fmt.Print("│")
		for range cols - 2 {
			fmt.Print(" ")
		}
		fmt.Print("│")
		fmt.Print(ansi.CursorMoveDown(1))
		fmt.Print(ansi.CursorMoveToColumn(0))
	}
	fmt.Print("└")
	for range cols - 2 {
		fmt.Print("─")
	}
	fmt.Print("┘")

	buf := make([]byte, 1)
eventLoop:
	for {
		must.Do2(os.Stdin.Read(buf))
		switch string(buf[0]) {
		case ansi.InputCtrlC:
			break eventLoop
		}
	}
}
