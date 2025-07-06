package collector

import (
	"context"
	"fmt"
	"github.com/DimKa163/go-metrics/internal/client"
	"github.com/DimKa163/go-metrics/internal/models"
	"math/rand"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

type Collector struct {
	*Config
	client client.MetricClient
}

func NewCollector(conf *Config) *Collector {
	return &Collector{conf, client.NewClient(conf.Addr)}
}

func (s *Collector) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	var count int64
	cl := client.NewClient(fmt.Sprintf("http://%s", s.Addr))
	pollTicker := time.NewTicker(time.Duration(s.PollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(s.ReportInterval) * time.Second)
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
			var metrics []*models.Metric
			for k, v := range values {
				metrics = append(metrics, models.CreateGauge(k, v))
			}
			metrics = append(metrics, models.CreateCounter("PollCount", count))
			if err := cl.BatchUpdate(metrics); err != nil {
				fmt.Println(err)
			}

		}
	}
}
