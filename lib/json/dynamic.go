package json

import (
	. "encoding/json/v2"
	"math"

	"github.com/tobiashort/th-utils/lib/must"
)

type DynamicJson struct {
	payload any
}

func Dynamic(data []byte) *DynamicJson {
	json := new(DynamicJson{})
	must.Do(Unmarshal(data, &json.payload))
	return json
}

func (json *DynamicJson) Value(key string) *DynamicJson {
	return new(DynamicJson{payload: json.payload.(map[string]any)[key]})
}

func (json *DynamicJson) Index(i int) *DynamicJson {
	return new(DynamicJson{payload: json.payload.([]any)[i]})
}

func (json *DynamicJson) String() string {
	return json.payload.(string)
}

func (json *DynamicJson) StringSlice() []string {
	a := json.payload.([]any)
	s := make([]string, len(a))
	for i, v := range a {
		s[i] = v.(string)
	}
	return s
}

func (json *DynamicJson) Int() int {
	f := json.payload.(float64)
	if math.Trunc(f) != f {
		panic("not an integer")
	}
	return int(f)
}

func (json *DynamicJson) IntSlice() []int {
	a := json.payload.([]any)
	s := make([]int, len(a))
	for i, v := range a {
		s[i] = v.(int)
	}
	return s
}

func (json *DynamicJson) Float64() float64 {
	return json.payload.(float64)
}

func (json *DynamicJson) Float64Slice() []float64 {
	a := json.payload.([]any)
	s := make([]float64, len(a))
	for i, v := range a {
		s[i] = v.(float64)
	}
	return s
}

func (json *DynamicJson) Bool() bool {
	return json.payload.(bool)
}

func (json *DynamicJson) BoolSlice() []bool {
	a := json.payload.([]any)
	s := make([]bool, len(a))
	for i, v := range a {
		s[i] = v.(bool)
	}
	return s
}
