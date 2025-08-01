package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func usage() {
	fmt.Print(`Usage: rand-string [-w WIDTH] [-l] [-u] [-n] [-s] [-a ALPHABET] [-c COUNT]
Generates random strings.


EXAMPLES

	./rand-string -l -u -n -w 12 -c 100
	Generates 100 twelve character wide random strings with
	lowercase and uppercase characters and numbers.


	./rand-string -a "abc123" -w 10 -c 50
	Generates 50 ten character wide random strings that match
	the given alphabet.


FLAGS

`)
	flag.PrintDefaults()
}

const (
	Symbols      = "!\"#$%&'()*+,-./:;<>=?@[\\]^_`{|}~"
	Lowercase    = "abcdefghijklmnopqrstuvwxyz"
	Uppercase    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Numbers      = "0123456789"
	Hexadecimals = Numbers + "abcdef"
)

func main() {
	var flagWidth int
	var flagLowercase bool
	var flagUppercase bool
	var flagNumbers bool
	var flagSymbols bool
	var flagHex bool
	var flagAlphabet string
	var flagCount int

	flag.Usage = usage
	flag.IntVar(&flagWidth, "w", 12, "width in characters")
	flag.BoolVar(&flagLowercase, "l", false, "lowercase characters [a-z]")
	flag.BoolVar(&flagUppercase, "u", false, "uppercase characters [A-Z]")
	flag.BoolVar(&flagNumbers, "n", false, "numbers [0-9]")
	flag.BoolVar(&flagSymbols, "s", false, fmt.Sprintf("symbols [%s]", Symbols))
	flag.BoolVar(&flagHex, "H", false, "hexadecimals [0-9a-f]")
	flag.StringVar(&flagAlphabet, "a", "", "alphabet to be used instead")
	flag.IntVar(&flagCount, "c", 1, "the amount of strings to be generated")
	flag.Parse()

	if flagWidth < 1 {
		fmt.Fprintf(os.Stderr, "w must be greater than 0")
		usage()
		os.Exit(1)
	}

	if flagCount < 1 {
		fmt.Fprintf(os.Stderr, "c must be greater than 0")
		usage()
		os.Exit(1)
	}

	var alphabet string

	if flagAlphabet != "" {
		alphabet = flagAlphabet
	} else if flagHex {
		alphabet = Hexadecimals
	} else {
		if flagLowercase {
			alphabet += Lowercase
		}
		if flagUppercase {
			alphabet += Uppercase
		}
		if flagNumbers {
			alphabet += Numbers
		}
		if flagSymbols {
			alphabet += Symbols
		}
	}

	if alphabet == "" {
		alphabet = Lowercase + Uppercase + Numbers + Symbols
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	for range flagCount {
		b := make([]byte, flagWidth)
		for i := range flagWidth {
			b[i] = alphabet[rnd.Intn(len(alphabet))]
		}
		fmt.Println(string(b))
	}
}
