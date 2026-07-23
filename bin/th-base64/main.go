package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"github.com/tobiashort/th-utils/lib/clap"
	"github.com/tobiashort/th-utils/lib/must"
)

type Args struct {
	Raw    bool   `clap:"desc='no padding'"`
	Url    bool   `clap:"desc='url and file safe'"`
	Decode bool   `clap:"desc='decode otherwise encode'"`
	Text   string `clap:"positional,desc='The string to encode/decode. Reads from stdin if not specified'"`
}

func main() {
	args := Args{}
	clap.Parse(&args)
	data := []byte(args.Text)
	if len(data) == 0 {
		data = must.Do2(io.ReadAll(os.Stdin))
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
