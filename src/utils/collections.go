package utils

func UniqueOnly[T comparable](data *[]T) bool {
	uniqueValues := make(map[T]struct{})
	for _, x := range *data {
		uniqueValues[x] = struct{}{}
	}
	return len(uniqueValues) == len(*data)
}