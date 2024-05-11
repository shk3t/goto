package utils

func Default[T any](value T, err error) T {
	if err != nil {
		var goDefault T
		return goDefault
	}
	return value
}