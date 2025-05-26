package common

import (
	"os"
	"strconv"
)

func ParseStringEnv(name string, defValue *string) {
	if envValue := os.Getenv(name); envValue != "" {
		defValue = &envValue
	}
}
func ParseIntEnv(name string, defValue *int) {
	if envValue := os.Getenv(name); envValue != "" {
		if value, err := strconv.Atoi(envValue); err == nil {
			*defValue = value
		}
	}
}
