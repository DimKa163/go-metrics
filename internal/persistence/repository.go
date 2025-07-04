package persistence

import (
	"context"
	"github.com/DimKa163/go-metrics/internal/models"
)

type Repository interface {
	Find(ctx context.Context, key string) (*models.Metric, error)

	GetAll(ctx context.Context) ([]models.Metric, error)

	Upsert(ctx context.Context, metric *models.Metric) error
}
