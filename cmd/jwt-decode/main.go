package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/tobiashort/clap-go"

	"github.com/tobiashort/th-utils/pkg/jwt"
)

type Args struct {
	JWT string `clap:"positional,desc='The JWT string. Reads from Stdin if not specified.'"`
}

func main() {
	args := Args{}
	clap.Parse(&args)

	input := ""
	if args.JWT != "" {
		input = args.JWT
	} else {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
		input = strings.TrimSpace(string(data))
	}

	decoded, err := jwt.Decode(input)
	if err != nil {
		panic(err)
	}

	fmt.Println(decoded)
}
