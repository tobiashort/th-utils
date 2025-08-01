package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	var raw bool
	var url bool
	var decode bool
	flag.BoolVar(&raw, "r", false, "no padding")
	flag.BoolVar(&url, "u", false, "url and file safe")
	flag.BoolVar(&decode, "d", false, "decode otherwise encode")
	flag.Parse()
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	var encoding *base64.Encoding
	if raw {
		if url {
			encoding = base64.RawURLEncoding
		} else {
			encoding = base64.RawStdEncoding
		}
	} else {
		if url {
			encoding = base64.URLEncoding
		} else {
			encoding = base64.StdEncoding
		}
	}
	if decode {
		decoded, err := encoding.DecodeString(string(data))
		fmt.Print(string(decoded))
		if err != nil {
			fmt.Print(err)
		}
	} else {
		fmt.Print(encoding.EncodeToString(data))
	}
}
