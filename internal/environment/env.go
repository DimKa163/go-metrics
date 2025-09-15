package environment

import (
	"os"
	"strconv"
)

func ParseIntEnv(name string, defValue *int) error {
	var val int
	var err error
	if envValue := os.Getenv(name); envValue != "" {
		if val, err = strconv.Atoi(envValue); err != nil {
			return err
		}
		*defValue = val
	}
	return nil
}

func ParseInt64Env(name string, defValue *int64) error {
	var val int64
	var err error
	if envValue := os.Getenv(name); envValue != "" {
		if val, err = strconv.ParseInt(envValue, 10, 64); err != nil {
			return err
		}
		*defValue = val
	}
	return nil
}

func ParseBoolEnv(name string, defValue *bool) error {
	var val bool
	var err error
	if envValue := os.Getenv(name); envValue != "" {
		if val, err = strconv.ParseBool(envValue); err != nil {
			return err
		}
		*defValue = val
	}
	return nil
}
