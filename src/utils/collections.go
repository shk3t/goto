package utils

func UniqueOnly[T comparable](data *[]T) bool {
	uniqueValues := make(map[T]struct{})
	for _, x := range *data {
		uniqueValues[x] = struct{}{}
	}
	return len(uniqueValues) == len(*data)
}

func GetAssertDefault[T any](data map[string]any, key string, defaultValue T) T {
	value, ok := data[key]
	if ok {
		return value.(T)
	} else {
		return defaultValue
	}
}