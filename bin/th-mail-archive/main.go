package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"github.com/tobiashort/th-utils/lib/ansi"
	"github.com/tobiashort/th-utils/lib/cfmt"
	"github.com/tobiashort/th-utils/lib/clap"
	"github.com/tobiashort/th-utils/lib/clog"
	"github.com/tobiashort/th-utils/lib/ellipsis"
	"github.com/tobiashort/th-utils/lib/must"
	"github.com/tobiashort/th-utils/lib/term"
)

type Args struct {
}

func NotMuchNew() chan string {
	out := make(chan string)
	cmd := exec.Command("notmuch", "new")
	stderr := must.Do2(cmd.StderrPipe())
	stdout := must.Do2(cmd.StdoutPipe())
	errDone, outDone := false, false
	go func() {
		cmd.Run()
		for {
			if errDone && outDone {
				close(out)
				break
			}
		}
	}()
	go func() {
		s := bufio.NewScanner(stderr)
		for s.Scan() {
			out <- s.Text()
		}
		errDone = true
	}()
	go func() {
		s := bufio.NewScanner(stdout)
		for s.Scan() {
			out <- s.Text()
		}
		outDone = true
	}()
	return out
}

func NotMuchSearch(search string) chan string {
	out := make(chan string)
	cmd := exec.Command("notmuch", "search", search)
	stderr := must.Do2(cmd.StderrPipe())
	stdout := must.Do2(cmd.StdoutPipe())
	errDone, outDone := false, false
	go func() {
		cmd.Run()
		for {
			if errDone && outDone {
				close(out)
				break
			}
		}
	}()
	go func() {
		s := bufio.NewScanner(stderr)
		for s.Scan() {
			out <- s.Text()
		}
		errDone = true
	}()
	go func() {
		s := bufio.NewScanner(stdout)
		for s.Scan() {
			out <- s.Text()
		}
		outDone = true
	}()
	return out
}

func run() int {
	args := Args{}
	clap.Parse(&args)

	_, err := exec.LookPath("notmuch")
	if err != nil {
		clog.Error("not found: notmuch")
		return 1
	}

	tty := must.Do2(term.OpenTTY())
	defer tty.Close()
	must.Do(term.MakeRaw(tty))
	defer term.Restore(tty)

	onKeystroke := make(chan byte, 1)
	go func() {
		for {
			buf := make([]byte, 1)
			must.Do2(tty.Read(buf))
			onKeystroke <- buf[0]
		}
	}()

	onResize := term.OnResize(tty)

	fmt.Print(ansi.ScreenAlternativeEnter)
	defer fmt.Print(ansi.ScreenAlternativeLeave)

	fmt.Print(ansi.CursorHide)
	defer fmt.Print(ansi.CursorShow)

	title := "n: new, s: search"
	output := []string{}
	dim := must.Do2(term.Size(tty))
	startLine := 0
	startCol := 0
	cursorLine := 0

draw:
	fmt.Print(ansi.EraseEntireScreen)
	fmt.Print(ansi.CursorMoveToHomePosition)
	cfmt.Printf("#R{%s}", title)
	for i := startLine; i < min(len(output), startLine+dim.Rows-1); i++ {
		fmt.Print(ansi.CursorMoveDown(1))
		fmt.Print(ansi.CursorMoveToColumn(1))
		text := output[i]
		text = ellipsis.Ellipsis(text, dim.Cols, cfmt.Sprint("#R{>}"), ellipsis.PosEnd)
		if i-startLine == cursorLine {
			cfmt.Printf("#yR{%s}", text)
		} else {
			fmt.Print(text)
		}
	}

eventLoop:
	for {
		select {
		case key := <-onKeystroke:
			switch string([]byte{key}) {
			case "j":
				if cursorLine == dim.Rows-2 {
					if startLine < len(output)-1 {
						startLine++
					}
				} else {
					cursorLine++
				}
				goto draw
			case "k":
				if cursorLine == 0 {
					if startLine > 0 {
						startLine--
					}
				} else {
					cursorLine--
				}
				goto draw
			case "h":
				if startCol > 0 {
					startCol--
					goto draw
				}
			case "l":
				if startCol < dim.Cols-1 {
					startCol++
					goto draw
				}
			case "n":
				title = cfmt.Sprint("#R{notmuch new}")
				output = []string{}
				startLine = 0
				startCol = 0
				cursorLine = 0
				fmt.Print(ansi.EraseEntireScreen)
				fmt.Print(ansi.CursorMoveToHomePosition)
				fmt.Print(title)
				fmt.Print(ansi.CursorMoveDown(1))
				fmt.Print(ansi.CursorMoveToColumn(1))
				for out := range NotMuchNew() {
					output = append(output, out)
					fmt.Print(ellipsis.Ellipsis(out, dim.Cols, cfmt.Sprint("#R{>}"), ellipsis.PosEnd))
					fmt.Print(ansi.CursorMoveDown(1))
					fmt.Print(ansi.CursorMoveToColumn(1))
					cursorLine++
				}
				goto draw
			case "s":
				output = []string{}
				startLine = 0
				startCol = 0
				cursorLine = 0
				search := ""
			searchInputLoop:
				for {
					title = cfmt.Sprintf("#R{notmuch search %s}", search)
					fmt.Print(ansi.EraseEntireScreen)
					fmt.Print(ansi.CursorMoveToHomePosition)
					fmt.Print(title)
					key := <-onKeystroke
					switch string([]byte{key}) {
					case ansi.InputCR:
						fallthrough
					case ansi.InputLF:
						break searchInputLoop
					case ansi.InputBackSpace:
						fallthrough
					case ansi.InputDelete:
						if search != "" {
							search = search[:len(search)-1]
						}
					default:
						search += string(key)
					}
				}
				fmt.Print(ansi.CursorMoveDown(1))
				fmt.Print(ansi.CursorMoveToColumn(1))
				count := 0
				for out := range NotMuchSearch(search) {
					output = append(output, out)
					if count < dim.Rows-1 {
						fmt.Print(ellipsis.Ellipsis(out, dim.Cols, cfmt.Sprint("#R{>}"), ellipsis.PosEnd))
						fmt.Print(ansi.CursorMoveDown(1))
						fmt.Print(ansi.CursorMoveToColumn(1))
					}
					count++
				}
				goto draw
			case "q":
				fallthrough
			case ansi.InputCtrlC:
				break eventLoop
			}
		case dim = <-onResize:
			goto draw
		}
	}

	return 0
}

func main() {
	os.Exit(run())
}
