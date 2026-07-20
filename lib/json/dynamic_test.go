package json_test

import (
	"fmt"

	"github.com/tobiashort/th-utils/lib/json"
)

func ExampleDynamic_exmple1() {
	s := `{
    "glossary": {
        "title": "example glossary",
		"GlossDiv": {
            "title": "S",
			"GlossList": {
                "GlossEntry": {
                    "ID": "SGML",
					"SortAs": "SGML",
					"GlossTerm": "Standard Generalized Markup Language",
					"Acronym": "SGML",
					"Abbrev": "ISO 8879:1986",
					"GlossDef": {
                        "para": "A meta-markup language, used to create markup languages such as DocBook.",
						"GlossSeeAlso": ["GML", "XML"]
                    },
					"GlossSee": "markup"
                }
            }
        }
    }
}`
	dj := json.Dynamic([]byte(s))
	fmt.Println(dj.Value("glossary").Value("title").String())
	fmt.Println(dj.Value("glossary").Value("GlossDiv").Value("GlossList").Value("GlossEntry").Value("GlossDef").Value("GlossSeeAlso").StringSlice()[1])
	// Output: example glossary
	// XML
}

func ExampleDynamic_example2() {
	s := `[
  "hello",
  42,
  3.14,
  true,
  false,
  null,
  {
    "name": "Alice",
    "age": 30
  },
  [
    1,
    2,
    3
  ]
]`
	dj := json.Dynamic([]byte(s))
	fmt.Println(dj.Index(3).Bool())
	fmt.Println(dj.Index(6).Value("age").Int())
	fmt.Println(dj.Index(7).Index(2).Int())
	// Output: true
	// 30
	// 3
}
