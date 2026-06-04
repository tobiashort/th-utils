package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/tobiashort/th-utils/lib/ansi"
	"github.com/tobiashort/th-utils/lib/clap"
	"github.com/tobiashort/th-utils/lib/clog"
	"github.com/tobiashort/th-utils/lib/ellipsis"
	"github.com/tobiashort/th-utils/lib/must"
	"github.com/tobiashort/th-utils/lib/term"
)

type Args struct {
	File string `clap:"positional,desc='The file to open. Reads from Stdin if not specified'"`
}

func main() {
	clog.Level = clog.LevelDebug

	args := Args{}
	clap.Parse(&args)

	var reader io.Reader
	if args.File != "" {
		reader = must.Do2(os.Open(args.File))
	} else {
		reader = os.Stdin
	}

	text := string(must.Do2(io.ReadAll(reader)))
	text = strings.TrimSuffix(text, "\n")
	text = strings.ReplaceAll(text, "\t", "    ")
	textLines := strings.Split(text, "\n")
	maxTextLines := len(textLines)
	maxTextCols := 0
	for _, textLine := range textLines {
		maxTextCols = max(maxTextCols, len([]rune(textLine)))
	}

	defer fmt.Print(ansi.ScreenAlternativeLeave)
	fmt.Print(ansi.ScreenAlternativeEnter)
	tty := must.Do2(term.OpenTTY())
	defer tty.Close()
	must.Do(term.MakeRaw(tty))
	defer term.Restore(tty)

	cols, lines := must.Do3(term.Size(tty))
	startCol := 0
	startLine := 0

draw:
	fmt.Print(ansi.EraseEntireScreen)
	fmt.Print(ansi.CursorMoveToHomePosition)
	fmt.Print(ansi.CursorHide)
	fmt.Print("┌")
	for range cols - 2 {
		fmt.Print("─")
	}
	fmt.Print("┐")
	fmt.Print(ansi.CursorMoveDown(1))
	fmt.Print(ansi.CursorMoveToColumn(0))
	for i := 0; i < min(maxTextLines, lines-3); i++ {
		line := textLines[startLine+i]
		line = fmt.Sprintf("%-*s", cols-2, line)
		line = line[startCol:]
		line = fmt.Sprintf("%-*s", cols-2, line)
		line = ellipsis.EllipsisSuffix(line, cols-2, ">>>")
		fmt.Print("│")
		fmt.Print(line)
		fmt.Print("│")
		fmt.Print(ansi.CursorMoveDown(1))
		fmt.Print(ansi.CursorMoveToColumn(0))
	}
	for i := maxTextLines; i < lines-3; i++ {
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
	fmt.Print(ansi.CursorMoveDown(1))
	fmt.Print(ansi.CursorMoveToColumn(0))
	fmt.Printf(" %d%%, %dl", 100*min(maxTextLines, (startLine+lines-3))/maxTextLines, maxTextLines)

	buf := make([]byte, 1)
eventLoop:
	for {
		must.Do2(tty.Read(buf))
		switch string(buf[0]) {
		case "h":
			startCol--
			startCol = max(startCol, 0)
			goto draw
		case "j":
			if maxTextLines > lines-3 {
				startLine++
				startLine = min(startLine, maxTextLines-lines+3)
				goto draw
			}
		case "k":
			startLine--
			startLine = max(startLine, 0)
			goto draw
		case "l":
			if maxTextCols > cols-2 {
				startCol++
				startCol = min(startCol, maxTextCols-cols+2)
				goto draw
			}
		case ansi.InputCtrlD:
			if maxTextLines > lines-3 {
				startLine += lines / 2
				startLine = min(startLine, maxTextLines-lines+3)
				goto draw
			}
		case ansi.InputCtrlU:
			startLine -= lines / 2
			startLine = max(startLine, 0)
			goto draw
		case "g":
			fmt.Print(ansi.CursorMoveTo(lines, 0))
			fmt.Print(ansi.EraseEntireLine)
			fmt.Print(" g")
			must.Do2(tty.Read(buf))
			switch string(buf[0]) {
			case "e":
				startLine = maxTextLines - lines + 3
			case "g":
				startLine = 0
			}
			goto draw
		case ansi.InputEscape:
			goto draw
		case "q":
			fallthrough
		case ansi.InputCtrlC:
			break eventLoop
		}
	}
}
