package main

func ExampleJoin() {
	args := Args{
		Files:     []string{"./testdata/file1", "./testdata/file2", "./testdata/file3"},
		Separator: " ",
	}
	join(args)
	// Output: a 0 *
	// b 1 %
	// c 2 &
	// d 3 #
	//  4 @
}
