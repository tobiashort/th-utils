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

type Formatter struct {
	ForceColors bool
}

var DefaultFormatter = Formatter{
	ForceColors: false,
}

var re = regexp.MustCompile(`#([A-Za-z]{1,2})\{((?:\\.|[^\\}])*)\}`)

func (f Formatter) inColor() bool {
	return term.IsTerminal() || f.ForceColors
}

func (f Formatter) Print(a ...any) {
	for i := range a {
		a[i] = f.clr(fmt.Sprint(a[i]), ansi.DecorReset)
	}
	fmt.Print(a...)
}

func (f Formatter) Printf(format string, a ...any) {
	fmt.Printf(f.clr(format, ansi.DecorReset), a...)
}

func (f Formatter) Println(a ...any) {
	for i := range a {
		a[i] = f.clr(fmt.Sprint(a[i]), ansi.DecorReset)
	}
	fmt.Println(a...)
}

func (f Formatter) Fprint(w io.Writer, a ...any) {
	for i := range a {
		a[i] = f.clr(fmt.Sprint(a[i]), ansi.DecorReset)
	}
	fmt.Fprint(w, a...)
}

func (f Formatter) Fprintf(w io.Writer, format string, a ...any) {
	fmt.Fprintf(w, f.clr(format, ansi.DecorReset), a...)
}

func (f Formatter) Fprintln(w io.Writer, a ...any) {
	for i := range a {
		a[i] = f.clr(fmt.Sprint(a[i]), ansi.DecorReset)
	}
	fmt.Fprintln(w, a...)
}

func (f Formatter) Sprint(a ...any) string {
	for i := range a {
		a[i] = f.clr(fmt.Sprint(a[i]), ansi.DecorReset)
	}
	return fmt.Sprint(a...)
}

func (f Formatter) Sprintf(format string, a ...any) string {
	return fmt.Sprintf(f.clr(format, ansi.DecorReset), a...)
}

func (f Formatter) Sprintln(a ...any) string {
	for i := range a {
		a[i] = f.clr(fmt.Sprint(a[i]), ansi.DecorReset)
	}
	return fmt.Sprintln(a...)
}

func (f Formatter) stoc(s string) (ansi.Decor, error) {
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

func (f Formatter) Cprint(color string, a ...any) {
	Cfprint(os.Stdout, color, a...)
}

func (f Formatter) Cfprint(w io.Writer, color string, a ...any) {
	c := must.Do2(f.stoc(color))
	for i := range a {
		a[i] = f.clr(fmt.Sprint(a[i]), c)
	}
	if f.inColor() {
		fmt.Fprint(w, c)
	}
	fmt.Fprint(w, a...)
	if f.inColor() {
		fmt.Fprint(w, ansi.DecorReset)
	}
}

func (f Formatter) Cprintf(color string, format string, a ...any) {
	Cfprintf(os.Stdout, color, format)
}

func (f Formatter) Cfprintf(w io.Writer, color string, format string, a ...any) {
	c := must.Do2(f.stoc(color))
	if f.inColor() {
		fmt.Fprint(w, c)
	}
	fmt.Fprintf(w, f.clr(format, c), a...)
	if f.inColor() {
		fmt.Fprint(w, ansi.DecorReset)
	}
}

func (f Formatter) Cprintln(color string, a ...any) {
	Cfprintln(os.Stdout, color, a...)
}

func (f Formatter) Cfprintln(w io.Writer, color string, a ...any) {
	c := must.Do2(f.stoc(color))
	for i := range a {
		a[i] = f.clr(fmt.Sprint(a[i]), c)
	}
	if f.inColor() {
		fmt.Fprint(w, c)
	}
	fmt.Fprintln(w, a...)
	if f.inColor() {
		fmt.Fprint(w, ansi.DecorReset)
	}
}

func (f Formatter) clr(str string, reset ansi.Decor) string {
	matches := re.FindAllStringSubmatch(str, -1)
	for _, match := range matches {
		if c, err := f.stoc(match[1]); err == nil {
			replace := match[0]
			replaceWith := match[2]
			replaceWith = strings.ReplaceAll(replaceWith, `\}`, `}`)
			replaceWith = strings.ReplaceAll(replaceWith, `\\`, `\`)
			if f.inColor() {
				str = strings.Replace(str, replace, c+replaceWith+reset, 1)
			} else {
				str = strings.Replace(str, replace, replaceWith, 1)
			}
		}
	}
	return str
}

func (f Formatter) Begin(decor ansi.Decor) {
	if f.inColor() {
		fmt.Print(decor)
	}
}

func (f Formatter) End() {
	if f.inColor() {
		fmt.Print(ansi.DecorReset)
	}
}

func Print(a ...any) {
	DefaultFormatter.Print(a...)
}

func Printf(format string, a ...any) {
	DefaultFormatter.Printf(format, a...)
}

func Println(a ...any) {
	DefaultFormatter.Println(a...)
}

func Fprint(w io.Writer, a ...any) {
	DefaultFormatter.Fprint(w, a...)
}

func Fprintf(w io.Writer, format string, a ...any) {
	DefaultFormatter.Fprintf(w, format, a...)
}

func Fprintln(w io.Writer, a ...any) {
	DefaultFormatter.Fprintln(w, a...)
}

func Sprint(a ...any) string {
	return DefaultFormatter.Sprint(a...)
}

func Sprintf(format string, a ...any) string {
	return DefaultFormatter.Sprintf(format, a...)
}

func Sprintln(a ...any) string {
	return DefaultFormatter.Sprintln(a...)
}

func Cprint(color string, a ...any) {
	DefaultFormatter.Cprint(color, a...)
}

func Cfprint(w io.Writer, color string, a ...any) {
	DefaultFormatter.Cfprint(w, color, a...)
}

func Cprintf(color string, format string, a ...any) {
	DefaultFormatter.Cprintf(color, format, a...)
}

func Cfprintf(w io.Writer, color string, format string, a ...any) {
	DefaultFormatter.Cfprintf(w, color, format, a...)
}

func Cprintln(color string, a ...any) {
	DefaultFormatter.Cprintln(color, a...)
}

func Cfprintln(w io.Writer, color string, a ...any) {
	DefaultFormatter.Cfprintln(w, color, a...)
}

func Begin(decor ansi.Decor) {
	DefaultFormatter.Begin(decor)
}

func End() {
	DefaultFormatter.End()
}
