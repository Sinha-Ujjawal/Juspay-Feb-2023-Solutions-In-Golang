package chanUtils

func AsSlice[T any](ch <-chan T) []T {
	var ret []T
	for t := range ch {
		ret = append(ret, t)
	}
	return ret
}
