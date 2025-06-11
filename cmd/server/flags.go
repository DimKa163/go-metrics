package main

import (
	"flag"
	"os"
)

var (
	addr     string
	logLevel string
)

func parseFlag() {
	flag.StringVar(&addr, "a", ":8080", "server address")
	flag.StringVar(&logLevel, "l", "info", "log level")
	flag.Parse()

	if envValue := os.Getenv("ADDRESS"); envValue != "" {
		addr = envValue
	}

	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		logLevel = envLogLevel
	}
}
