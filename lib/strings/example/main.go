package main

import (
	"fmt"

	strings2 "github.com/tobiashort/th-utils/lib/strings"
)

func main() {
	fmt.Println("-------------")
	s := `Lorem ipsum dolor sit amet,
		 |consectetur adipiscing elit.
		 |Curabitur justo tellus, facilisis nec efficitur dictum,
		 |fermentum vitae ligula. Sed eu convallis sapien.`
	fmt.Println(s)
	fmt.Println("-------------")
	fmt.Println(strings2.Dedent(s))
	fmt.Println("-------------")
}
