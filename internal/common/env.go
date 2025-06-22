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

func ParseInt64Env(name string, defValue *int64) {
	if envValue := os.Getenv(name); envValue != "" {
		if value, err := strconv.ParseInt(envValue, 10, 64); err == nil {
			*defValue = value
		}
	}
}

func ParseBoolEnv(name string, defValue *bool) {
	if envValue := os.Getenv(name); envValue != "" {
		if value, err := strconv.ParseBool(envValue); err == nil {
			*defValue = value
		}
	}
}
