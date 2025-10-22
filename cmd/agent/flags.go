package main

import (
	"github.com/DimKa163/go-metrics/app/collector"
	"github.com/DimKa163/go-metrics/internal/environment"
)

func ParseFlags(config *collector.Config) {
	environment.BindStringArg("a", "localhost:8082", "keeper address")
	environment.BindStringEnv("ADDRESS")
	environment.BindIntArg("r", 10, "report interval in seconds")
	environment.BindIntEnv("REPORT_INTERVAL")
	environment.BindIntArg("p", 2, "poll interval in seconds")
	environment.BindIntEnv("POLL_INTERVAL")
	environment.BindStringArg("k", "", "key")
	environment.BindStringEnv("KEY")
	environment.BindIntArg("l", 4, "rate limit")
	environment.BindIntEnv("RATE_LIMIT")
	environment.BindStringArg("c", "", "config")
	environment.BindStringEnv("CONFIG")
	environment.BindStringArg("crypto-key", "", "crypto key")
	environment.BindStringEnv("CRYPTO_KEY")
	environment.BindBooleanArg("ga", false, "grpc address")
	environment.BindBooleanEnv("USE_GRPC")
	environment.Parse(config)
}
