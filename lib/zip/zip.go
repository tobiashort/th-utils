package zip

import "reflect"

func Zip[T any](strct T, slices ...any) []T {
	capacity := 0
	for _, it := range slices {
		length := reflect.ValueOf(it).Len()
		if length > capacity {
			capacity = length
		}
	}

	ret := make([]T, capacity)

	for i, it := range slices {
		slice := reflect.ValueOf(it)
		for i2 := 0; i2 < slice.Len(); i2++ {
			it2 := slice.Index(i2)
			val := &ret[i2]
			reflect.ValueOf(val).Elem().Field(i).Set(it2)
		}
	}

	return ret
}
