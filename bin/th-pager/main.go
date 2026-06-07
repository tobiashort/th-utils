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

func MaxTokenCols(tokens [][]Token) int {
	return slices.Max(slices2.Map(tokens, func(tokenRow []Token) int { return len(tokenRow) }))
}

func Line(tokenRow []Token) string {
	b := strings.Builder{}
	for _, token := range tokenRow {
		b.WriteString(token.Literal)
	}
	return b.String()
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
	tokens := Parse(text)
	maxTokenRows := len(tokens)
	maxTokenCols := MaxTokenCols(tokens)

	defer fmt.Print(ansi.ScreenAlternativeLeave)
	fmt.Print(ansi.ScreenAlternativeEnter)
	tty := must.Do2(term.OpenTTY())
	defer tty.Close()
	must.Do(term.MakeRaw(tty))
	defer term.Restore(tty)

	ttyDim := must.Do2(term.Size(tty))
	ttyCols := ttyDim.Cols
	ttyRows := ttyDim.Rows
	startCol := 0
	startLine := 0
	lineNumbers := false
	searchTerm := ""
	occurrences := []Occurrence{}
	occurrenceIndex := 0

	dimCh := term.OnResize(tty)

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
	for i := 0; i < min(maxTokenRows, ttyRows-2); i++ {
		tokenRow := tokens[startLine+i]

		j := 0
		col := 0
		for ; col < len(tokenRow) && j < startCol; col++ {
			if tokenRow[col].Type == TokenAnsi {
				continue
			}
			j++
		}

		tokenRow = append(slices2.Filter(tokenRow[:col], func(t Token) bool { return t.Type == TokenAnsi }), tokenRow[col:]...)

		line := Line(tokenRow)

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
	for i := maxTokenRows; i < ttyRows-2; i++ {
		fmt.Print(ansi.EraseEntireLine)
		fmt.Print(ansi.CursorMoveDown(1))
		fmt.Print(ansi.CursorMoveToColumn(0))
	}
	fmt.Print(ansi.CursorMoveDown(1))
	fmt.Print(ansi.CursorMoveToColumn(0))
	cfmt.Printf("#R{ %dl, %d%% }", maxTokenRows, 100*min(maxTokenRows, (startLine+ttyRows-2))/maxTokenRows)

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
				if maxTokenRows > ttyRows-2 {
					startLine++
					startLine = min(startLine, maxTokenRows-ttyRows+2)
					goto draw
				}
			case "k":
				startLine--
				startLine = max(startLine, 0)
				goto draw
			case "l":
				if maxTokenCols > ttyCols {
					startCol++
					startCol = min(startCol, maxTokenCols-ttyCols)
					goto draw
				}
			case ansi.InputCtrlD:
				if maxTokenRows > ttyRows-2 {
					startLine += ttyRows / 2
					startLine = min(startLine, maxTokenRows-ttyRows+2)
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
					if maxTokenRows > ttyRows-2 {
						startLine = maxTokenRows - ttyRows + 2
					}
				case "l":
					if maxTokenCols > ttyCols {
						startCol = maxTokenCols - ttyCols
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
					maxTokenCols = MaxTokenCols(tokens)
					maxTokenCols += utf8.RuneCountInString(fmt.Sprintf(" %3d  ", maxTokenCols))
				} else {
					maxTokenCols = MaxTokenCols(tokens)
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
					startLine = max(0, min(maxTokenRows-ttyRows+2, occurrences[occurrenceIndex].Line))
					startCol = max(0, occurrences[occurrenceIndex].Col+10-ttyCols)
					startCol = max(0, min(startCol, maxTokenCols-ttyCols))
				}
				goto draw
			case "p":
				if len(occurrences) > 0 {
					if occurrenceIndex-1 >= 0 {
						occurrenceIndex--
					} else {
						occurrenceIndex = len(occurrences) - 1
					}
					startLine = max(0, min(maxTokenRows-ttyRows+2, occurrences[occurrenceIndex].Line))
					startCol = max(0, occurrences[occurrenceIndex].Col+10-ttyCols)
					startCol = max(0, min(startCol, maxTokenCols-ttyCols))
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
						for i := startLine; i < len(tokens); i++ {
							tokenRow := tokens[i]
							line := Line(tokenRow)
							for _, index := range strings2.AllIndexes(line, searchTerm) {
								occurrences = append(occurrences, Occurrence{Line: i, Col: index})
							}
						}
						for i := 0; i < startLine; i++ {
							tokenRow := tokens[i]
							line := Line(tokenRow)
							for _, index := range strings2.AllIndexes(line, searchTerm) {
								occurrences = append(occurrences, Occurrence{Line: i, Col: index})
							}
						}
						if len(occurrences) == 0 {
							fmt.Print(ansi.CursorMoveTo(ttyRows, 0))
							fmt.Print(ansi.EraseEntireLine)
							cfmt.Print("#R{ not found }")
						} else {
							startLine = max(0, min(maxTokenRows-ttyRows+2, occurrences[occurrenceIndex].Line))
							startCol = max(0, occurrences[occurrenceIndex].Col+10-ttyCols)
							startCol = max(0, min(startCol, maxTokenCols-ttyCols))
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
		case ttyDim = <-dimCh:
			ttyCols = ttyDim.Cols
			ttyRows = ttyDim.Rows
			goto draw
		}
	}
}
