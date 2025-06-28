package files

import (
	"bufio"
	"encoding/json"
	"github.com/DimKa163/go-metrics/internal/models"
	"os"
)

type Filer struct {
	path string
}

func NewFiler(path string) *Filer {
	return &Filer{
		path: path,
	}
}

func (f *Filer) Restore() ([]models.Metric, error) {
	file, err := os.OpenFile(f.path, os.O_CREATE|os.O_RDWR, 0755)
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
	file, err := os.OpenFile(f.path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
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
