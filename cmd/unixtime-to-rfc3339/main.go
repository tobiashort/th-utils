package main

import (
	"fmt"
	"time"

	"github.com/tobiashort/clap-go"
)

type Args struct {
	Secs int `clap:"positional,mandatory,description='The seconds to convert to RFC3339 format'"`
}

func main() {
	args := Args{}
	clap.Parse(&args)
	fmt.Println(time.Unix(int64(args.Secs), 0).Format(time.RFC3339))
}
