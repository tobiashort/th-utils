package must

func Do(err error) {
	if err != nil {
		panic(err)
	}
}

func Do2[T any](val T, err error) T {
	Do(err)
	return val
}

func Do3[T1 any, T2 any](val1 T1, val2 T2, err error) (T1, T2) {
	Do(err)
	return val1, val2
}
