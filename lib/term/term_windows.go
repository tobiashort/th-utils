//go:build windows

package term

/*
#include <io.h>
#include <stdio.h>
#include <windows.h>

DWORD term_mode;

int term_make_raw(HANDLE h_stdin) {
	//HANDLE h_stdin = GetStdHandle(STD_INPUT_HANDLE);
	if (!GetConsoleMode(h_stdin, &term_mode)) {
		return 1;
	}
	DWORD raw_mode = term_mode;
	raw_mode &= ~(ENABLE_ECHO_INPUT | ENABLE_LINE_INPUT | ENABLE_PROCESSED_INPUT);
	raw_mode |= ENABLE_VIRTUAL_TERMINAL_INPUT;
	if (!SetConsoleMode(h_stdin, raw_mode)) {
		return 2;
	}
	return 0;
}

int term_restore(HANDLE h_stdin) {
	//HANDLE h_stdin = GetStdHandle(STD_INPUT_HANDLE);
	if (!SetConsoleMode(h_stdin, term_mode)) {
		return 1;
	}
	return 0;
}

int term_size(int *cols, int *rows) {
    CONSOLE_SCREEN_BUFFER_INFO info;

    HANDLE h = GetStdHandle(STD_OUTPUT_HANDLE);
    if (h == INVALID_HANDLE_VALUE) {
        return 1;
    }

    if (!GetConsoleScreenBufferInfo(h, &info)) {
        return 2;
    }

    *cols = info.srWindow.Right - info.srWindow.Left + 1;
    *rows = info.srWindow.Bottom - info.srWindow.Top + 1;

    return 0;
}
*/
import "C"

import (
	"fmt"
	"os"
	"syscall"
	"time"
	"unsafe"
)

// Checks whether stdout is a terminal or not
func IsTerminal() bool {
	return C._isatty(C.int(1)) != 0
}

func OpenTTY() (*os.File, error) {
	return os.OpenFile("CONIN$", os.O_RDWR, 0)
}

func MakeRaw(tty *os.File) error {
	h := syscall.Handle(tty.Fd())
	ret := int(C.term_make_raw(C.HANDLE(unsafe.Pointer(h))))
	switch ret {
	case 0:
		return nil
	case 1:
		return fmt.Errorf("Unable to get current console mode")
	case 2:
		return fmt.Errorf("Unable to set console mode to raw")
	default:
		panic("unreachable")
	}
}

func Restore(tty *os.File) error {
	h := syscall.Handle(tty.Fd())
	ret := int(C.term_restore(C.HANDLE(unsafe.Pointer(h))))
	switch ret {
	case 0:
		return nil
	case 1:
		return fmt.Errorf("Unable to set console mode to raw")
	default:
		panic("unreachable")
	}
}

func Size(tty *os.File) (Dim, error) {
	var cols, rows C.int
	ret := C.term_size(&cols, &rows)
	if ret != 0 {
		return Dim{}, fmt.Errorf("failed to get terminal size (code %d)", int(ret))
	}
	return Dim{Cols: int(cols), Rows: int(rows)}, nil
}

func OnResize(tty *os.File) chan Dim {
	ch := make(chan Dim, 1)
	go func() {
		dimPrev := Dim{}
		for {
			dim, err := Size(tty)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			} else {
				if dim.Cols != dimPrev.Cols || dim.Rows != dimPrev.Rows {
					ch <- dim
					dimPrev = dim
				}
			}
			time.Sleep(250 * time.Millisecond)
		}
	}()
	return ch
}
