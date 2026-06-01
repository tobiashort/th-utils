//go:build windows

package term

/*
#include <io.h>
#include <stdio.h>
#include <windows.h>

DWORD term_mode;

int term_make_raw() {
	HANDLE h_stdin = GetStdHandle(STD_INPUT_HANDLE);
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

int term_restore() {
	HANDLE h_stdin = GetStdHandle(STD_INPUT_HANDLE);
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

import "fmt"

// Checks whether stdout is a terminal or not
func IsTerminal() bool {
	return C._isatty(C.int(1)) != 0
}

func MakeRaw() error {
	ret := int(C.term_make_raw())
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

func Restore() error {
	ret := int(C.term_restore())
	switch ret {
	case 0:
		return nil
	case 1:
		return fmt.Errorf("Unable to set console mode to raw")
	default:
		panic("unreachable")
	}
}

func Size() (int, int, error) {
	var cols, rows C.int
	ret := C.term_size(&cols, &rows)
	if ret != 0 {
		return 0, 0, fmt.Errorf("failed to get terminal size (code %d)", int(ret))
	}
	return int(cols), int(rows), nil
}
