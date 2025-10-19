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
	_, err := g.Update(ctx, &proto.UpdateMetric{
		Metric: &proto.Metric{
			Id:    name,
			Type:  proto.MetricType_GAUGE,
			Value: value,
		},
	})
	return err
}

func (g *GrpcMetricClient) UpdateCounter(ctx context.Context, name string, value int64) error {
	_, err := g.Update(ctx, &proto.UpdateMetric{
		Metric: &proto.Metric{
			Id:    name,
			Type:  proto.MetricType_COUNTER,
			Delta: value,
		},
	})
	return err
}

func (g *GrpcMetricClient) BatchUpdate(ctx context.Context, metric []*contracts.Metric) error {
	in := make([]*proto.Metric, len(metric))
	for i, m := range metric {
		in[i] = &proto.Metric{
			Id: m.ID,
		}
		switch m.Type {
		case models.GaugeType:
			in[i].Type = proto.MetricType_GAUGE
			in[i].Value = m.Value
		case models.CounterType:
			in[i].Type = proto.MetricType_COUNTER
			in[i].Delta = m.Delta
		}
	}
	_, err := g.MetricsClient.BatchUpdate(ctx, &proto.BatchUpdateMetric{
		Metric: in,
	})
	return err
}

func NewGRPCMetricClient(client *grpc.ClientConn) *GrpcMetricClient {
	return &GrpcMetricClient{
		MetricsClient: proto.NewMetricsClient(client),
	}
}
