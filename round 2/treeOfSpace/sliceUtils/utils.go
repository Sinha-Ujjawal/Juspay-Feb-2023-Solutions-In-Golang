package sliceUtils

func Map[T any, U any](xs []T, mapFn func(T) U) []U {
	ret := make([]U, len(xs))
	for i, x := range xs {
		ret[i] = mapFn(x)
	}
	return ret
}

func Equal[T comparable](xs []T, ys []T) bool {
	if len(xs) != len(ys) {
		return false
	}
	for i := 0; i < len(xs); i++ {
		if xs[i] != ys[i] {
			return false
		}
	}
	return true
}

func AsChannel[T any](xs []T) <-chan T {
	ch := make(chan T, len(xs))
	go func() {
		for i := 0; i < len(xs); i++ {
			ch <- xs[i]
		}
		close(ch)
	}()
	return ch
}
