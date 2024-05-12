package utils

func Insert(data []any, index int, value any) []any {
	if len(data) == index {
		return append(data, value)
	}
	data = append(data[:index+1], data[index:]...)
	data[index] = value
	return data
}

func UniqueOnly[T comparable](data *[]T) bool {
	uniqueValues := map[T]struct{}{}
	for _, x := range *data {
		uniqueValues[x] = struct{}{}
	}
	return len(uniqueValues) == len(*data)
}

func Difference(left []string, right []string) []string {
	rightMap := make(map[string]struct{}, len(right))
	for _, x := range right {
		rightMap[x] = struct{}{}
	}
	result := []string{}
	for _, x := range left {
		if _, found := rightMap[x]; !found {
			result = append(result, x)
		}
	}
	return result
}

func Intersection[T comparable](left []T, right []T) []T {
	result := []T{}
	leftMap := make(map[T]struct{})
	for _, v := range left {
		leftMap[v] = struct{}{}
	}
	for _, v := range right {
		if _, ok := leftMap[v]; ok {
			result = append(result, v)
		}
	}
	return result
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