package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestFormat(t *testing.T) {
	input := []byte(`["hello world",42,3.14159,true,false,null,{"id":1,"name":"Alice","active":true,"scores":[10,20,30],"profile":{"age":28,"city":"Zurich"}},[1,2,3,"mixed",false],{"event":"login","timestamp":"2026-03-27T12:34:56Z","meta":{"ip":"192.168.0.1","device":"mobile"}},-17,0,"",["nested",["deeply",["nested"]]]]`)

	for i := range 10 {
		t.Run(fmt.Sprintf("Run%d", i), func(t *testing.T) {
			actual := format(input)
			expected := []byte(`[
  "hello world",
  42,
  3.14159,
  true,
  false,
  null,
  {
    "id": 1,
    "name": "Alice",
    "active": true,
    "scores": [
      10,
      20,
      30
    ],
    "profile": {
      "age": 28,
      "city": "Zurich"
    }
  },
  [
    1,
    2,
    3,
    "mixed",
    false
  ],
  {
    "event": "login",
    "timestamp": "2026-03-27T12:34:56Z",
    "meta": {
      "ip": "192.168.0.1",
      "device": "mobile"
    }
  },
  -17,
  0,
  "",
  [
    "nested",
    [
      "deeply",
      [
        "nested"
      ]
    ]
  ]
]`)

			if !bytes.Equal(actual, expected) {
				t.Fatal(actual)
			}
		})
	}
}
