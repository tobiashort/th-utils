package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/tobiashort/clap-go"
)

const (
	Symbols      = "!\"#$%&'()*+,-./:;<>=?@[\\]^_`{|}~"
	Lowercase    = "abcdefghijklmnopqrstuvwxyz"
	Uppercase    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Numbers      = "0123456789"
	Hexadecimals = Numbers + "abcdef"
)

type Args struct {
	Width     int    `clap:"default-value=12,description='Width in characters'"`
	Lowercase bool   `clap:"description='Use lowercase characters'"`
	Uppercase bool   `clap:"description='Use uppercase characters'"`
	Numbers   bool   `clap:"description='Use numbers'"`
	Symbols   bool   `clap:"description='Use symbols'"`
	Hex       bool   `clap:"short=x,conflicts-with='Alphabet,Lowercase,Uppercase,Numbers,Symbols',description='Use hexadecimals as the alphabet'"`
	Alphabet  string `clap:"conflicts-with='Hex,Lowercase,Uppercase,Numbers,Symbols',description='The custom alphabet to be used'"`
	Amount    int    `clap:"short=c,default-value=1,description='The amount of strings to be generated'"`
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

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	for range args.Amount {
		b := make([]byte, args.Width)
		for i := range args.Width {
			b[i] = alphabet[rnd.Intn(len(alphabet))]
		}
		fmt.Println(string(b))
	}
}
