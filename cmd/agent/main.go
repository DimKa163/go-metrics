package main

import (
	"fmt"
	"github.com/DimKa163/go-metrics/internal/client"
	"runtime"
	"time"
)

func main() {
	count := int64(0)
	cl := client.NemClient("http://localhost:8080/")
	for {
		memStats := &runtime.MemStats{}

		runtime.ReadMemStats(memStats)

		if count%5 == 0 {
			report(cl, memStats, count)
		}

		time.Sleep(2 * time.Second)
		count++
	}
}

func report(cl client.MetricClient, memStats *runtime.MemStats, count int64) {
	execute(cl.UpdateGauge, "Alloc", float64(memStats.Alloc))
	execute(cl.UpdateGauge, "BuckHashSys", float64(memStats.BuckHashSys))
	execute(cl.UpdateGauge, "Frees", float64(memStats.Frees))
	execute(cl.UpdateGauge, "GCCPUFraction", memStats.GCCPUFraction)
	execute(cl.UpdateGauge, "GCSys", float64(memStats.GCSys))
	execute(cl.UpdateGauge, "HeapAlloc", float64(memStats.HeapAlloc))
	execute(cl.UpdateGauge, "HeapIdle", float64(memStats.HeapIdle))
	execute(cl.UpdateGauge, "HeapInuse", float64(memStats.HeapInuse))
	execute(cl.UpdateGauge, "HeapObjects", float64(memStats.HeapObjects))
	execute(cl.UpdateGauge, "HeapReleased", float64(memStats.HeapReleased))
	execute(cl.UpdateGauge, "HeapSys", float64(memStats.HeapSys))
	execute(cl.UpdateGauge, "LastGC", float64(memStats.LastGC))
	execute(cl.UpdateGauge, "MCacheInuse", float64(memStats.Lookups))
	execute(cl.UpdateGauge, "MCacheInuse", float64(memStats.MCacheInuse))
	execute(cl.UpdateGauge, "MCacheSys", float64(memStats.MCacheSys))
	execute(cl.UpdateGauge, "MSpanSys", float64(memStats.MSpanSys))
	execute(cl.UpdateGauge, "Mallocs", float64(memStats.Mallocs))
	execute(cl.UpdateGauge, "NumForcedGC", float64(memStats.NextGC))
	execute(cl.UpdateGauge, "NumForcedGC", float64(memStats.NumForcedGC))
	execute(cl.UpdateGauge, "NumGC", float64(memStats.NumGC))
	execute(cl.UpdateGauge, "OtherSys", float64(memStats.OtherSys))
	execute(cl.UpdateGauge, "StackInuse", float64(memStats.PauseTotalNs))
	execute(cl.UpdateGauge, "StackInuse", float64(memStats.StackInuse))
	execute(cl.UpdateGauge, "StackSys", float64(memStats.StackSys))
	execute(cl.UpdateGauge, "Sys", float64(memStats.Sys))
	execute(cl.UpdateGauge, "TotalAlloc", float64(memStats.TotalAlloc))
	execute(cl.UpdateCounter, "PollCount", count)
}

func execute[T float64 | int64](handler func(name string, value T) error, name string, value T) {
	err := handler(name, value)
	if err != nil {
		fmt.Println(err)
	}
}
