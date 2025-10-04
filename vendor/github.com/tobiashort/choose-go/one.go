package choose

import (
	"fmt"
	"os"
	"strings"

	"github.com/tobiashort/ansi-go"
	"github.com/tobiashort/cfmt-go"
	"github.com/tobiashort/utils-go/must"

	"golang.org/x/term"
)

func ToOptions[T any](s []T) []string {
	return ToOptionsFunc(s, func(v T) string {
		return fmt.Sprintf("%v", v)
	})
}

func ToOptionsFunc[T any](s []T, f func(v T) string) []string {
	r := make([]string, len(s))
	for i, v := range s {
		r[i] = f(v)
	}
	return r
}

func One(prompt string, options []string) (int, string, bool) {
	fd := int(os.Stdin.Fd())
	oldState := must.Do2(term.MakeRaw(fd))
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	ok := false
	selectedIndex := 0
	maxLines := 5
	selectedLine := 0
	search := strings.Builder{}

	type option struct {
		origIndex int
		value     string
	}

draw:
	filtered := make([]option, 0)
	if search.String() == "" {
		for i, s := range options {
			filtered = append(filtered, option{origIndex: i, value: s})
		}
	} else {
		for i, s := range options {
			lSearch := strings.ToLower(search.String())
			lOption := strings.ToLower(s)
			if strings.Contains(lOption, lSearch) {
				filtered = append(filtered, option{origIndex: i, value: s})
			}
		}
	}

	fmt.Printf("%s\r\n", prompt)
	if len(filtered) > 0 {
		for index := selectedIndex - selectedLine; index < min(selectedIndex+(maxLines-selectedLine), len(filtered)); index++ {
			option := filtered[index]
			if index == selectedIndex {
				cfmt.Printf("#yB{▌ %s}\r\n", option.value)
			} else {
				fmt.Printf("  %s\r\n", option.value)
			}
		}
	}
	cfmt.Printf("  #b{%d/%d}\r\n", min(selectedIndex+1, len(filtered)), len(filtered))
	cfmt.Printf("#bB{>} %s", search.String())

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
	fmt.Print("\r")
	for range min(maxLines, len(filtered)) + 2 {
		fmt.Print(ansi.EraseEntireLine)
		fmt.Print(ansi.CursorMoveUp(1))
	}
	goto draw

done:
	fmt.Print("\r")
	for range min(maxLines, len(filtered)) + 2 {
		fmt.Print(ansi.EraseEntireLine)
		fmt.Print(ansi.CursorMoveUp(1))
	}
	if selectedIndex >= 0 && selectedIndex < len(filtered) {
		return filtered[selectedIndex].origIndex, filtered[selectedIndex].value, ok
	}
	return -1, "", false
}
