package persistence

import "github.com/DimKa163/go-metrics/internal/models"

type Repository interface {
	Find(key string) *models.Metric

	Get(key string) (*models.Metric, error)

	GetAll() []models.Metric

	Upsert(metric *models.Metric) error
}
