package persistence

import (
	"errors"
	"github.com/DimKa163/go-metrics/internal/files"
	"github.com/DimKa163/go-metrics/internal/models"
	"io"
	"sync"
)

var ErrValueNotFound = errors.New("value not found")

type StoreOption struct {
	Restore bool
	UseSYNC bool
}

type MemStorage struct {
	metrics map[string]*models.Metric
	mutex   *sync.RWMutex
	option  StoreOption
}

func ConfigureStore(filer files.Filer, options StoreOption) (Repository, error) {
	store, err := NewStore(filer, options)
	if err != nil {
		return nil, err
	}
	if options.UseSYNC {
		return NewSyncStore(filer, store)
	}
	return store, nil
}

func NewStore(filer files.Filer, options StoreOption) (Repository, error) {
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
		metrics: data,
		option:  options,
		mutex:   &sync.RWMutex{},
	}, nil
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
		return &*val, nil
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
	if _, ok := s.metrics[metric.ID]; ok {
		delete(s.metrics, metric.ID)
	}
	s.metrics[metric.ID] = metric
	return nil
}

type SyncMemStorage struct {
	Repository
	filer files.Filer
}

func NewSyncStore(filer files.Filer, inner Repository) (Repository, error) {
	return &SyncMemStorage{
		Repository: inner,
		filer:      filer,
	}, nil
}
func (s *SyncMemStorage) Find(key string) *models.Metric {
	return s.Repository.Find(key)
}

func (s *SyncMemStorage) Get(key string) (*models.Metric, error) {
	return s.Repository.Get(key)
}

func (s *SyncMemStorage) GetAll() []models.Metric {
	return s.Repository.GetAll()

}
func (s *SyncMemStorage) Upsert(metric *models.Metric) error {
	err := s.Repository.Upsert(metric)
	if err != nil {
		return err
	}
	return s.filer.Dump(s.Repository.GetAll())
}
