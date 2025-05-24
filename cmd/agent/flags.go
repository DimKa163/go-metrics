package main

import "flag"

var addr string
var reportInterval int
var pollInterval int

func parseFlags() {
	flag.StringVar(&addr, "a", "localhost:8080", "agent address")
	flag.IntVar(&reportInterval, "r", 10, "report interval in seconds")
	flag.IntVar(&pollInterval, "p", 2, "poll interval in seconds")
	flag.Parse()
}
