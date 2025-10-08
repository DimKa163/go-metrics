package main

import (
	"github.com/DimKa163/go-metrics/app/keeper"
	"github.com/DimKa163/go-metrics/internal/environment"
)

func ParseFlags(config *keeper.Config) error {
	environment.BindBooleanArg("r", true, "restore data")
	environment.BindBooleanEnv("RESTORE")
	environment.BindStringArg("f", "dump", "file to store data")
	environment.BindStringEnv("FILE_STORAGE_PATH")
	environment.BindStringArg("l", "info", "log level")
	environment.BindStringEnv("LOG_LEVEL")
	environment.BindInt64Arg("i", 300, "store interval in seconds")
	environment.BindInt64Env("STORE_INTERVAL")
	environment.BindStringArg("a", ":8080", "keeper address")
	environment.BindStringEnv("ADDRESS")
	environment.BindStringArg("d", "", "keeper database")
	environment.BindStringEnv("DATABASE_DSN")
	environment.BindStringArg("k", "", "keeper key")
	environment.BindStringEnv("KEY")
	environment.BindStringArg("c", "", "config")
	environment.BindStringEnv("CONFIG")
	environment.BindStringArg("crypto-key", "", "crypto key")
	environment.BindStringEnv("CRYPTO_KEY")
	environment.Parse(config)
	return nil
}
