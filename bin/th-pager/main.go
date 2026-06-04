package main

import (
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/tobiashort/th-utils/lib/ansi"
	"github.com/tobiashort/th-utils/lib/cfmt"
	"github.com/tobiashort/th-utils/lib/clap"
	"github.com/tobiashort/th-utils/lib/clog"
	"github.com/tobiashort/th-utils/lib/ellipsis"
	"github.com/tobiashort/th-utils/lib/must"
	slices2 "github.com/tobiashort/th-utils/lib/slices"
	strings2 "github.com/tobiashort/th-utils/lib/strings"
	"github.com/tobiashort/th-utils/lib/term"
)

type Args struct {
	File string `clap:"positional,desc='The file to open. Reads from Stdin if not specified'"`
}

type Occurrence struct {
	Line int
	Col  int
}

var (
	textLines []string
)

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
	text = ansi.Strip(text)
	text = strings.ReplaceAll(text, "\r", "")
	text = strings.TrimSuffix(text, "\n")
	text = strings.ReplaceAll(text, "\t", "    ")
	textLines = strings.Split(text, "\n")
	maxTextLines := len(textLines)
	maxTextCols := slices.Max(slices2.Map(textLines, func(line string) int { return utf8.RuneCountInString(line) }))

	defer fmt.Print(ansi.ScreenAlternativeLeave)
	fmt.Print(ansi.ScreenAlternativeEnter)
	tty := must.Do2(term.OpenTTY())
	defer tty.Close()
	must.Do(term.MakeRaw(tty))
	defer term.Restore(tty)

	cols, lines := must.Do3(term.Size(tty))
	startCol := 0
	startLine := 0

	lineNumbers := false

	searchTerm := ""
	occurrences := []Occurrence{}
	occurrenceIndex := 0

