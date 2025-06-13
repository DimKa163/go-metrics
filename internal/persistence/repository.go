package persistence

import "github.com/DimKa163/go-metrics/internal/models"

type Repository interface {
	Find(metricType string, key string) *models.Metric
	Get(metricType string, key string) (*models.Metric, error)
	GetAll() []models.Metric
	Create(metric *models.Metric) error
}
