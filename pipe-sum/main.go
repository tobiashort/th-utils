package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var sum float64

	for {
		val, err := reader.ReadString('\n')
		val = strings.TrimSpace(val)

		if err != nil && err != io.EOF {
			panic(err)
		}

		if val != "" {
			f, err := strconv.ParseFloat(val, 64)

			if err != nil {
				panic(err)
			}

			sum += f
		}

		if err == io.EOF {
			break
		}
	}

	fmt.Println(sum)
}
