package chanUtils

import "sync"

func MapPar[In any, Out any](in <-chan In, mapFn func(i In, clb func(Out))) <-chan Out {
	out := make(chan Out)
	go func() {
		wg := sync.WaitGroup{}
		for i := range in {
			go func(i In) {
				mapFn(
					i,
					func(o Out) {
						out <- o
						wg.Done()
					},
				)
			}(i)
			wg.Add(1)
		}
		wg.Wait()
		close(out)
	}()
	return out
}

func Map[In any, Out any](in <-chan In, mapFn func(In) Out) <-chan Out {
	out := make(chan Out)
	go func() {
		for i := range in {
			out <- mapFn(i)
		}
		close(out)
	}()
	return out
}

func Take[T any](ch <-chan T, n uint) <-chan T {
	out := make(chan T)
	go func() {
		for i := uint(0); i < n; i++ {
			out <- <-ch
		}
		close(out)
	}()
	return out
}

func TakeAsList[T any](ch <-chan T, n uint) []T {
	var ret []T
	for i := uint(0); i < n; i++ {
		ret = append(ret, <-ch)
	}
	return ret
}

func AsSlice[T any](ch <-chan T) []T {
	var ret []T
	for t := range ch {
		ret = append(ret, t)
	}
	return ret
}
