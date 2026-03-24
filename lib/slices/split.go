package slices

func Split[T any](a []T, fn func(a T) bool) [][]T {
	r := make([][]T, 0)
	if len(a) == 0 {
		return r
	}
	s := make([]T, 0)
	for _, v := range a {
		if fn(v) {
			r = append(r, s)
			s = make([]T, 0)
			continue
		}
		s = append(s, v)
	}
	r = append(r, s)
	return r
}
