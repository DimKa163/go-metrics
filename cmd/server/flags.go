package main

import (
	"flag"
	"github.com/DimKa163/go-metrics/app/keeper"
	"github.com/DimKa163/go-metrics/internal/environment"
	"os"
)

func ParseFlags(config *keeper.Config) {
	flag.StringVar(&config.Addr, "a", ":8080", "keeper address")
	flag.StringVar(&config.LogLevel, "l", "info", "log level")
	flag.Int64Var(&config.StoreInterval, "i", 300, "store interval")
	flag.StringVar(&config.Path, "f", "dump", "file to store data")
	flag.BoolVar(&config.Restore, "r", true, "restore data")
	flag.StringVar(&config.DatabaseDSN, "d", "", "database connection string")
	flag.StringVar(&config.Key, "k", "secret_key", "key")
	flag.Parse()

	if envValue := os.Getenv("ADDRESS"); envValue != "" {
		config.Addr = envValue
	}

	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		config.LogLevel = envLogLevel
	}

	if databaseString := os.Getenv("DATABASE_DSN"); databaseString != "" {
		config.DatabaseDSN = databaseString
	}

	if envValue := os.Getenv("KEY"); envValue != "" {
		config.Key = envValue
	}

	environment.ParseInt64Env("STORE_INTERVAL", &config.StoreInterval)

	if envPath := os.Getenv("FILE_STORAGE_PATH"); envPath != "" {
		config.Path = envPath
	}

	environment.ParseBoolEnv("RESTORE", &config.Restore)
}
