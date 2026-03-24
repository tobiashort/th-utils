package main

import (
	"fmt"

	. "github.com/tobiashort/th-utils/lib/option"
)

func foo() Option[string] {
	return Some("barr")
	// return None[string]()
}

func main() {
	switch f := foo(); f {
	case Some("bar"):
		fmt.Println("bar")
	case Some("baz"):
		fmt.Println("baz")
	case None[string]():
		fmt.Println("none")
	default:
		fmt.Println(f.Val)
	}
}
