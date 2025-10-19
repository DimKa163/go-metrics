package gclient

import (
	"context"
	"github.com/DimKa163/go-metrics/internal/gc/proto"
	"github.com/DimKa163/go-metrics/internal/mhttp/contracts"
	"github.com/DimKa163/go-metrics/internal/models"
	"google.golang.org/grpc"
)

type GrpcMetricClient struct {
	proto.MetricsClient
}

func (g *GrpcMetricClient) UpdateGauge(ctx context.Context, name string, value float64) error {
	var metricRequest proto.UpdateMetric
	var metric proto.Metric
	metric.SetId(name)
	metric.SetType(proto.MetricType_GAUGE)
	metric.SetValue(value)
	metricRequest.SetMetric(&metric)
	_, err := g.Update(ctx, &metricRequest)
	return err
}

func (g *GrpcMetricClient) UpdateCounter(ctx context.Context, name string, value int64) error {
	var metricRequest proto.UpdateMetric
	var metric proto.Metric
	metric.SetId(name)
	metric.SetType(proto.MetricType_COUNTER)
	metric.SetDelta(value)
	metricRequest.SetMetric(&metric)
	_, err := g.Update(ctx, &metricRequest)
	return err
}

func (g *GrpcMetricClient) BatchUpdate(ctx context.Context, metrics []*contracts.Metric) error {
	var metricRequest proto.BatchUpdateMetric
	in := make([]*proto.Metric, len(metrics))
	for i, m := range metrics {
		var metric proto.Metric
		in[i] = &metric
		metric.SetId(m.ID)
		switch m.Type {
		case models.GaugeType:
			metric.SetType(proto.MetricType_GAUGE)
			metric.SetValue(m.Value)
		case models.CounterType:
			metric.SetType(proto.MetricType_COUNTER)
			metric.SetDelta(m.Delta)
		}
	}
	metricRequest.SetMetric(in)
	_, err := g.MetricsClient.BatchUpdate(ctx, &metricRequest)
	return err
}

func NewGRPCMetricClient(client *grpc.ClientConn) *GrpcMetricClient {
	return &GrpcMetricClient{
		MetricsClient: proto.NewMetricsClient(client),
	}
}
