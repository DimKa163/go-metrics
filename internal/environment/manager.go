package environment

import (
	"encoding/json"
	"errors"
	"flag"
	"io"
	"os"
	"reflect"
)

type Arg interface {
	Value() any
	IsDefault() bool
}

type arg struct {
	value any
}

func (a *arg) IsDefault() bool {
	if a.value == nil {
		return true
	}
	v := reflect.ValueOf(a.value)
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return true
		}
		v = v.Elem()
	}
	return v.IsZero()
}

func (a *arg) Value() any {
	return a.value
}

var args map[string]Arg = make(map[string]Arg)

var envArgs map[string]Arg = make(map[string]Arg)

func BindIntArg(name string, value int, usage string) {
	args[name] = &arg{flag.Int(name, value, usage)}
}

func BindInt64Arg(name string, value int64, usage string) {
	args[name] = &arg{flag.Int64(name, value, usage)}
}
func BindStringArg(name string, value string, usage string) {
	args[name] = &arg{flag.String(name, value, usage)}
}

func BindBooleanArg(name string, value bool, usage string) {
	args[name] = &arg{flag.Bool(name, value, usage)}
}

func BindIntEnv(name string) {
	var val *int
	ParseIntEnv(name, val)
	envArgs[name] = &arg{val}
}

func BindInt64Env(name string) {
	var val *int64
	ParseInt64Env(name, val)
	envArgs[name] = &arg{val}
}

func BindBooleanEnv(name string) {
	var val *bool
	ParseBoolEnv(name, val)
	envArgs[name] = &arg{val}
}

func BindStringEnv(name string) {
	var val *string
	if envValue := os.Getenv(name); envValue != "" {
		val = &envValue
	}
	envArgs[name] = &arg{val}
}

func Parse[Config any](config *Config) error {
	flag.Parse()
	buffer, err := getConfig()
	if err != nil {
		return err
	}
	if buffer != nil {
		err = json.Unmarshal(buffer, config)
		if err != nil {
			return err
		}
	}

	v := reflect.ValueOf(config).Elem()
	tt := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldTag := tt.Field(i).Tag
		if !field.CanSet() {
			continue
		}
		switch field.Kind() {
		case reflect.String:
			val := getString(fieldTag)
			if val != nil {
				field.SetString(*val)
			}
		case reflect.Int:
			val := getInt(fieldTag)
			if val != nil {
				field.SetInt(int64(*val))
			}
		case reflect.Int64:
			val := getInt64(fieldTag)
			if val != nil {
				field.SetInt(*val)
			}
		case reflect.Bool:
			val := getBool(fieldTag)
			if val != nil {
				field.SetBool(*val)
			}
		}
	}
	return nil
}
func getString(fieldTag reflect.StructTag) *string {
	val, ok := envArgs[fieldTag.Get("envArg")]
	if ok && !val.IsDefault() {
		return val.Value().(*string)
	}
	val, ok = args[fieldTag.Get("arg")]
	if ok && !val.IsDefault() {
		return val.Value().(*string)
	}
	return nil
}

func getInt(fieldTag reflect.StructTag) *int {
	val, ok := envArgs[fieldTag.Get("envArg")]
	if ok && !val.IsDefault() {
		return val.Value().(*int)
	}
	val, ok = args[fieldTag.Get("arg")]
	if ok && !val.IsDefault() {
		return val.Value().(*int)
	}
	return nil
}

func getInt64(fieldTag reflect.StructTag) *int64 {
	val, ok := envArgs[fieldTag.Get("envArg")]
	if ok && !val.IsDefault() {
		return val.Value().(*int64)
	}
	val, ok = args[fieldTag.Get("arg")]
	if ok && !val.IsDefault() {
		return val.Value().(*int64)
	}
	return nil
}

func getBool(fieldTag reflect.StructTag) *bool {
	val, ok := envArgs[fieldTag.Get("envArg")]
	if ok && !val.IsDefault() {
		return val.Value().(*bool)
	}
	val, ok = args[fieldTag.Get("arg")]
	if ok && !val.IsDefault() {
		return val.Value().(*bool)
	}
	return nil
}

func getConfig() ([]byte, error) {
	ar, ok := envArgs["CONFIG"]
	if !ok {
		return nil, nil
	}
	if ar.IsDefault() {
		return nil, nil
	}
	path := ar.Value().(*string)
	if path == nil {
		ar = args["c"]
		if !ok {
			return nil, nil
		}
		if ar.IsDefault() {
			return nil, nil
		}
		path = ar.Value().(*string)
	}
	if _, err := os.Stat(*path); os.IsNotExist(err) {
		return nil, errors.New("config file not exist")
	}
	file, err := os.OpenFile(*path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	buffer, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}
