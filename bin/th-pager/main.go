package main

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"slices"
	"strings"
	"syscall"
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

func TextLines(text string) []string {
	return strings.Split(text, "\n")
}

func TextColumns(text string) int {
	lines := TextLines(text)
	lineLengths := slices2.Map(lines, func(line string) int { return utf8.RuneCountInString(line) })
	return slices.Max(lineLengths)
}

func PrepareText(text string) string {
	tty := must.Do2(term.OpenTTY())
	defer tty.Close()
	ttyCols, _ := must.Do3(term.Size(tty))
	stripped := ansi.Strip(text)
	if TextColumns(stripped) > ttyCols {
		text = stripped
	}
	text = strings.ReplaceAll(text, "\r", "")
	text = strings.TrimSuffix(text, "\n")
	text = strings.ReplaceAll(text, "\t", "    ")
	return text
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
	text = PrepareText(text)

	defer fmt.Print(ansi.ScreenAlternativeLeave)
	fmt.Print(ansi.ScreenAlternativeEnter)
	tty := must.Do2(term.OpenTTY())
	defer tty.Close()
	must.Do(term.MakeRaw(tty))
	defer term.Restore(tty)

	ttyCols, ttyRows := must.Do3(term.Size(tty))
	textLines := TextLines(text)
	maxTextCols := TextColumns(text)
	maxTextLines := len(textLines)
	startCol := 0
	startLine := 0
	lineNumbers := false
	searchTerm := ""
	occurrences := []Occurrence{}
	occurrenceIndex := 0

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGWINCH)

	bufCh := make(chan byte, 1)
	go func() {
		for {
			buf := make([]byte, 1)
			must.Do2(tty.Read(buf))
			bufCh <- buf[0]
		}
	}()

draw:
	fmt.Print(ansi.EraseEntireScreen)
	fmt.Print(ansi.CursorMoveToHomePosition)
	fmt.Print(ansi.CursorHide)
	if args.File != "" {
		cfmt.Printf("#R{ %s }", ellipsis.Ellipsis(args.File, ttyCols))
	} else {
		cfmt.Print("#R{ th-pager }")
	}
	fmt.Print(ansi.CursorMoveDown(1))
	fmt.Print(ansi.CursorMoveToColumn(0))
	for i := 0; i < min(maxTextLines, ttyRows-2); i++ {
		line := textLines[startLine+i]
		line = fmt.Sprintf("%-*s", maxTextCols, line)
		line = line[startCol:]
		if lineNumbers {
			line = cfmt.Sprintf("#R{ %3d } %s", startLine+i+1, line)
		}
		line = strings.TrimRight(line, " ")
		line = ellipsis.EllipsisSuffix(line, ttyCols, cfmt.Sprintf("#R{>>>}"))
		if searchTerm != "" {
			line = strings.ReplaceAll(line, searchTerm, cfmt.Sprintf("#R{%s}", searchTerm))
		}
		fmt.Print(line)
		fmt.Print(ansi.CursorMoveDown(1))
		fmt.Print(ansi.CursorMoveToColumn(0))
	}
	for i := maxTextLines; i < ttyRows-2; i++ {
		fmt.Print(ansi.EraseEntireLine)
		fmt.Print(ansi.CursorMoveDown(1))
		fmt.Print(ansi.CursorMoveToColumn(0))
	}
	fmt.Print(ansi.CursorMoveDown(1))
	fmt.Print(ansi.CursorMoveToColumn(0))
	cfmt.Printf("#R{ %dl, %d%% }", maxTextLines, 100*min(maxTextLines, (startLine+ttyRows-2))/maxTextLines)
eventLoop:
	for {
		select {
		case input := <-bufCh:
			switch string([]byte{input}) {
			case "h":
				startCol--
				startCol = max(startCol, 0)
				goto draw
			case "j":
				if maxTextLines > ttyRows-2 {
					startLine++
					startLine = min(startLine, maxTextLines-ttyRows+2)
					goto draw
				}
			case "k":
				startLine--
				startLine = max(startLine, 0)
				goto draw
			case "l":
				if maxTextCols > ttyCols {
					startCol++
					startCol = min(startCol, maxTextCols-ttyCols)
					goto draw
				}
			case ansi.InputCtrlD:
				if maxTextLines > ttyRows-2 {
					startLine += ttyRows / 2
					startLine = min(startLine, maxTextLines-ttyRows+2)
					goto draw
				}
			case ansi.InputCtrlU:
				startLine -= ttyRows / 2
				startLine = max(startLine, 0)
				goto draw
			case "g":
				fmt.Print(ansi.CursorMoveTo(ttyRows, 0))
				fmt.Print(ansi.EraseEntireLine)
				cfmt.Print("#R{ g }")
				input = <-bufCh
				switch string([]byte{input}) {
				case "e":
					if maxTextLines > ttyRows-2 {
						startLine = maxTextLines - ttyRows + 2
					}
				case "l":
					if maxTextCols > ttyCols {
						startCol = maxTextCols - ttyCols
					}
				case "h":
					startCol = 0
				case "g":
					startLine = 0
				}
				goto draw
			case "N":
				lineNumbers = !lineNumbers
				if lineNumbers {
					maxTextCols = TextColumns(text)
					maxTextCols += utf8.RuneCountInString(fmt.Sprintf(" %3d  ", maxTextLines))
				} else {
					maxTextCols = TextColumns(text)
				}
				startCol = 0
				goto draw
			case "n":
				if len(occurrences) > 0 {
					if occurrenceIndex+1 < len(occurrences) {
						occurrenceIndex++
					} else {
						occurrenceIndex = 0
					}
					startLine = min(maxTextLines-ttyRows+2, occurrences[occurrenceIndex].Line)
					startCol = max(0, occurrences[occurrenceIndex].Col+10-ttyCols)
					startCol = min(startCol, maxTextCols-ttyCols)
				}
				goto draw
			case "p":
				if len(occurrences) > 0 {
					if occurrenceIndex-1 >= 0 {
						occurrenceIndex--
					} else {
						occurrenceIndex = len(occurrences) - 1
					}
					startLine = min(maxTextLines-ttyRows+2, occurrences[occurrenceIndex].Line)
					startCol = max(0, occurrences[occurrenceIndex].Col+10-ttyCols)
					startCol = min(startCol, maxTextCols-ttyCols)
				}
				goto draw
			case "/":
				fmt.Print(ansi.CursorMoveTo(ttyRows, 0))
				fmt.Print(ansi.EraseEntireLine)
				cfmt.Print("#R{ / }")
				searchTermNew := ""
				for {
					input = <-bufCh
					switch string([]byte{input}) {
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
							fmt.Print(ansi.CursorMoveTo(ttyRows, 0))
							fmt.Print(ansi.EraseEntireLine)
							cfmt.Print("#R{ not found }")
						} else {
							startLine = min(maxTextLines-ttyRows+2, occurrences[occurrenceIndex].Line)
							startCol = max(0, occurrences[occurrenceIndex].Col+10-ttyCols)
							startCol = min(startCol, maxTextCols-ttyCols)
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
						searchTermNew += string([]byte{input})
					}
					fmt.Print(ansi.CursorMoveTo(ttyRows, 0))
					fmt.Print(ansi.EraseEntireLine)
					cfmt.Printf("#R{ /%s }", searchTermNew)
				}
			case "q":
				fallthrough
			case ansi.InputCtrlC:
				break eventLoop
			}
		case _ = <-signalCh:
			ttyCols, ttyRows = must.Do3(term.Size(tty))
			goto draw
		}
	}
}
