package env

import (
	"errors"
	"os"
)

func Get(key string, defaultValue string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		if defaultValue == "" {
			return "", errors.New("cannot find a thing")
		}
		return defaultValue, nil
	}

	return value, nil
}
