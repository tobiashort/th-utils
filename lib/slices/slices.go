package slices

import "github.com/tobiashort/th-utils/lib/iter"

func Skip[A any](slice []A, skip int) []A {
	return iter.From(slice).Skip(skip).Collect()
}

func Filter[A any](slice []A, filter func(A) bool) []A {
	return iter.From(slice).Filter(filter).Collect()
}

func Count[A any](slice []A, filter func(A) bool) int {
	return len(Filter(slice, filter))
}

func Map[A, B any](slice []A, mapper func(A) B) []B {
	return iter.From(slice).Map(mapper).Collect()
}

func Reduce[A any](slice []A, initialValue A, reduce func(accumulator, value A) A) A {
	return iter.From(slice).Reduce(initialValue, reduce)
}

func ForEach[A any](slice []A, do func(A)) {
	iter.From(slice).ForEach(do)
}
