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

func (iter Iterator[A]) Skip(skip int) Iterator[A] {
	return Iterator[A]{
		seq: func(yield func(int, A) bool) {
			for index, item := range iter.seq {
				if index >= skip {
					if !yield(index, item) {
						break
					}
				}
			}
		},
	}
}

func (iter Iterator[A]) Filter(filter func(A) bool) Iterator[A] {
	return Iterator[A]{
		seq: func(yield func(int, A) bool) {
			for index, item := range iter.seq {
				if filter(item) {
					if !yield(index, item) {
						return
					}
				}
			}
		},
	}
}

func (iter Iterator[A]) Map[B any](mapper func(A) B) Iterator[B] {
	return Iterator[B]{
		seq: func(yield func(int, B) bool) {
			for index, item := range iter.seq {
				if !yield(index, mapper(item)) {
					return
				}
			}
		},
	}
}

func (iter Iterator[A]) Reduce(initialValue A, reduce func(accumulator, value A) A) A {
	accumulator := initialValue
	for _, item := range iter.seq {
		accumulator = reduce(accumulator, item)
	}
	return accumulator
}

func (iter Iterator[A]) ForEach(do func(A)) {
	for _, item := range iter.seq {
		do(item)
	}
}

func (iter Iterator[A]) Collect() []A {
	var slice []A
	for _, item := range iter.seq {
		slice = append(slice, item)
	}
	return slice
}
