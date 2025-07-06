package mem

import (
	"context"
	"errors"
	"github.com/DimKa163/go-metrics/internal/files"
	"github.com/DimKa163/go-metrics/internal/models"
	"io"
	"sync"
)

type StoreOption struct {
	Restore bool
	UseSYNC bool
}

type MemoryStore struct {
	metrics map[string]*models.Metric
	filer   *files.Filer
	mutex   *sync.RWMutex
	option  StoreOption
}

func NewStore(filer *files.Filer, options StoreOption) (*MemoryStore, error) {
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
	return &MemoryStore{
		metrics: data,
		option:  options,
		filer:   filer,
		mutex:   &sync.RWMutex{},
	}, nil
}

func (s *MemoryStore) Find(_ context.Context, key string) (*models.Metric, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if val, ok := s.metrics[key]; ok {
		return val, nil
	}
	return nil, nil
}

func (s *MemoryStore) GetAll(_ context.Context) ([]models.Metric, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	var result []models.Metric
	for _, metric := range s.metrics {
		result = append(result, *metric)
	}
	return result, nil

}

func (s *MemoryStore) Upsert(_ context.Context, metric *models.Metric) error {
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

func (s *MemoryStore) BatchUpsert(_ context.Context, metrics []models.Metric) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, metric := range metrics {
		delete(s.metrics, metric.ID)
		s.metrics[metric.ID] = &metric
	}
	if s.option.UseSYNC {
		var result []models.Metric
		for _, met := range s.metrics {
			result = append(result, *met)
		}
		return s.filer.Dump(result)
	}
	return nil
}
