package json_test

import (
	"fmt"
	"testing"

	"github.com/tobiashort/th-utils/pkg/json"
	"github.com/tobiashort/utils-go/strings"
)

func TestToGoStruct(t *testing.T) {
	examples := []string{
		strings.Dedent(`{
  					   |    "string": "Hello, world!",
					   |    "number_integer": 42,
					   |    "number_float": 3.14159,
					   |    "boolean_true": true,
					   |    "boolean_false": false,
					   |    "null_value": null,
					   |    "array": [1, "two", false, null, {"nestedKey": "nestedValue"}],
					   |    "object": {
					   |        "id": 123,
					   |        "name": "Alice",
					   |        "active": true,
					   |        "roles": ["admin", "editor"],
					   |        "profile": {
					   |            "age": 30,
					   |            "score": 99.5
					   |        }
					   |    }
					   |}`),
	}

	for i, example := range examples {
		t.Run(fmt.Sprintf("example%d", i), func(t *testing.T) {
			_, err := json.ToGoStruct([]byte(example))
			if err != nil {
				t.Log(err)
				t.Fail()
			}
		})
	}
}
