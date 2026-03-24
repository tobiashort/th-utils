package orderedmap

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestLen(t *testing.T) {
	m := NewOrderedMap[int, string]()
	m.Put(1, "apple")
	m.Put(2, "banana")
	m.Put(3, "citrus")
	m.Put(4, "dragonfruit")
	if m.Len() != 4 {
		t.Errorf("expected 4, got %d", m.Len())
	}
}

func TestKeys(t *testing.T) {
	m := NewOrderedMap[int, string]()
	m.Put(1, "apple")
	m.Put(2, "banana")
	m.Put(3, "citrus")
	m.Put(4, "dragonfruit")
	if !reflect.DeepEqual(m.Keys(), []int{1, 2, 3, 4}) {
		t.Error(m.Keys())
	}
}

func TestValues(t *testing.T) {
	m := NewOrderedMap[int, string]()
	m.Put(1, "apple")
	m.Put(2, "banana")
	m.Put(3, "citrus")
	m.Put(4, "dragonfruit")
	if !reflect.DeepEqual(m.Values(), []string{"apple", "banana", "citrus", "dragonfruit"}) {
		t.Error(m.Values())
	}
}

func TestDelete(t *testing.T) {
	m := NewOrderedMap[int, string]()
	m.Put(1, "apple")
	m.Put(2, "banana")
	m.Put(3, "citrus")
	m.Put(4, "dragonfruit")
	m.Del(3)
	if !reflect.DeepEqual(m.Keys(), []int{1, 2, 4}) {
		t.Error(m.Keys())
	}
}

func TestRange(t *testing.T) {
	m := NewOrderedMap[int, string]()
	m.Put(1, "apple")
	m.Put(2, "banana")
	m.Put(3, "citrus")
	m.Put(4, "dragonfruit")
	keys := make([]int, 0)
	values := make([]string, 0)
	for key, value := range m.Iterate() {
		keys = append(keys, key)
		values = append(values, value)
	}
	if !reflect.DeepEqual(keys, []int{1, 2, 3, 4}) {
		t.Error(keys)
	}
	if !reflect.DeepEqual(values, []string{"apple", "banana", "citrus", "dragonfruit"}) {
		t.Error(values)
	}
}

func TestUnmarshalMarshalStringAny(t *testing.T) {
	dataIn := []byte(`{
  "id": "001",
  "name": "Test Object",
  "active": true,
  "settings": {
    "theme": "dark",
    "notifications": {
      "email": true,
      "sms": false,
      "push": {
        "enabled": true,
        "frequency": "daily",
        "next": {
          "first": "first",
          "second": "second"
        }
      }
    }
  },
  "users": [
    {
      "id": "user_01",
      "name": "Alice",
      "roles": [
        "admin",
        "editor"
      ],
      "preferences": {
        "language": "en",
        "timezone": "UTC",
        "dashboard": {
          "widgets": [
            "stats",
            "tasks",
            "notifications"
          ],
          "layout": "grid"
        }
      }
    },
    {
      "id": "user_02",
      "name": "Bob",
      "roles": [
        "viewer"
      ],
      "preferences": {
        "language": "fr",
        "timezone": "CET",
        "dashboard": {
          "widgets": [
            "news",
            "calendar"
          ],
          "layout": "list"
        }
      }
    }
  ],
  "logs": [
    {
      "timestamp": "2025-01-20T10:00:00Z",
      "level": "info",
      "message": "System started"
    },
    {
      "timestamp": "2025-01-20T10:05:00Z",
      "level": "warning",
      "message": "High memory usage detected"
    }
  ]
}`)
	var m OrderedMap[string, any]
	err := json.Unmarshal(dataIn, &m)
	if err != nil {
		t.Error(err)
	}
	dataOut, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(dataOut, dataIn) {
		t.Fatalf("Not equal:\n%s\n---\n%s\n---\n%+v", dataIn, dataOut, m)
	}
}

func TestMarshalUnmarshalStringStruct(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	p1 := Person{Name: "John", Age: 70}
	p2 := Person{Name: "Helen", Age: 55}
	p3 := Person{Name: "Sam", Age: 56}

	m := NewOrderedMap[string, Person]()
	m.Put("p3", p3)
	m.Put("p1", p1)
	m.Put("p2", p2)

	data, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != `{"p3":{"Name":"Sam","Age":56},"p1":{"Name":"John","Age":70},"p2":{"Name":"Helen","Age":55}}` {
		t.Fatal(string(data))
	}

	var mAfter OrderedMap[string, Person]
	err = json.Unmarshal(data, &mAfter)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(mAfter.keys, []string{"p3", "p1", "p2"}) {
		t.Fatal("Not equal", mAfter.keys)
	}

	p1After, _ := mAfter.Get("p1")
	p2After, _ := mAfter.Get("p2")
	p3After, _ := mAfter.Get("p3")

	if !reflect.DeepEqual(p1After, p1) {
		t.Fatal("Not equal", p3After)
	}

	if !reflect.DeepEqual(p2After, p2) {
		t.Fatal("Not equal", p3After)
	}

	if !reflect.DeepEqual(p3After, p3) {
		t.Fatal("Not equal", p3After)
	}
}
