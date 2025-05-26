package common

import (
	"os"
	"strconv"
)

func ParseIntEnv(name string, defValue *int) {
	if envValue := os.Getenv(name); envValue != "" {
		if value, err := strconv.Atoi(envValue); err == nil {
			*defValue = value
		}
	}
}
