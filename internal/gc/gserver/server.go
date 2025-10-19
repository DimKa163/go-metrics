package gserver

import (
	"context"
	"errors"
	"github.com/DimKa163/go-metrics/internal/gc/proto"
	"github.com/DimKa163/go-metrics/internal/models"
	"github.com/DimKa163/go-metrics/internal/persistence"
	"github.com/DimKa163/go-metrics/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MetricServer struct {
	app *usecase.MetricService
	proto.UnimplementedMetricsServer
}

func NewMetricServer(appService *usecase.MetricService) *MetricServer {
	return &MetricServer{app: appService}
}

func (ms *MetricServer) Register(server *grpc.Server) {
	proto.RegisterMetricsServer(server, ms)
}
func (ms *MetricServer) Get(ctx context.Context, in *proto.GetMetric) (*proto.MetricResponse, error) {
	var resp proto.MetricResponse
	m, err := ms.app.Get(ctx, in.Name)
	if err != nil {
		if errors.Is(err, persistence.ErrMetricNotFound) {
			return nil, status.Error(codes.NotFound, "resource not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	resp.Metric = toOut(m)
	return &resp, nil
}

func (ms *MetricServer) Update(ctx context.Context, in *proto.UpdateMetric) (*proto.MetricResponse, error) {
	var resp proto.MetricResponse
	m, err := ms.app.Upsert(ctx, toIn(in.Metric))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	resp.Metric = toOut(m)
	return &resp, nil
}

func (ms *MetricServer) BatchUpdate(ctx context.Context, in *proto.BatchUpdateMetric) (*proto.BatchUpdateMetricResponse, error) {
	var resp proto.BatchUpdateMetricResponse
	metrics := make([]*models.Metric, len(in.Metric))
	for i, m := range in.Metric {
		metrics[i] = toIn(m)
	}
	err := ms.app.BatchUpdate(ctx, metrics)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &resp, nil
}

func toOut(metric *models.Metric) *proto.Metric {
	var m proto.Metric
	switch metric.Type {
	case models.CounterType:
		m.Type = proto.MetricType_COUNTER
		m.Delta = metric.Delta
	case models.GaugeType:
		m.Type = proto.MetricType_GAUGE
		m.Value = metric.Value
	}
	m.Id = metric.ID
	return &m
}

func toIn(metric *proto.Metric) *models.Metric {
	var m models.Metric
	switch metric.Type {
	case proto.MetricType_COUNTER:
		m.Type = models.CounterType
		m.Delta = metric.Delta
	case proto.MetricType_GAUGE:
		m.Type = models.GaugeType
		m.Value = metric.Value
	}
	m.ID = metric.Id
	return &m
}
