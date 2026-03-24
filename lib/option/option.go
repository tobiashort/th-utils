package option

type Option[T any] struct {
	Val  T
	None bool
}

func Some[T any](val T) Option[T] {
	return Option[T]{
		Val:  val,
		None: false,
	}
}

func None[T any]() Option[T] {
	return Option[T]{
		None: true,
	}
}
