package iter

import (
	"iter"
)

type Iterator[A any] struct {
	seq iter.Seq2[int, A]
}

func From[A any](slice []A) Iterator[A] {
	return Iterator[A]{
		seq: func(yield func(int, A) bool) {
			for index, item := range slice {
				if !yield(index, item) {
					return
				}
			}
		},
	}
}

func Skip[A any](iterator Iterator[A], skip int) Iterator[A] {
	return Iterator[A]{
		seq: func(yield func(int, A) bool) {
			for index, item := range iterator.seq {
				if index >= skip {
					if !yield(index, item) {
						break
					}
				}
			}
		},
	}
}

func Filter[A any](iterator Iterator[A], filter func(A) bool) Iterator[A] {
	return Iterator[A]{
		seq: func(yield func(int, A) bool) {
			for index, item := range iterator.seq {
				if filter(item) {
					if !yield(index, item) {
						return
					}
				}
			}
		},
	}
}

func Map[A, B any](iterator Iterator[A], mapper func(A) B) Iterator[B] {
	return Iterator[B]{
		seq: func(yield func(int, B) bool) {
			for index, item := range iterator.seq {
				if !yield(index, mapper(item)) {
					return
				}
			}
		},
	}
}

func Reduce[A any](iterator Iterator[A], initialValue A, reduce func(accumulator, value A) A) A {
	accumulator := initialValue
	for _, item := range iterator.seq {
		accumulator = reduce(accumulator, item)
	}
	return accumulator
}

func ForEach[A any](iterator Iterator[A], do func(A)) {
	for _, item := range iterator.seq {
		do(item)
	}
}

func Collect[A any](iterator Iterator[A]) []A {
	var slice []A
	for _, item := range iterator.seq {
		slice = append(slice, item)
	}
	return slice
}
