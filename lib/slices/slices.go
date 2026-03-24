package slices

import "github.com/tobiashort/th-utils/lib/iter"

func Skip[A any](slice []A, skip int) []A {
	iterator := iter.From(slice)
	iterator = iter.Skip(iterator, skip)
	return iter.Collect(iterator)
}

func Filter[A any](slice []A, filter func(A) bool) []A {
	iterator := iter.From(slice)
	iterator = iter.Filter(iterator, filter)
	return iter.Collect(iterator)
}

func Map[A, B any](slice []A, mapper func(A) B) []B {
	iterator := iter.From(slice)
	iterator2 := iter.Map(iterator, mapper)
	return iter.Collect(iterator2)
}

func Reduce[A any](slice []A, initialValue A, reduce func(accumulator, value A) A) A {
	iterator := iter.From(slice)
	return iter.Reduce(iterator, initialValue, reduce)
}

func ForEach[A any](slice []A, do func(A)) {
	iterator := iter.From(slice)
	iter.ForEach(iterator, do)
}
