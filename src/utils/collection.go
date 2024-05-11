package utils

import "errors"

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
	diff := []string{}
	for _, x := range left {
		if _, found := rightMap[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

func GetAssertDefault[T any](data map[string]any, key string, defaultValue T) T {
	value, ok := data[key]
	if ok {
		value, ok := value.(T)
		if ok {
			return value
		}
	}
	return defaultValue
}

func GetAssertError[T any](data map[string]any, key string, errorContext string) (T, error) {
	var goDefault T
	errorMessage := ""

	value, ok := data[key]
	if ok {
		value, ok := value.(T)
		if ok {
			return value, nil
		}

		errorMessage = "`" + key + "` has bad format"
	} else {
		errorMessage = "`" + key + "` is not specified"
	}

	if errorContext != "" {
		errorMessage = errorContext + ": " + errorMessage
	}
	return goDefault, errors.New(errorMessage)
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