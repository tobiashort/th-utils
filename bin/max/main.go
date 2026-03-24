package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"

	"github.com/tobiashort/clap-go"
	"github.com/tobiashort/utils-go/assert"
)

type Args struct {
	IgnoreNaN bool `clap:"desc='Ignores not numeric values and treats them as 0.'"`
}

func main() {
	args := Args{}
	clap.Description("Reads from Stdin and prints the maximum.")
	clap.Parse(&args)

	max := big.NewFloat(0)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		num, ok := new(big.Float).SetString(text)
		if ok {
			if num.Cmp(max) > 0 {
				max = num
			}
		} else {
			if !args.IgnoreNaN {
				fmt.Fprintf(os.Stderr, "%s is not a number", text)
				os.Exit(1)
			}
		}
	}
	assert.Nil(scanner.Err(), "scanner error: %w", scanner.Err())

	fmt.Println(max)
}
