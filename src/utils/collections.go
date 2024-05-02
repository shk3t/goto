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

func MapKeys[K comparable, V any](data map[K]V) []K {
	keys := make([]K, len(data))
	i := 0
	for k := range data {
		keys[i] = k
		i++
	}
	return keys
}

func MapValues[K comparable, V any](data map[K]V) []V {
	values := make([]V, len(data))
	i := 0
	for _, v := range data {
		values[i] = v
		i++
	}
	return values
}