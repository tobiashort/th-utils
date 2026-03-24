package main

import (
	"fmt"

	"github.com/tobiashort/th-utils/lib/must"
)

func foo() error {
	return nil
}

func foo2() (string, error) {
	return "bar", nil
}

func foo3() (string, int, error) {
	return "bar", len("bar"), nil
}

func main() {
	must.Do(foo())

	str := must.Do2(foo2())
	fmt.Println(str)

	str, len := must.Do3(foo3())
	fmt.Println(str, len)
}
