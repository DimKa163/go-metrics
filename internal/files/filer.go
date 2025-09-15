package files

import (
	"bufio"
	"context"
	"encoding/json"
	"os"

	"github.com/cenkalti/backoff/v5"

	"github.com/DimKa163/go-metrics/internal/models"
)

type Filer struct {
	path     string
	attempts []int
}

func NewFiler(path string, attempts []int) *Filer {
	return &Filer{
		path:     path,
		attempts: attempts,
	}
}

func (f *Filer) Restore() ([]models.Metric, error) {
	file, err := f.openFile(os.O_CREATE | os.O_RDWR)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	buf := bufio.NewReader(file)
	var metrics []models.Metric
	if err := json.NewDecoder(buf).Decode(&metrics); err != nil {
		return nil, err
	}
	return metrics, nil
}

func (f *Filer) Dump(metrics []models.Metric) error {
	file, err := f.openFile(os.O_CREATE | os.O_WRONLY | os.O_TRUNC)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	if err := file.Truncate(0); err != nil {
		return err
	}
	if _, err := file.Seek(0, 0); err != nil {
		return err
	}
	if err := json.NewEncoder(file).Encode(&metrics); err != nil {
		return err
	}
	return nil
}

func (f *Filer) openFile(flag int) (*os.File, error) {
	seconds := f.attempts
	attempt := 0
	return backoff.Retry(context.Background(), func() (*os.File, error) {
		file, err := os.OpenFile(f.path, flag, 0644)
		if err != nil {
			if attempt > len(seconds)-1 {
				return nil, backoff.Permanent(err)
			}
			at := attempt
			attempt++
			return nil, backoff.RetryAfter(seconds[at])
		}
		return file, nil
	})
}
