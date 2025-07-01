package persistence

import (
	"context"
	"github.com/DimKa163/go-metrics/internal/models"
)

type Repository interface {
	Ping(ctx context.Context) error

	Find(key string) *models.Metric

	Get(key string) (*models.Metric, error)

	GetAll() []models.Metric

	Upsert(metric *models.Metric) error
}
