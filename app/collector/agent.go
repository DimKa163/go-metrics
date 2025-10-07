// Package collector collect runtime metric
package collector

import (
	"context"
	"fmt"
	"github.com/DimKa163/go-metrics/internal/crypto"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/DimKa163/go-metrics/internal/client"
	"github.com/DimKa163/go-metrics/internal/client/tripper"
	"github.com/DimKa163/go-metrics/internal/models"
	"github.com/DimKa163/go-metrics/internal/runtime"
)

type Collector struct {
	*Config
	wg sync.WaitGroup
	client.MetricClient
}

func NewCollector(conf *Config) (*Collector, error) {
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

	if conf.PublicKeyFilePath != "" {
		encrypter, err := crypto.NewEncrypter(conf.PublicKeyFilePath)
		if err != nil {
			return nil, err
		}
		tripperFc = append(tripperFc, func(transport http.RoundTripper) http.RoundTripper {
			return tripper.NewCryptoTripper(transport, encrypter)
		})
	}
	return &Collector{Config: conf, MetricClient: client.NewClient(fmt.Sprintf("http://%s", conf.Addr), tripperFc)}, nil
}

// Run worker
func (c *Collector) Run(buildVersion string, buildDate string, buildCommit string) error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	var count int64
	values := make(map[string]float64)
	jobs := make(chan *models.Metric, c.Limit*4)
	var err error
	for i := 0; i < c.Limit; i++ {
		c.wg.Add(1)
		go c.worker(ctx, jobs)
	}
	pollTicker := time.NewTicker(time.Duration(c.PollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(c.ReportInterval) * time.Second)
	printBuildInfo(buildVersion, buildDate, buildCommit)
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
			for k, v := range values {

				jobs <- models.CreateGauge(k, v)
			}
			jobs <- models.CreateCounter("PollCount", count)
			fmt.Println("Start waiting for jobs to finish")
			fmt.Println("End")
			close(jobs)
		}
	}
}

func (c *Collector) worker(ctx context.Context, ch <-chan *models.Metric) {
	defer c.wg.Done()
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

func printBuildInfo(buildVersion string, buildDate string, buildCommit string) {
	fmt.Printf("Build version: %s\n", ifNan(buildVersion))
	fmt.Printf("Build date: %s\n", ifNan(buildDate))
	fmt.Printf("Build commit: %s\n", ifNan(buildCommit))
}

func ifNan(value string) string {
	if value == "" {
		return "N/A"
	}
	return value
}
