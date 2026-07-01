package main

import "fmt"

func Example_linesToTake() {
	for _, l := range linesToTake(3, "6,8-10,12", 2, 20) {
		fmt.Println(l)
	}
	// Output: 1
	// 2
	// 3
	// 6
	// 8
	// 9
	// 10
	// 12
	// 19
	// 20
}

func Example_run() {
	args := Args{}
	args.Head = 3
	args.Lines = "6,8-10,12"
	args.Tail = 2
	args.File = "./testdata/testdata.txt"
	run(args)
	// Output: take
	// take
	// take
	// take
	// take
	// take
	// take
	// take
	// take
	// take
}
