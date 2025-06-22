package main

import (
	"context"
	"fmt"
	"github.com/DimKa163/go-metrics/internal/client"
	"math/rand"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func main() {
	parseFlags()
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	_ = run(ctx)
}

func run(ctx context.Context) error {
	var count int64
	cl := client.NewClient(fmt.Sprintf("http://%s", addr))
	pollTicker := time.NewTicker(time.Duration(pollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(reportInterval) * time.Second)
	values := make(map[string]float64)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-pollTicker.C:
			var memStats runtime.MemStats
			runtime.ReadMemStats(&memStats)
			values["Alloc"] = float64(memStats.Alloc)
			values["BuckHashSys"] = float64(memStats.BuckHashSys)
			values["Frees"] = float64(memStats.Frees)
			values["GCCPUFraction"] = float64(memStats.GCCPUFraction)
			values["GCSys"] = float64(memStats.GCSys)
			values["HeapAlloc"] = float64(memStats.HeapAlloc)
			values["HeapIdle"] = float64(memStats.HeapIdle)
			values["HeapInuse"] = float64(memStats.HeapInuse)
			values["HeapObjects"] = float64(memStats.HeapObjects)
			values["HeapReleased"] = float64(memStats.HeapReleased)
			values["HeapSys"] = float64(memStats.HeapSys)
			values["LastGC"] = float64(memStats.LastGC)
			values["MCacheInuse"] = float64(memStats.MCacheInuse)
			values["MCacheSys"] = float64(memStats.MCacheSys)
			values["MSpanSys"] = float64(memStats.MSpanSys)
			values["Mallocs"] = float64(memStats.Mallocs)
			values["NextGC"] = float64(memStats.NextGC)

			values["NumForcedGC"] = float64(memStats.NumForcedGC)
			values["NumGC"] = float64(memStats.NumGC)
			values["OtherSys"] = float64(memStats.OtherSys)

			values["PauseTotalNs"] = float64(memStats.PauseTotalNs)
			values["StackInuse"] = float64(memStats.StackInuse)
			values["StackSys"] = float64(memStats.StackSys)
			values["Sys"] = float64(memStats.Sys)

			values["TotalAlloc"] = float64(memStats.TotalAlloc)
			values["MSpanInuse"] = float64(memStats.MSpanInuse)
			values["Lookups"] = float64(memStats.Lookups)
			values["RandomValue"] = rand.Float64()
			count++
		case <-reportTicker.C:
			for k, v := range values {
				execute(cl.UpdateGauge, k, v)
			}
			execute(cl.UpdateCounter, "PollCount", count)

		}
	}
}
func execute[T float64 | int64](handler func(name string, value T) error, name string, value T) {
	err := handler(name, value)
	if err != nil {
		fmt.Println(err)
	}
}
