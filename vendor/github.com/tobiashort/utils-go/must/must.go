package must

import "github.com/tobiashort/utils-go/assert"

func Do(err error) {
	assert.Nilf(err, "error: %w", err)
}

func Do2[T any](val T, err error) T {
	assert.Nilf(err, "error: %w", err)
	return val
}

func Do3[T1 any, T2 any](val1 T1, val2 T2, err error) (T1, T2) {
	assert.Nilf(err, "error: %w", err)
	return val1, val2
}
