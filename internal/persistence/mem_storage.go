package persistence

import (
	"errors"
	"github.com/DimKa163/go-metrics/internal/models"
)

var ErrValueNotFound = errors.New("value not found")
var ErrValueAlreadyExist = errors.New("value already exist")

type MemStorage struct {
	metrics map[models.MetricType]map[string]*models.Metric
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		metrics: make(map[models.MetricType]map[string]*models.Metric),
	}
}

func (s *MemStorage) Find(mType models.MetricType, key string) *models.Metric {
	if val, ok := s.metrics[mType][key]; ok {
		return val
	}
	return nil
}

func (s *MemStorage) Get(mType models.MetricType, key string) (*models.Metric, error) {
	if val, ok := s.metrics[mType]; ok {
		if val, ok := val[key]; ok {
			return val, nil
		}
	}
	return nil, ErrValueNotFound
}

func (s *MemStorage) GetAll() []models.Metric {
	var result []models.Metric
	for _, metricGroup := range s.metrics {
		for _, metric := range metricGroup {
			result = append(result, *metric)
		}
	}
	return result

}

func (s *MemStorage) Create(metric *models.Metric) error {
	if metricGroup, ok := s.metrics[metric.Type]; ok {
		if _, ok := metricGroup[metric.ID]; ok {
			return ErrValueAlreadyExist
		}
		metricGroup[metric.ID] = metric
		return nil
	}
	s.metrics[metric.Type] = make(map[string]*models.Metric)
	s.metrics[metric.Type][metric.ID] = metric
	return nil
}
