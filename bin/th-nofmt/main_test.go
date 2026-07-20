package main

func Example_fmt1() {
	src := `   package main
func main() {
	type a struct{s string; i int}
} `
	Fmt(src)
	// Output: package main
	//
	// func main() {
	// 	type a struct {
	// 		s string
	// 		i int
	// 	}
	// }
}

func Example_fmt2() {
	src := `   package main
func main() {
	//nofmt:enable
	type a struct{s string; i int}
	type b struct{s string; i int}
	//nofmt:disable
} `
	Fmt(src)
	// Output: package main
	//
	// func main() {
	// 	//nofmt:enable
	// 	type a struct{s string; i int}
	// 	type b struct{s string; i int}
	// 	//nofmt:disable
	// }
}

func Example_fmt3() {
	src := `   package main
func main() {
	//nofmt:enable
	type a struct{s string; i int}
	//nofmt:disable
} `
	Fmt(src)
	// Output: package main
	//
	// func main() {
	// 	type a struct{s string; i int} //nofmt
	// }
}

func Example_fmt4() {
	src := `   package main
func main() {
	type a struct{s string; i int} //nofmt
} `
	Fmt(src)
	// Output: package main
	//
	// func main() {
	// 	type a struct{s string; i int} //nofmt
	// }
}
