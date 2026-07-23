package slices

import "fmt"

func ExampleZip() {
	a := []int{1, 2, 3, 4}
	b := []string{"apple", "banana", "citron"}
	c := []float64{0.2, 0.1, 0.5, 0.2}
	ret := Zip(struct {First  int; Second string; Third float64}{}, a, b, c) //nofmt
	fmt.Println(ret)
	// Output: [{1 apple 0.2} {2 banana 0.1} {3 citron 0.5} {4  0.2}]
}
