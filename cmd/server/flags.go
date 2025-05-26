package main

import (
	"flag"
	"github.com/DimKa163/go-metrics/internal/common"
)

var addr string

func parseFlag() {
	flag.StringVar(&addr, "a", ":8080", "server address")
	flag.Parse()

	common.ParseStringEnv("ADDRESS", &addr)
}
