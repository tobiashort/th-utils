package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/tobiashort/th-utils/lib/assert"
	"github.com/tobiashort/th-utils/lib/clap"
)

type Args struct {
	Count int `clap:"desc='Number of lines to skip'"`
}

func main() {
	args := Args{}
	clap.Parse(&args)

	count := 0
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if count < args.Count {
			count++
			continue
		}
		fmt.Println(scanner.Text())
	}
	assert.Nil(scanner.Err(), "scanner error")
}
