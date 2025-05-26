package main

import (
	"flag"
	"github.com/DimKa163/go-metrics/internal/common"
	"os"
)

var addr string
var reportInterval int
var pollInterval int

func parseFlags() {
	flag.StringVar(&addr, "a", "localhost:8081", "agent address")
	flag.IntVar(&reportInterval, "r", 10, "report interval in seconds")
	flag.IntVar(&pollInterval, "p", 2, "poll interval in seconds")
	flag.Parse()
	if envValue := os.Getenv("ADDRESS"); envValue != "" {
		addr = envValue
	}

	common.ParseIntEnv("REPORT_INTERVAL", &reportInterval)

	common.ParseIntEnv("POLL_INTERVAL", &pollInterval)
}
