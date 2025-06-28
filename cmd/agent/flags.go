package main

import (
	"flag"
	"github.com/DimKa163/go-metrics/app/collector"
	"github.com/DimKa163/go-metrics/internal/environment"
	"os"
)

func ParseFlags(config *collector.Config) {
	flag.StringVar(&config.Addr, "a", "localhost:8080", "agent address")
	flag.IntVar(&config.ReportInterval, "r", 10, "report interval in seconds")
	flag.IntVar(&config.PollInterval, "p", 2, "poll interval in seconds")
	flag.Parse()
	if envValue := os.Getenv("ADDRESS"); envValue != "" {
		config.Addr = envValue
	}

	environment.ParseIntEnv("REPORT_INTERVAL", &config.ReportInterval)

	environment.ParseIntEnv("POLL_INTERVAL", &config.PollInterval)
}
