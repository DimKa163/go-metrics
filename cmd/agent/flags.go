package main

import (
	"flag"
	"fmt"
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

	parseAddress(&addr)

	parseReportInterval(&reportInterval)

	parsePollInterval(&pollInterval)
}

func parseAddress(addr *string) {
	if envAddr := os.Getenv("Addr"); envAddr != "" {
		addr = &envAddr
	}
}

func parseReportInterval(reportInterval *int) {
	if envReportInterval := os.Getenv("ReportInterval"); envReportInterval != "" {
		var (
			value int
			err   error
		)
		if value, err = strconv.Atoi(envReportInterval); err != nil {
			fmt.Println(err)
			return
		}
		reportInterval = &value
	}
}

func parsePollInterval(pollInterval *int) {
	if envReportInterval := os.Getenv("PollInterval"); envReportInterval != "" {
		var (
			value int
			err   error
		)
		if value, err = strconv.Atoi(envReportInterval); err != nil {
			fmt.Println(err)
			return
		}
		pollInterval = &value
	}
}
