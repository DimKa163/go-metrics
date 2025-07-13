package collector

import (
	"context"
	"fmt"
	"github.com/DimKa163/go-metrics/internal/client"
	"github.com/DimKa163/go-metrics/internal/client/tripper"
	"github.com/DimKa163/go-metrics/internal/models"
	"github.com/DimKa163/go-metrics/internal/runtime"
	"net/http"
	"os/signal"
	"syscall"
	"time"
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

func (c *Collector) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	var count int64
	pollTicker := time.NewTicker(time.Duration(c.PollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(c.ReportInterval) * time.Second)
	values, err := runtime.ReadStat()
	if err != nil {
		return err
	}
	jobs := make(chan *models.Metric, len(values)+1)
	for i := 0; i < c.Limit; i++ {
		go c.worker(ctx, jobs)
	}
	for {
		select {
		case <-ctx.Done():
			close(jobs)
			return ctx.Err()
		case <-pollTicker.C:
			values, _ = runtime.ReadStat()
			count++
		case <-reportTicker.C:
			for k, v := range values {
				jobs <- models.CreateGauge(k, v)
			}
			jobs <- models.CreateCounter("PollCount", count)
		}
	}
}

func (c *Collector) worker(ctx context.Context, ch <-chan *models.Metric) {
	for {
		select {
		case <-ctx.Done():
			return
		case metric := <-ch:
			fmt.Println(metric)
			if metric.Type == models.CounterType {
				_ = c.UpdateCounter(metric.ID, *metric.Delta)
			} else if metric.Type == models.GaugeType {
				_ = c.UpdateGauge(metric.ID, *metric.Value)
			}
		}
	}

}
