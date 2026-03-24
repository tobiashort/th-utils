package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/tobiashort/clap-go"
	"github.com/tobiashort/utils-go/must"
)

type Args struct {
	Secs string `clap:"positional,desc='The seconds to convert to RFC3339 format. Reads from Stdin if not specified.'"`
}

func main() {
	args := Args{}
	clap.Parse(&args)

	input := args.Secs
	if input == "" {
		input = string(must.Do2(io.ReadAll(os.Stdin)))
	}

	secs := must.Do2(strconv.Atoi(input))
	fmt.Println(time.Unix(int64(secs), 0).Format(time.RFC3339))
}
