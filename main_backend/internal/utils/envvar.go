package utils

import (
	"fmt"
	"os"
)

func GetEnvAssert(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("%s environment variable is not set", key))
	}
	return value
}
