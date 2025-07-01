package persistence

import (
	"context"
	"errors"
	"github.com/DimKa163/go-metrics/internal/files"
	"github.com/DimKa163/go-metrics/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
	"sync"
)

var ErrValueNotFound = errors.New("value not found")

type StoreOption struct {
	Restore bool
	UseSYNC bool
}

type MemStorage struct {
	pg      *pgxpool.Pool
	metrics map[string]*models.Metric
	filer   *files.Filer
	mutex   *sync.RWMutex
	option  StoreOption
}

func NewStore(pg *pgxpool.Pool, filer *files.Filer, options StoreOption) (Repository, error) {
	data := make(map[string]*models.Metric)
	if options.Restore {
		metrics, err := filer.Restore()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return nil, err
			}
		}
		for _, metric := range metrics {
			data[metric.ID] = &metric
		}
	}
	return &MemStorage{
		pg:      pg,
		metrics: data,
		option:  options,
		filer:   filer,
		mutex:   &sync.RWMutex{},
	}, nil
}

func (s *MemStorage) Ping(ctx context.Context) error {
	return s.pg.Ping(ctx)
}

func (s *MemStorage) Find(key string) *models.Metric {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if val, ok := s.metrics[key]; ok {
		return val
	}
	return nil
}

func (s *MemStorage) Get(key string) (*models.Metric, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if val, ok := s.metrics[key]; ok {
		return val, nil
	}
	return nil, ErrValueNotFound
}

func (s *MemStorage) GetAll() []models.Metric {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	var result []models.Metric
	for _, metric := range s.metrics {
		result = append(result, *metric)
	}
	return result

}

func (s *MemStorage) Upsert(metric *models.Metric) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.metrics, metric.ID)
	s.metrics[metric.ID] = metric
	if s.option.UseSYNC {
		var result []models.Metric
		for _, met := range s.metrics {
			result = append(result, *met)
		}
		return s.filer.Dump(result)
	}
	return nil
}
