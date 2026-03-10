//go:build windows

package term

/*
#include <io.h>
#include <stdio.h>
#include <windows.h>

DWORD term_mode;

int term_make_raw() {
	HANDLE stdin = GetStdHandle(STD_INPUT_HANDLE);
	if (!GetConsoleMode(stdin, &term_mode)) {
		return 1;
	}
	DWORD raw_mode = term_mode;
	raw_mode &= ~(ENABLE_ECHO_INPUT | ENABLE_LINE_INPUT | ENABLE_PROCESSED_INPUT);
	if (!SetConsoleMode(stdin, raw_mode)) {
		return 2;
	}
	return 0;
}

int term_restore() {
	HANDLE stdin = GetStdHandle(STD_INPUT_HANDLE);
	if (!SetConsoleMode(stdin, term_mode)) {
		return 1;
	}
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
