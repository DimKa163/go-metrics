package client

import (
	"context"
	"github.com/DimKa163/go-metrics/internal/mhttp/contracts"
)

type MetricClient interface {
	UpdateGauge(ctx context.Context, name string, value float64) error
	UpdateCounter(ctx context.Context, name string, value int64) error

	BatchUpdate(ctx context.Context, metrics []*contracts.Metric) error
}
