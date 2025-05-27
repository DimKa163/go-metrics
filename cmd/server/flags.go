package main

import (
	"flag"
	"os"
)

var addr string

func parseFlag() {
	flag.StringVar(&addr, "a", ":8080", "server address")
	flag.Parse()

	if envValue := os.Getenv("ADDRESS"); envValue != "" {
		addr = envValue
	}
}
