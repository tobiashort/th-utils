//go:build linux || darwin

package term

/*
#include <unistd.h>
#include <termios.h>

struct termios oldt, rawt;

int term_make_raw() {
	if (tcgetattr(STDIN_FILENO, &oldt) != 0) {
		return 1;
	}
	rawt = oldt;
	cfmakeraw(&rawt);
	if (tcsetattr(STDIN_FILENO, TCSANOW, &rawt) != 0) {
		return 2;
	}
	return 0;
}

int term_restore() {
	if (tcsetattr(STDIN_FILENO, TCSANOW, &oldt) != 0) {
		return 1;
	}
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

func MakeRaw() error {
	ret := int(C.term_make_raw())
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

func Restore() error {
	ret := int(C.term_restore())
	switch ret {
	case 0:
		return nil
	case 1:
		return fmt.Errorf("Unable to set term mode to raw")
	default:
		panic("unreachable")
	}
}
