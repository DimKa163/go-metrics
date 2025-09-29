package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/DimKa163/go-metrics/internal/logging"
	"github.com/DimKa163/go-metrics/internal/models"
	"github.com/DimKa163/go-metrics/internal/persistence"
)

var ErrMetricNotFound = errors.New("metric not found")

type MetricService struct {
	repository persistence.Repository
}

func NewMetricService(repository persistence.Repository) *MetricService {
	return &MetricService{repository: repository}
}

// Get get metric
func (ms *MetricService) Get(ctx context.Context, id string) (models.Metric, error) {
	model, err := ms.repository.Find(ctx, id)
	if err != nil {
		if errors.Is(err, persistence.ErrMetricNotFound) {
			return models.Metric{}, ErrMetricNotFound
		}
		return models.Metric{}, fmt.Errorf("db unhandled error %w", err)
	}
	return *model, nil
}

// GetAll get all metric
func (ms *MetricService) GetAll(ctx context.Context) ([]models.Metric, error) {
	return ms.repository.GetAll(ctx)
}

// Upsert create/update metric
func (ms *MetricService) Upsert(ctx context.Context, newMetric models.Metric) (models.Metric, error) {
	m, err := ms.processMetric(ctx, newMetric)
	if err != nil {
		return models.Metric{}, err
	}
	err = ms.repository.Upsert(ctx, &m)
	if err != nil {
		return models.Metric{}, fmt.Errorf("db unhandled error %w", err)
	}
	return m, nil
}

// BatchUpdate create/update metrics
func (ms *MetricService) BatchUpdate(ctx context.Context, metricList []models.Metric) error {
	var err error
	mapMetric := make(map[string]models.Metric)
	for _, metric := range metricList {
		it, ok := mapMetric[metric.ID]
		if ok {
			switch it.Type {
			case models.GaugeType:
				mapMetric[metric.ID] = metric
			case models.CounterType:
				*it.Delta = *metric.Delta + *it.Delta
				mapMetric[metric.ID] = it
			}
			continue
		}
		mapMetric[metric.ID] = metric
	}
	resultList := make([]models.Metric, 0)
	var m models.Metric
	for _, metric := range mapMetric {
		m, err = ms.processMetric(ctx, metric)
		if err != nil {
			return err
		}
		resultList = append(resultList, m)
	}
	if err = ms.repository.BatchUpsert(ctx, resultList); err != nil {
		return fmt.Errorf("db unhandled error %w", err)
	}
	return nil
}

func (ms *MetricService) processMetric(ctx context.Context, metric models.Metric) (models.Metric, error) {
	m, err := ms.repository.Find(ctx, metric.ID)
	if err != nil && !errors.Is(err, persistence.ErrMetricNotFound) {
		return models.Metric{}, fmt.Errorf("db unhandled error %w", err)
	}
	if m == nil {
		logging.Log.Info("metric not found. adding new metric")
		m = &metric
	} else {
		logging.Log.Info("metric found. updating metric")
		m.Update(metric)
	}
	return *m, nil
}
