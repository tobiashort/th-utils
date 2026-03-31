package ansi

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/tobiashort/th-utils/lib/must"
)

type CursorControl = string

const (
	CursorHide               CursorControl = "\033[?25l"
	CursorShow               CursorControl = "\033[?25h"
	CursorMoveToHomePosition CursorControl = "\033[H"
	CursorSavePosition       CursorControl = "\033[s"
	CursorRestorePosition    CursorControl = "\033[u"
	CursorCurrentPosition    CursorControl = "\033[6n"
)

func CursorMoveTo(line, column int) CursorControl {
	return fmt.Sprintf("\033[%d;%dH", line, column)
}

func CursorMoveUp(lines int) CursorControl {
	return fmt.Sprintf("\033[%dA", lines)
}

func CursorMoveDown(lines int) CursorControl {
	return fmt.Sprintf("\033[%dB", lines)
}

func CursorMoveRight(columns int) CursorControl {
	return fmt.Sprintf("\033[%dC", columns)
}

func CursorMoveLeft(columns int) CursorControl {
	return fmt.Sprintf("\033[%dD", columns)
}

func CursorMoveToColumn(column int) CursorControl {
	return fmt.Sprintf("\033[%dG", column)
}

func CursorGetCurrentPosition() (row int, col int) {
	fmt.Fprint(os.Stderr, CursorCurrentPosition)
	reader := bufio.NewReader(os.Stdin)
	response := must.Do2(reader.ReadString('R'))
	response = strings.TrimPrefix(response, "\x1b[")
	response = strings.TrimSuffix(response, "R")
	fmt.Sscanf(response, "%d;%d", &row, &col)
	return
}
