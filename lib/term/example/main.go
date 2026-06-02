package main

import (
	"fmt"

	"github.com/tobiashort/th-utils/lib/ansi"
	"github.com/tobiashort/th-utils/lib/clog"
	"github.com/tobiashort/th-utils/lib/must"
	"github.com/tobiashort/th-utils/lib/term"
)

func main() {
	clog.Infof("isTerminal=%v", term.IsTerminal())

	tty := must.Do2(term.OpenTTY())
	defer tty.Close()

	must.Do(term.MakeRaw(tty))
	clog.Info("Switched to raw mode")
	fmt.Print(ansi.CursorMoveToColumn(0))

	must.Do(term.Restore(tty))
	clog.Info("Restored to original mode")

	cols, lines := must.Do3(term.Size(tty))
	clog.Infof("cols=%d lines=%d", cols, lines)
}
