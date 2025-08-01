package main

import (
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
)

func usage() {
	fmt.Fprint(os.Stderr, `Usage: cols [-n COLS] [FILE]

  ATTENTION
	Be aware that this program loads the whole content
	of the file into memory.

  FILE
	Optional file otherwise reads from stdin.

`)

	flag.PrintDefaults()
}

func textToCols(in string, nCols int) string {
	text := strings.TrimSpace(in)
	text = strings.ReplaceAll(text, "\r\n", "\n")
	lines := strings.Split(text, "\n")
	rows := int(math.Ceil(float64(len(lines)) / float64(nCols)))
	table := make([][]string, nCols)
	for i := range nCols {
		table[i] = make([]string, rows)
	}

	col := 0
	row := 0
	for i, line := range lines {
		col = (i / rows) % nCols
		row = i % rows
		table[col][row] = line
	}

	for _, col := range table {
		width := 0
		for _, colText := range col {
			if len(colText) > width {
				width = len(colText)
			}
		}
		format := fmt.Sprintf("%%-%ds", width)
		for i, colText := range col {
			col[i] = fmt.Sprintf(format, colText)
		}
	}

	out := ""
	for row := 0; row < rows; row++ {
		for col := 0; col < nCols; col++ {
			out += table[col][row]
			if col < nCols-1 {
				out += "    "
			} else {
				out += "\n"
			}
		}
	}
	return out
}

func main() {
	nCols := flag.Int("n", 4, "Number of columns")
	flag.Usage = usage
	flag.Parse()

	var file *os.File

	if flag.NArg() > 1 {
		usage()
		os.Exit(1)
	} else if flag.NArg() == 1 {
		var err error
		file, err = os.Open(flag.Arg(0))
		if err != nil {
			panic(err)
		}
	} else {
		file = os.Stdin
	}

	if *nCols < 2 {
		io.Copy(os.Stdout, file)
		os.Exit(0)
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	text := string(fileBytes)
	fmt.Print(textToCols(text, *nCols))
}
