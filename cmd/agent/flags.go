package main

import (
	"flag"
	"os"
	"strconv"
)

var addr string
var reportInterval int
var pollInterval int

func parseFlags() {
	flag.StringVar(&addr, "a", "localhost:8080", "agent address")
	flag.IntVar(&reportInterval, "r", 10, "report interval in seconds")
	flag.IntVar(&pollInterval, "p", 2, "poll interval in seconds")
	flag.Parse()

	if envAddress := os.Getenv("ADDRESS"); envAddress != "" {
		addr = envAddress
	}

	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		var (
			value int
			err   error
		)
		if value, err = strconv.Atoi(envReportInterval); err == nil {
			reportInterval = value
		}

	}

	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		var (
			value int
			err   error
		)
		if value, err = strconv.Atoi(envPollInterval); err == nil {
			pollInterval = value
		}

	}
}
