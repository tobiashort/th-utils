//go:build ignore

package main

import (
	"fmt"
	"os"
	"os/exec"
	"unicode"

	"github.com/tobiashort/th-utils/lib/ansi"
	"github.com/tobiashort/th-utils/lib/must"
	"github.com/tobiashort/th-utils/lib/term"
)

func main() {
	tty := must.Do2(term.OpenTTY())
	must.Do(term.MakeRaw(tty))
	defer term.Restore(tty)
	src := must.Do2(os.Create("width.go"))
	fmt.Fprintln(src, "package unicode")
	fmt.Fprintln(src, "var widths = map[rune]int{")
	for r := rune(0); r <= unicode.MaxRune; r++ {
		fmt.Print(ansi.EraseEntireLine)
		fmt.Print(ansi.CursorMoveToColumn(1))
		fmt.Printf("%c", r)
		_, col := ansi.CursorGetCurrentPosition()
		width := col - 1
		if col > 1 {
			fmt.Fprintf(src, "rune(0x%x): %d,\n", r, width)
		}
	}
	fmt.Fprintln(src, "}")
	fmt.Fprint(src, `func Width(r rune) int {
	if w, ok := widths[r]; ok {
		return w
	} else {
		return 1
	}
}

func WidthString(s string) int {
	sum := 0
	for _, r := range []rune(s) {
		sum += Width(r)
	}
	return sum
}
`)
	must.Do(exec.Command("gofmt", "-w", "width.go").Run())
}
