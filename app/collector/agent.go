// Package collector collect runtime metric
package collector

import (
	"context"
	"fmt"
	"github.com/DimKa163/go-metrics/internal/crypto"
	"github.com/DimKa163/go-metrics/internal/gc/gclient"
	"github.com/DimKa163/go-metrics/internal/mhttp/contracts"
	"github.com/DimKa163/go-metrics/internal/mhttp/hclient"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/DimKa163/go-metrics/internal/client"
	"github.com/DimKa163/go-metrics/internal/models"
	"github.com/DimKa163/go-metrics/internal/runtime"
)

type Collector struct {
	*Config
	wg sync.WaitGroup
	client.MetricClient
	jobs chan *contracts.Metric
}

func NewCollector(conf *Config) (*Collector, error) {
	cl, err := CreateClient(conf)
	if err != nil {
		return nil, err
	}
	return &Collector{Config: conf, MetricClient: cl}, nil
}

// Run worker
func (c *Collector) Run(buildVersion string, buildDate string, buildCommit string) error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()
	var count int64
	values := make(map[string]float64)
	c.jobs = make(chan *contracts.Metric, c.Limit*4)
	var err error
	for i := 0; i < c.Limit; i++ {
		go c.worker(ctx)
	}
	pollTicker := time.NewTicker(time.Duration(c.PollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(c.ReportInterval) * time.Second)
	printBuildInfo(buildVersion, buildDate, buildCommit)
	for {
		select {
		case <-ctx.Done():
			c.wg.Wait()
			close(c.jobs)
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
				c.wg.Add(1)
				c.jobs <- contracts.CreateGauge(k, v)
			}
			c.wg.Add(1)
			c.jobs <- contracts.CreateCounter("PollCount", count)

		}
	}
}

func (c *Collector) worker(ctx context.Context) {
	for metric := range c.jobs {
		if metric == nil {
			continue
		}
		fmt.Println(metric)
		if metric.Type == models.CounterType {
			if err := c.UpdateCounter(ctx, metric.ID, metric.Delta); err != nil {
				fmt.Println(err)
				continue
			}
		} else if metric.Type == models.GaugeType {
			if err := c.UpdateGauge(ctx, metric.ID, metric.Value); err != nil {
				fmt.Println(err)
				continue
			}
		}
		c.wg.Done()
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

func CreateClient(conf *Config) (client.MetricClient, error) {
	if conf.UseGrpc {
		return CreateGRPCMetricClient(conf.Addr,
			gclient.UnaryRetryInterceptor(),
			gclient.UnaryLoggingInterceptor(),
			gclient.UnaryIdentifyInterceptor())
	}

	tripperFc := []hclient.RequestHandler{
		hclient.UseRetryHandler(),
		hclient.UseGzipHandler(),
		hclient.UseIdentifyHandler(),
	}
	if conf.Key != "" {
		tripperFc = append(tripperFc, hclient.UseHashHandler(conf.Key))
	}

	if conf.PublicKeyFilePath != "" {
		encrypter, err := crypto.NewEncrypter(conf.PublicKeyFilePath)
		if err != nil {
			return nil, err
		}
		tripperFc = append(tripperFc, hclient.UseCryptoHandler(encrypter))
	}
	return CreateHTTPMetricClient(hclient.HTTP, conf.Addr, tripperFc...), nil
}

func CreateGRPCMetricClient(addr string, interceptors ...grpc.UnaryClientInterceptor) (*gclient.GrpcMetricClient, error) {
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(interceptors...))
	if err != nil {
		return nil, err
	}
	return gclient.NewGRPCMetricClient(conn), nil
}

func CreateHTTPMetricClient(protocol, addr string, handlers ...hclient.RequestHandler) *hclient.HTTPMetricClient {
	return hclient.NewClient(fmt.Sprintf("%s://%s", protocol, addr), handlers...)
}
