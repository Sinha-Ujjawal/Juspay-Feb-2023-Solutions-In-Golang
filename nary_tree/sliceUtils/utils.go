package sliceUtils

import "math"

func Map[T any, U any](xs []T, mapFn func(T) U) []U {
	ret := []U{}
	for _, x := range xs {
		ret = append(ret, mapFn(x))
	}
	return ret
}

func ZipWith[T any, U any, V any](ts []T, us []U, zipFn func(T, U) V) []V {
	l := uint(math.Min(float64(len(ts)), float64(len(us))))
	ret := []V{}
	for i := uint(0); i < l; i++ {
		ret = append(ret, zipFn(ts[i], us[i]))
	}
	return ret
}

type Pair[T any, U any] struct {
	First  T
	Second U
}

func Zip[T any, U any](ts []T, us []U) []Pair[T, U] {
	return ZipWith(ts, us, func(t T, u U) Pair[T, U] {
		return Pair[T, U]{First: t, Second: u}
	})
}

func And[T any](ts []T, checkFn func(T) bool) bool {
	for _, t := range ts {
		if !checkFn(t) {
			return false
		}
	}
	return true
}

func Or[T any](ts []T, checkFn func(T) bool) bool {
	for _, t := range ts {
		if checkFn(t) {
			return true
		}
	}
	return false
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

func AsChannel[T any](xs []T) <-chan *T {
	ch := make(chan *T, len(xs))
	go func() {
		for i := 0; i < len(xs); i++ {
			ch <- &xs[i]
		}
		close(ch)
	}()
	return ch
}
