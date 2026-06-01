package main

import (
	"fmt"

	"github.com/tobiashort/th-utils/lib/term"
)

func main() {
	fmt.Println("Is a tty:", term.IsTerminal())

	err := term.MakeRaw()
	if err != nil {
		panic(err)
	}
	fmt.Println("Switched to raw mode")

	err = term.Restore()
	if err != nil {
		panic(err)
	}
	fmt.Println("Restored to original mode")

	cols, lines, err := term.Size()
	if err != nil {
		panic(err)
	}
	fmt.Printf("cols:%d,lines:%d\n", cols, lines)
}
