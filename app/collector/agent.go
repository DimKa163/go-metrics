// Package collector collect runtime metric
package collector

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/DimKa163/go-metrics/internal/client"
	"github.com/DimKa163/go-metrics/internal/client/tripper"
	"github.com/DimKa163/go-metrics/internal/models"
	"github.com/DimKa163/go-metrics/internal/runtime"
)

type Collector struct {
	*Config
	client.MetricClient
}

func NewCollector(conf *Config) *Collector {
	tripperFc := []func(transport http.RoundTripper) http.RoundTripper{
		func(transport http.RoundTripper) http.RoundTripper {
			return tripper.NewRetryRoundTripper(transport)
		},
		func(transport http.RoundTripper) http.RoundTripper {
			return tripper.NewGzip(transport)
		},
	}
	if conf.Key != "" {
		tripperFc = append(tripperFc, func(transport http.RoundTripper) http.RoundTripper {
			return tripper.NewHashTripper(transport, conf.Key)
		})
	}
	return &Collector{conf, client.NewClient(fmt.Sprintf("http://%s", conf.Addr), tripperFc)}
}

// Run worker
func (c *Collector) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	var count int64
	values := make(map[string]float64)
	var err error
	pollTicker := time.NewTicker(time.Duration(c.PollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(c.ReportInterval) * time.Second)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-pollTicker.C:
			err = runtime.ReadMemoryStats(values)
			if err != nil {
				fmt.Printf("Error reading stats: %v\n", err)
				continue
			}
			err = runtime.ReadCPUStats(values)
			if err != nil {
				fmt.Printf("Error reading stats: %v\n", err)
				continue
			}
			count++
		case <-reportTicker.C:
			jobs := make(chan *models.Metric, len(values))
			for i := 0; i < c.Limit; i++ {
				go c.worker(ctx, jobs)
			}
			for k, v := range values {
				jobs <- models.CreateGauge(k, v)
			}
			jobs <- models.CreateCounter("PollCount", count)
			close(jobs)
		}
	}
}

func (c *Collector) worker(ctx context.Context, ch <-chan *models.Metric) {
	for {
		select {
		case <-ctx.Done():
			return
		case metric := <-ch:
			if metric == nil {
				continue
			}
			fmt.Println(metric)
			if metric.Type == models.CounterType {
				if err := c.UpdateCounter(metric.ID, *metric.Delta); err != nil {
					fmt.Println(err)
					continue
				}
			} else if metric.Type == models.GaugeType {
				if err := c.UpdateGauge(metric.ID, *metric.Value); err != nil {
					fmt.Println(err)
					continue
				}
			}
		}
	}

}
