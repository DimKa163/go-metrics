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
	//var configPath string
	//flag.StringVar(config.Addr, "a", ":8080", "keeper address")
	//flag.StringVar(config.LogLevel, "l", "info", "log level")
	//flag.Int64Var(config.StoreInterval, "i", 300, "store interval")
	//flag.StringVar(config.Path, "f", "dump", "file to store data")
	//flag.BoolVar(config.Restore, "r", true, "restore data")
	//flag.StringVar(config.DatabaseDSN, "d", "", "database connection string")
	//flag.StringVar(config.Key, "k", "", "key")
	//flag.StringVar(&configPath, "config", "", "key")
	//flag.Parse()
	//
	//if envValue := os.Getenv("ADDRESS"); envValue != "" {
	//	config.Addr = &envValue
	//}
	//
	//if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
	//	config.LogLevel = &envLogLevel
	//}
	//
	//if databaseString := os.Getenv("DATABASE_DSN"); databaseString != "" {
	//	config.DatabaseDSN = &databaseString
	//}
	//
	//if envValue := os.Getenv("KEY"); envValue != "" {
	//	config.Key = &envValue
	//}
	//
	//environment.ParseInt64Env("STORE_INTERVAL", config.StoreInterval)
	//
	//if envPath := os.Getenv("FILE_STORAGE_PATH"); envPath != "" {
	//	config.Path = &envPath
	//}
	//
	//environment.ParseBoolEnv("RESTORE", config.Restore)
	//
	//if envValue := os.Getenv("CONFIG"); envValue != "" {
	//	configPath = envValue
	//}
	//
	//if configPath != "" {
	//	return environment.LoadConfig(configPath, func(cfg environment.Configuration) error {
	//		if config.Addr == nil {
	//			add, err := cfg.GetString("address")
	//			if err != nil {
	//				return err
	//			}
	//			config.Addr = &add
	//		}
	//		if config.Restore == nil {
	//			rest, err := cfg.GetBool("restore")
	//			if err != nil {
	//				return err
	//			}
	//			config.Restore = &rest
	//		}
	//		return nil
	//	})
	//}
	return nil
}
