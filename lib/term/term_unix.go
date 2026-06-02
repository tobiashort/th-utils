//go:build linux || darwin || freebsd

package term

/*
#include <stdio.h>
#include <unistd.h>
#include <termios.h>
#include <sys/ioctl.h>

struct termios oldt, rawt;

int term_make_raw(int fd) {
	if (tcgetattr(fd, &oldt) != 0) {
		return 1;
	}
	rawt = oldt;
	cfmakeraw(&rawt);
	if (tcsetattr(fd, TCSANOW, &rawt) != 0) {
		return 2;
	}
	return 0;
}

int term_restore(int fd) {
	if (tcsetattr(fd, TCSANOW, &oldt) != 0) {
		return 1;
	}
	return 0;
}

int term_size(int fd, int *cols, int *rows) {
    struct winsize ws;

    if (ioctl(fd, TIOCGWINSZ, &ws) != 0) {
        return 1;
    }

    *cols = ws.ws_col;
    *rows = ws.ws_row;

    return 0;
}
*/
import "C"

import (
	"fmt"
	"os"
)

// Checks whether stdout is a terminal or not
func IsTerminal() bool {
	return int(C.isatty(C.int(os.Stdout.Fd()))) == 1
}

func OpenTTY() (*os.File, error) {
	return os.OpenFile("/dev/tty", os.O_RDWR, 0)
}

func MakeRaw(tty *os.File) error {
	ret := int(C.term_make_raw(C.int(tty.Fd())))
	switch ret {
	case 0:
		return nil
	case 1:
		return fmt.Errorf("Unable to get current term mode")
	case 2:
		return fmt.Errorf("Unable to set term mode to raw")
	default:
		panic("unreachable")
	}
}

func Restore(tty *os.File) error {
	ret := int(C.term_restore(C.int(tty.Fd())))
	switch ret {
	case 0:
		return nil
	case 1:
		return fmt.Errorf("Unable to set term mode to raw")
	default:
		panic("unreachable")
	}
}

func Size(tty *os.File) (int, int, error) {
	var cols, rows C.int
	fd := C.int(tty.Fd())
	if ret := C.term_size(fd, &cols, &rows); ret != 0 {
		C.perror(C.CString("ioctl"))
		return 0, 0, fmt.Errorf("failed to get terminal size (code %d)", int(ret))
	}
	return int(cols), int(rows), nil
}
