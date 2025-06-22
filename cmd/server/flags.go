package main

import (
	"flag"
	"github.com/DimKa163/go-metrics/internal/common"
	"os"
)

func getServerConfig() ServerConfig {
	serverConfig := ServerConfig{}
	flag.StringVar(&serverConfig.Addr, "a", ":8080", "server address")
	flag.StringVar(&serverConfig.LogLevel, "l", "info", "log level")
	flag.Int64Var(&serverConfig.StoreInterval, "i", 300, "store interval")
	flag.StringVar(&serverConfig.Path, "f", "dump", "file to store data")
	flag.BoolVar(&serverConfig.Restore, "r", true, "restore data")
	flag.Parse()

	if envValue := os.Getenv("ADDRESS"); envValue != "" {
		serverConfig.Addr = envValue
	}

	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		serverConfig.LogLevel = envLogLevel
	}

	common.ParseInt64Env("STORE_INTERVAL", &serverConfig.StoreInterval)

	if envPath := os.Getenv("FILE_STORAGE_PATH"); envPath != "" {
		serverConfig.Path = envPath
	}

	common.ParseBoolEnv("RESTORE", &serverConfig.Restore)
	return serverConfig
}
