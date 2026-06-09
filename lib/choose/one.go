package choose

import (
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/tobiashort/th-utils/lib/ansi"
	"github.com/tobiashort/th-utils/lib/cfmt"
	"github.com/tobiashort/th-utils/lib/ellipsis"
	"github.com/tobiashort/th-utils/lib/must"
	"github.com/tobiashort/th-utils/lib/term"
)

type Chooser struct {
	Writer    io.Writer
	Formatter cfmt.Formatter
	SortFunc  func(o1, o2 Option, search string) int
}

var DefaultChooser = Chooser{
	Writer:    os.Stderr,
	Formatter: cfmt.DefaultFormatter,
	SortFunc:  nil,
}

type Option struct {
	Index int
	Value string
}

func ToOptions[T any](s []T) []Option {
	return ToOptionsFunc(s, func(v T) string {
		return fmt.Sprintf("%v", v)
	})
}

func ToOptionsFunc[T any](s []T, f func(v T) string) []Option {
	r := make([]Option, len(s))
	for i, v := range s {
		r[i] = Option{
			Index: i,
			Value: f(v),
		}
	}
	return r
}

func (c Chooser) One(prompt string, options []Option) (Option, bool) {
	tty := must.Do2(term.OpenTTY())
	must.Do(term.MakeRaw(tty))
	defer func() {
		must.Do(term.Restore(tty))
		tty.Close()
	}()

	dim := must.Do2(term.Size(tty))

	ok := false
	selectedIndex := 0
	maxLines := 5
	selectedLine := 0
	search := strings.Builder{}

draw:
	filtered := make([]Option, 0)
	if search.String() == "" {
		for _, option := range options {
			filtered = append(filtered, option)
		}
	} else {
		for _, option := range options {
			lSearch := strings.ToLower(search.String())
			lOption := strings.ToLower(option.Value)
			if strings.Contains(lOption, lSearch) {
				filtered = append(filtered, option)
			}
		}
	}

	if c.SortFunc != nil {
		slices.SortStableFunc(filtered, func(o1, o2 Option) int {
			return c.SortFunc(o1, o2, search.String())
		})
	}

	fmt.Fprintf(c.Writer, "%s\r\n", prompt)
	if len(filtered) > 0 {
		for index := selectedIndex - selectedLine; index < min(selectedIndex+(maxLines-selectedLine), len(filtered)); index++ {
			option := filtered[index]
			if index == selectedIndex {
				c.Formatter.Fprintf(c.Writer, "#yB{▌ %s}\r\n", ellipsis.Ellipsis(option.Value, dim.Cols-2, "...", ellipsis.PosCenter))
			} else {
				fmt.Fprintf(c.Writer, "  %s\r\n", ellipsis.Ellipsis(option.Value, dim.Cols-2, "...", ellipsis.PosCenter))
			}
		}
	}
	c.Formatter.Fprintf(c.Writer, "  #b{%d/%d}\r\n", min(selectedIndex+1, len(filtered)), len(filtered))
	c.Formatter.Fprintf(c.Writer, "#bB{>} %s", search.String())

	buf := make([]byte, 3)
	for {
		n := must.Do2(os.Stdin.Read(buf))
		switch string(buf[:n]) {
		case ansi.InputTab:
			fallthrough
		case ansi.InputKeyDown:
			if selectedLine < maxLines-1 {
				selectedLine++
			}
			if selectedIndex < len(filtered)-1 {
				selectedIndex++
			} else {
				selectedLine = 0
				selectedIndex = 0
			}
			goto redraw
		case ansi.InputShiftTab:
			fallthrough
		case ansi.InputKeyUp:
			if selectedLine > 0 {
				selectedLine--
			}
			if selectedIndex > 0 {
				selectedIndex--
			} else {
				selectedLine = maxLines
				selectedIndex = len(filtered) - 1
			}
			goto redraw
		case ansi.InputCR:
			fallthrough
		case ansi.InputLF:
			fallthrough
		case ansi.InputCRLF:
			ok = true
			goto done
		case ansi.InputEscape:
			ok = false
			goto done
		case ansi.InputCtrlC:
			ok = false
			goto done
		case ansi.InputBackSpace:
			fallthrough
		case ansi.InputDelete:
			s := search.String()
			if s != "" {
				s = s[:len(s)-1]
				search.Reset()
				search.WriteString(s)
				selectedIndex = 0
				selectedLine = 0
			}
			goto redraw
		default:
			search.Write(buf[:n])
			selectedIndex = 0
			selectedLine = 0
			goto redraw
		}
	}

redraw:
	fmt.Fprint(c.Writer, "\r")
	for range min(maxLines, len(filtered)) + 2 {
		fmt.Fprint(c.Writer, ansi.EraseEntireLine)
		fmt.Fprint(c.Writer, ansi.CursorMoveUp(1))
	}
	goto draw

done:
	fmt.Fprint(c.Writer, "\r")
	for range min(maxLines, len(filtered)) + 2 {
		fmt.Fprint(c.Writer, ansi.EraseEntireLine)
		fmt.Fprint(c.Writer, ansi.CursorMoveUp(1))
	}
	fmt.Fprint(c.Writer, ansi.EraseEntireLine)
	if selectedIndex >= 0 && selectedIndex < len(filtered) {
		return filtered[selectedIndex], ok
	}
	return Option{Index: -1, Value: ""}, false
}

func One(prompt string, options []Option) (Option, bool) {
	return DefaultChooser.One(prompt, options)
}
