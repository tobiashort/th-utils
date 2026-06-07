package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/tobiashort/th-utils/lib/must"
)

func TestParse(t *testing.T) {
	text := string(must.Do2(os.ReadFile("testdata/testfile.txt")))
	tokens := Parse(text)
	for i := 0; i < len(tokens); i++ {
		for j := 0; j < len(tokens[i]); j++ {
			t := tokens[i][j]
			if t.Type == TokenAnsi {
				fmt.Print(t.Type)
			} else {
				fmt.Print(tokens[i][j].Literal)
			}
		}
		fmt.Println()
	}
}
