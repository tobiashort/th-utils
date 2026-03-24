package main

import (
	"fmt"

	"github.com/tobiashort/th-utils/lib/ternary"
)

func main() {
	age := 18
	access := ternary.IfThenElse(age < 18, "Access denied", "Access granted")
	fmt.Println(access)
}
