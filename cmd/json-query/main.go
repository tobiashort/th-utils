package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/expr-lang/expr"

	"github.com/tobiashort/clap-go"
	"github.com/tobiashort/utils-go/must"
	strings2 "github.com/tobiashort/utils-go/strings"
)

type Args struct {
	Expr string `clap:"mandatory,positional,desc='The expression (see https://expr-lang.org/docs/language-definition)'"`
	File string `clap:"positional,desc='The json file to read. Reads from Stdin if not specified.'"`
}

func main() {
	args := Args{}
	clap.Example(
		strings2.Dedent(
			`Expr‑lang can only operate on maps, not on arrays.
			|Consequently, JSON arrays such as the one below are interpreted
			|as a map with a single key named "items".
            |
            |File array.json:
            |----------------
            |[
            |  {
            |    "name": "Alex",
            |    "age": 21,
            |    "groups": ["admin", "support", "user"],
            |    "profile": {
            |      "complete": true
            |    }
            |  },
            |  "bla",
            |  {
            |    "name": "Judith",
            |    "age": 59,
            |    "groups": ["support", "user"],
            |    "profile": {
            |      "complete": true
            |    }
            |  },
            |  {
            |    "name": "Willow",
            |    "age": 12,
            |    "groups": ["user"],
            |    "profile": {
            |      "complete": false
            |    }
            |  }
            |]
            |
            |Running the query:
            |-----------------
            |th-json-query 'filter(items, type(#) == "map") | filter("user" in .groups) | filter(.profile.complete == true) | map(.name) | join("\n")' array.json
            |
            |Produces:
            |---------
            |Alex
            |Judith
            |`))
	clap.Parse(&args)

	var data []byte
	if args.File != "" {
		data = must.Do2(os.ReadFile(args.File))
	} else {
		data = must.Do2(io.ReadAll(os.Stdin))
	}

	var env any
	if strings.HasPrefix(strings.TrimSpace(string(data)), "[") {
		var v []any
		must.Do(json.Unmarshal(data, &v))
		env = map[string]any{"items": v}
	} else {
		var v map[string]any
		must.Do(json.Unmarshal(data, &v))
		env = v
	}

	prog := must.Do2(expr.Compile(args.Expr, expr.Env(env)))
	out := must.Do2(expr.Run(prog, env))
	fmt.Printf("%v", out)
}