draw:
	fmt.Print(ansi.EraseEntireScreen)
	fmt.Print(ansi.CursorMoveToHomePosition)
	fmt.Print(ansi.CursorHide)
	if args.File != "" {
		cfmt.Printf("#R{ %s }", ellipsis.Ellipsis(args.File, cols))
	} else {
		cfmt.Print("#R{ th-pager }")
	}
	fmt.Print(ansi.CursorMoveDown(1))
	fmt.Print(ansi.CursorMoveToColumn(0))
	for i := 0; i < min(maxTextLines, lines-2); i++ {
		line := textLines[startLine+i]
		line = fmt.Sprintf("%-*s", maxTextCols, line)
		line = line[startCol:]
		if lineNumbers {
			line = cfmt.Sprintf("#R{ %3d } %s", startLine+i+1, line)
		}
		line = strings.TrimRight(line, " ")
		line = ellipsis.EllipsisSuffix(line, cols, cfmt.Sprintf("#R{>>>}"))
		if searchTerm != "" {
			line = strings.ReplaceAll(line, searchTerm, cfmt.Sprintf("#R{%s}", searchTerm))
		}
		fmt.Print(line)
		fmt.Print(ansi.CursorMoveDown(1))
		fmt.Print(ansi.CursorMoveToColumn(0))
	}
	for i := maxTextLines; i < lines-2; i++ {
		fmt.Print(ansi.EraseEntireLine)
		fmt.Print(ansi.CursorMoveDown(1))
		fmt.Print(ansi.CursorMoveToColumn(0))
	}
	fmt.Print(ansi.CursorMoveDown(1))
	fmt.Print(ansi.CursorMoveToColumn(0))
	cfmt.Printf("#R{ %dl, %d%% }", maxTextLines, 100*min(maxTextLines, (startLine+lines-2))/maxTextLines)

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
			if maxTextLines > lines-2 {
				startLine++
				startLine = min(startLine, maxTextLines-lines+2)
				goto draw
			}
		case "k":
			startLine--
			startLine = max(startLine, 0)
			goto draw
		case "l":
			if maxTextCols > cols {
				startCol++
				startCol = min(startCol, maxTextCols-cols)
				goto draw
			}
		case ansi.InputCtrlD:
			if maxTextLines > lines-2 {
				startLine += lines / 2
				startLine = min(startLine, maxTextLines-lines+2)
				goto draw
			}
		case ansi.InputCtrlU:
			startLine -= lines / 2
			startLine = max(startLine, 0)
			goto draw
		case "g":
			fmt.Print(ansi.CursorMoveTo(lines, 0))
			fmt.Print(ansi.EraseEntireLine)
			cfmt.Print("#R{ g }")
			must.Do2(tty.Read(buf))
			switch string(buf[0]) {
			case "e":
				if maxTextLines > lines-2 {
					startLine = maxTextLines - lines + 2
				}
			case "l":
				if maxTextCols > cols {
					startCol = maxTextCols - cols
				}
			case "h":
				startCol = 0
			case "g":
				startLine = 0
			}
			goto draw
		case "N":
			lineNumbers = !lineNumbers
			goto draw
		case "n":
			if occurrenceIndex+1 < len(occurrences) {
				occurrenceIndex++
			} else {
				occurrenceIndex = 0
			}
			startLine = min(maxTextLines-lines+2, occurrences[occurrenceIndex].Line)
			startCol = max(0, occurrences[occurrenceIndex].Col+10-cols)
			startCol = min(startCol, maxTextCols-cols)
			goto draw
		case "p":
			if occurrenceIndex-1 >= 0 {
				occurrenceIndex--
			} else {
				occurrenceIndex = len(occurrences) - 1
			}
			startLine = min(maxTextLines-lines+2, occurrences[occurrenceIndex].Line)
			startCol = max(0, occurrences[occurrenceIndex].Col+10-cols)
			startCol = min(startCol, maxTextCols-cols)
			goto draw
		case "/":
			fmt.Print(ansi.CursorMoveTo(lines, 0))
			fmt.Print(ansi.EraseEntireLine)
			cfmt.Print("#R{ / }")
			searchTermNew := ""
			for {
				must.Do2(tty.Read(buf))
				switch string(buf[0]) {
				case ansi.InputCR:
					fallthrough
				case ansi.InputLF:
					searchTerm = searchTermNew
					occurrences = []Occurrence{}
					occurrenceIndex = 0
					for i := startLine; i < len(textLines); i++ {
						line := textLines[i]
						for _, index := range strings2.AllIndexes(line, searchTerm) {
							occurrences = append(occurrences, Occurrence{Line: i, Col: index})
						}
					}
					for i := 0; i < startLine; i++ {
						line := textLines[i]
						for _, index := range strings2.AllIndexes(line, searchTerm) {
							occurrences = append(occurrences, Occurrence{Line: i, Col: index})
						}
					}
					if len(occurrences) == 0 {
						fmt.Print(ansi.CursorMoveTo(lines, 0))
						fmt.Print(ansi.EraseEntireLine)
						cfmt.Print("#R{ not found }")
					} else {
						startLine = min(maxTextLines-lines+2, occurrences[occurrenceIndex].Line)
						startCol = max(0, occurrences[occurrenceIndex].Col+10-cols)
						startCol = min(startCol, maxTextCols-cols)
						goto draw
					}
				case ansi.InputDelete:
					fallthrough
				case ansi.InputBackSpace:
					if searchTermNew != "" {
						searchTermNew = searchTermNew[:len(searchTermNew)-1]
					}
				case ansi.InputEscape:
					goto draw
				default:
					searchTermNew += string(buf[0])
				}
				fmt.Print(ansi.CursorMoveTo(lines, 0))
				fmt.Print(ansi.EraseEntireLine)
				cfmt.Printf("#R{ /%s }", searchTermNew)
			}
		case "q":
			fallthrough
		case ansi.InputCtrlC:
			break eventLoop
		}
	}
}
