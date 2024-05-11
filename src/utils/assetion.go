package utils

import "errors"

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