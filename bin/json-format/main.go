package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/tobiashort/th-utils/lib/clap"
	"github.com/tobiashort/th-utils/lib/must"
	"github.com/tobiashort/th-utils/lib/orderedmap"
)

type Args struct {
	JSON string `clap:"positional,desc='The JSON string. Reads from Stdin if not specified.'"`
}

func unmarshal(r json.RawMessage) any {
	r = bytes.TrimSpace(r)
	switch r[0] {
	case '[':
		var v []any
		var temp []json.RawMessage
		must.Do(json.Unmarshal(r, &temp))
		for _, t := range temp {
			v = append(v, unmarshal(t))
		}
		return v
	case '{':
		var v orderedmap.OrderedMap[string, any]
		must.Do(json.Unmarshal(r, &v))
		return v
	default:
		var v any
		must.Do(json.Unmarshal(r, &v))
		return v
	}
}

func format(input []byte) []byte {
	unmarshalled := unmarshal(input)
	return must.Do2(json.MarshalIndent(unmarshalled, "", "  "))
}

func main() {
	args := Args{}
	clap.Parse(&args)

	var input []byte
	if args.JSON == "" {
		var err error
		input, err = io.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
	} else {
		input = []byte(args.JSON)
	}

	output := format(input)
	fmt.Print(string(output))
}
