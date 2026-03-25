package cfmt

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/tobiashort/th-utils/lib/ansi"
	"github.com/tobiashort/th-utils/lib/must"
	"github.com/tobiashort/th-utils/lib/term"
)

var re = regexp.MustCompile(`(?:^|[^\\])#([A-Za-z]{1,2})\{((?:\\.|[^\\}])*)\}`)

func Print(a ...any) {
	for i := range a {
		a[i] = clr(fmt.Sprint(a[i]), ansi.DecorReset)
	}
	fmt.Print(a...)
}

func Printf(format string, a ...any) {
	fmt.Printf(clr(format, ansi.DecorReset), a...)
}

func Println(a ...any) {
	for i := range a {
		a[i] = clr(fmt.Sprint(a[i]), ansi.DecorReset)
	}
	fmt.Println(a...)
}

func Fprint(w io.Writer, a ...any) {
	for i := range a {
		a[i] = clr(fmt.Sprint(a[i]), ansi.DecorReset)
	}
	fmt.Fprint(w, a...)
}

func Fprintf(w io.Writer, format string, a ...any) {
	fmt.Fprintf(w, clr(format, ansi.DecorReset), a...)
}

func Fprintln(w io.Writer, a ...any) {
	for i := range a {
		a[i] = clr(fmt.Sprint(a[i]), ansi.DecorReset)
	}
	fmt.Fprintln(w, a...)
}

func Sprint(a ...any) string {
	for i := range a {
		a[i] = clr(fmt.Sprint(a[i]), ansi.DecorReset)
	}
	return fmt.Sprint(a...)
}

func Sprintf(format string, a ...any) string {
	return fmt.Sprintf(clr(format, ansi.DecorReset), a...)
}

func Sprintln(a ...any) string {
	for i := range a {
		a[i] = clr(fmt.Sprint(a[i]), ansi.DecorReset)
	}
	return fmt.Sprintln(a...)
}

func stoc(s string) (ansi.Decor, error) {
	//nofmt:enable
	switch s {
	case "r":  return ansi.DecorRed, nil
	case "g":  return ansi.DecorGreen, nil
	case "y":  return ansi.DecorYellow, nil
	case "b":  return ansi.DecorBlue, nil
	case "p":  return ansi.DecorPurple, nil
	case "c":  return ansi.DecorCyan, nil
	case "B":  return ansi.DecorBold, nil
	case "rB": return ansi.DecorRed + ansi.DecorBold, nil
	case "gB": return ansi.DecorGreen + ansi.DecorBold, nil
	case "yB": return ansi.DecorYellow + ansi.DecorBold, nil
	case "bB": return ansi.DecorBlue + ansi.DecorBold, nil
	case "pB": return ansi.DecorPurple + ansi.DecorBold, nil
	case "cB": return ansi.DecorCyan + ansi.DecorBold, nil
	case "U":  return ansi.DecorUnderline, nil
	case "rU": return ansi.DecorRed + ansi.DecorUnderline, nil
	case "gU": return ansi.DecorGreen + ansi.DecorUnderline, nil
	case "yU": return ansi.DecorYellow + ansi.DecorUnderline, nil
	case "bU": return ansi.DecorBlue + ansi.DecorUnderline, nil
	case "pU": return ansi.DecorPurple + ansi.DecorUnderline, nil
	case "cU": return ansi.DecorCyan + ansi.DecorUnderline, nil
	case "R":  return ansi.DecorReversed, nil
	case "rR": return ansi.DecorRed + ansi.DecorReversed, nil
	case "gR": return ansi.DecorGreen + ansi.DecorReversed, nil
	case "yR": return ansi.DecorYellow + ansi.DecorReversed, nil
	case "bR": return ansi.DecorBlue + ansi.DecorReversed, nil
	case "pR": return ansi.DecorPurple + ansi.DecorReversed, nil
	case "cR": return ansi.DecorCyan + ansi.DecorReversed, nil
	default: return "", fmt.Errorf("cannot map string '%s' to ansi Decorcolor", s)
	}
	//nofmt:disable
}

func Cprint(color string, a ...any) {
	Cfprint(os.Stdout, color, a...)
}

func Cfprint(w io.Writer, color string, a ...any) {
	c := must.Do2(stoc(color))
	for i := range a {
		a[i] = clr(fmt.Sprint(a[i]), c)
	}
	if term.IsTerminal() {
		fmt.Fprint(w, c)
	}
	fmt.Fprint(w, a...)
	if term.IsTerminal() {
		fmt.Fprint(w, ansi.DecorReset)
	}
}

func Cprintf(color string, format string, a ...any) {
	Cfprintf(os.Stdout, color, format)
}

func Cfprintf(w io.Writer, color string, format string, a ...any) {
	c := must.Do2(stoc(color))
	if term.IsTerminal() {
		fmt.Fprint(w, c)
	}
	fmt.Fprintf(w, clr(format, c), a...)
	if term.IsTerminal() {
		fmt.Fprint(w, ansi.DecorReset)
	}
}

func Cprintln(color string, a ...any) {
	Cfprintln(os.Stdout, color, a...)
}

func Cfprintln(w io.Writer, color string, a ...any) {
	c := must.Do2(stoc(color))
	for i := range a {
		a[i] = clr(fmt.Sprint(a[i]), c)
	}
	if term.IsTerminal() {
		fmt.Fprint(w, c)
	}
	fmt.Fprintln(w, a...)
	if term.IsTerminal() {
		fmt.Fprint(w, ansi.DecorReset)
	}
}

func clr(str string, reset ansi.Decor) string {
	matches := re.FindAllStringSubmatch(str, -1)
	for _, match := range matches {
		if c, err := stoc(match[1]); err == nil {
			if term.IsTerminal() {
				str = strings.Replace(str, match[0], c+match[2]+reset, 1)
			} else {
				str = strings.Replace(str, match[0], match[1], 1)
			}
		}
	}
	str = strings.ReplaceAll(str, `\}`, `}`)
	str = strings.ReplaceAll(str, `\{`, `{`)
	str = strings.ReplaceAll(str, `\#`, `#`)
	str = strings.ReplaceAll(str, `\\`, `\`)
	return str
}

func Begin(decor ansi.Decor) {
	if term.IsTerminal() {
		fmt.Print(decor)
	}
}

func End() {
	if term.IsTerminal() {
		fmt.Print(ansi.DecorReset)
	}
}
