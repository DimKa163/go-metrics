package environment

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseIntEnv(t *testing.T) {
	_ = os.Setenv("TEST_INT", "42")
	defer os.Unsetenv("TEST_INT")
	var err error
	var val1 int
	err = ParseIntEnv("TEST_INT", &val1)
	assert.Equal(t, val1, 42)
	assert.NoError(t, err)

	var val2 int
	err = ParseIntEnv("NOT_EXIST", &val2)
	assert.Equal(t, val2, 0)
	assert.NoError(t, err)

	_ = os.Setenv("TEST_INT", "oops")
	var val3 int
	err = ParseIntEnv("TEST_INT", &val3)
	assert.Equal(t, val3, 0)
	assert.Error(t, err)
}

func TestParseInt64Env(t *testing.T) {
	_ = os.Setenv("TEST_INT64", "922337203685477580")
	defer os.Unsetenv("TEST_INT64")
	var err error

	var val1 int64
	err = ParseInt64Env("TEST_INT64", &val1)
	assert.Equal(t, val1, int64(922337203685477580))
	assert.NoError(t, err)

	var val2 int64
	err = ParseInt64Env("NOT_EXISTS", &val2)
	assert.Equal(t, val2, int64(0))
	assert.NoError(t, err)

	_ = os.Setenv("TEST_INT64", "bad")
	var val3 int64
	err = ParseInt64Env("TEST_INT64", &val3)
	assert.Equal(t, val3, int64(0))
	assert.Error(t, err)
}

func TestParseBoolEnv(t *testing.T) {
	_ = os.Setenv("TEST_BOOL", "true")
	defer os.Unsetenv("TEST_BOOL")
	var err error
	var val1 bool
	err = ParseBoolEnv("TEST_BOOL", &val1)
	assert.Equal(t, true, val1, "expected true, got false")
	assert.NoError(t, err)
	var val2 bool
	err = ParseBoolEnv("NOT_EXISTS", &val2)
	assert.Equal(t, false, val2)
	assert.NoError(t, err)
	_ = os.Setenv("TEST_BOOL", "notabool")
	var val3 bool
	err = ParseBoolEnv("TEST_BOOL", &val3)
	assert.Equal(t, false, val3, "expected false, got true")
	assert.Error(t, err)
}
