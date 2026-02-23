package main

import (
	"fmt"
	"os"
	"time"

	"github.com/tobiashort/clap-go"
	"github.com/tobiashort/th-utils/pkg/random"
)

const (
	Symbols      = "!\"#$%&'()*+,-./:;<>=?@[\\]^_`{|}~"
	Lowercase    = "abcdefghijklmnopqrstuvwxyz"
	Uppercase    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Numbers      = "0123456789"
	Hexadecimals = Numbers + "abcdef"
)

type Args struct {
	Width     int    `clap:"default=12,desc='Width in characters'"`
	Lowercase bool   `clap:"desc='Use lowercase characters'"`
	Uppercase bool   `clap:"desc='Use uppercase characters'"`
	Numbers   bool   `clap:"desc='Use numbers'"`
	Symbols   bool   `clap:"desc='Use symbols'"`
	Hex       bool   `clap:"short=x,conflicts='Alphabet,Lowercase,Uppercase,Numbers,Symbols',desc='Use hexadecimals as the alphabet'"`
	Alphabet  string `clap:"conflicts='Hex,Lowercase,Uppercase,Numbers,Symbols',desc='The custom alphabet to be used'"`
	Amount    int    `clap:"short=c,default=1,desc='The amount of strings to be generated'"`
}

func main() {
	args := Args{}

	clap.Example(`$ rand-string -lun -w 12 -c 100
Generates 100 twelve character wide random strings with
lowercase and uppercase characters and numbers.

$ rand-string -a "abc123" -w 10 -c 50
Generates 50 ten character wide random strings that match
the given alphabet.`)

	clap.Parse(&args)

	if args.Width < 1 {
		fmt.Fprintf(os.Stderr, "w must be greater than 0")
		os.Exit(1)
	}

	if args.Amount < 1 {
		fmt.Fprintf(os.Stderr, "c must be greater than 0")
		os.Exit(1)
	}

	var alphabet string

	if args.Alphabet != "" {
		alphabet = args.Alphabet
	} else if args.Hex {
		alphabet = Hexadecimals
	} else {
		if args.Lowercase {
			alphabet += Lowercase
		}
		if args.Uppercase {
			alphabet += Uppercase
		}
		if args.Numbers {
			alphabet += Numbers
		}
		if args.Symbols {
			alphabet += Symbols
		}
	}

	if alphabet == "" {
		alphabet = Lowercase + Uppercase + Numbers + Symbols
	}

	seed := time.Now().UnixNano()

	for range args.Amount {
		fmt.Println(random.String(alphabet, args.Width, seed))
	}
}
