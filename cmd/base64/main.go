package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"github.com/tobiashort/clap-go"
)

type Args struct {
	Raw    bool `clap:"description='no padding'"`
	Url    bool `clap:"description='url and file safe'"`
	Decode bool `clap:"description='decode otherwise encode'"`
}

func main() {
	args := Args{}
	clap.Parse(&args)
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	var encoding *base64.Encoding
	if args.Raw {
		if args.Url {
			encoding = base64.RawURLEncoding
		} else {
			encoding = base64.RawStdEncoding
		}
	} else {
		if args.Url {
			encoding = base64.URLEncoding
		} else {
			encoding = base64.StdEncoding
		}
	}
	if args.Decode {
		decoded, err := encoding.DecodeString(string(data))
		fmt.Print(string(decoded))
		if err != nil {
			fmt.Print(err)
		}
	} else {
		fmt.Print(encoding.EncodeToString(data))
	}
}
