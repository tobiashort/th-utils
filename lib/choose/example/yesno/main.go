package main

import (
	"fmt"

	"github.com/tobiashort/th-utils/lib/choose"
)

func main() {
	if choose.YesNo("Should I?", choose.DEFAULT_NONE) {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}

	if choose.YesNo("Should he?", choose.DEFAULT_YES) {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}

	if choose.YesNo("Should they?", choose.DEFAULT_NO) {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}
