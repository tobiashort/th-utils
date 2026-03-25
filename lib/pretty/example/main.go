package main

import (
	"fmt"

	"github.com/tobiashort/th-utils/lib/pretty"
)

func main() {
	fmt.Println(pretty.Sprint([]string{"a", "b", "c"}))
	fmt.Println()
	fmt.Println(pretty.Sprint([][]string{{"a", "b", "c"}, {"d", "e", "f"}}))
	fmt.Println()
	fmt.Println(pretty.Sprint(map[string]any{"a": "apple", "b": "banana", "c": "cidre", "d": map[string]int{"amount": 42, "weight": 23}}))
	fmt.Println()

	type inner struct {
		a map[string]int
		b [3]float32
	}

	type strct struct {
		A string
		B string
		c string
		i inner
	}

	s := strct{
		A: "apple",
		B: "banana",
		c: "cidre",
		i: inner{
			a: map[string]int{
				"bla": 1,
				"foo": 2,
				"bar": 43,
			},
			b: [3]float32{1.1, 1.2, 1.3},
		},
	}
	fmt.Println("s:", pretty.Sprint(s))
}
