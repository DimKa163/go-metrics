package files

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"

	"github.com/DimKa163/go-metrics/internal/models"
)

func TestDumpAndRestore(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "test_dump.json")
	attempts := []int{1, 3, 5}
	f := NewFiler(filePath, attempts)

	// данные для дампа
	delta := int64(54)
	metrics := []models.Metric{
		models.Metric{
			ID:    "TestCounter",
			Type:  models.CounterType,
			Delta: &delta,
		},
	}

	// проверяем Dump
	if err := f.Dump(metrics); err != nil {
		t.Fatalf("Dump failed: %v", err)
	}

	// проверяем Restore
	restored, err := f.Restore()
	assert.NoError(t, err)
	assert.Equal(t, len(metrics), len(restored))
	assert.Equal(t, metrics[0].ID, restored[0].ID)
	assert.Equal(t, metrics[0].Type, restored[0].Type)
	assert.Equal(t, metrics[0].Delta, restored[0].Delta)
}

func TestRestore_FileNotExist(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "not_exist.json")

	attempts := []int{1, 3, 5}
	f := NewFiler(filePath, attempts)

	restored, err := f.Restore()

	assert.Error(t, err)
	assert.Nil(t, restored)
}

func TestDump_Overwrite(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "dump.json")

	attempts := []int{1, 3, 5}
	f := NewFiler(filePath, attempts)

	delta := int64(54)
	metrics := []models.Metric{
		models.Metric{
			ID:    "TestCounter",
			Type:  models.CounterType,
			Delta: &delta,
		},
	}

	// проверяем Dump
	if err := f.Dump(metrics); err != nil {
		t.Fatalf("Dump failed: %v", err)
	}

	value := float64(54.3)
	metrics = []models.Metric{
		models.Metric{
			ID:    "TestCounter",
			Type:  models.GaugeType,
			Value: &value,
		},
	}

	// проверяем Dump
	if err := f.Dump(metrics); err != nil {
		t.Fatalf("Dump failed: %v", err)
	}

	restored, err := f.Restore()
	assert.NoError(t, err)
	assert.Equal(t, len(metrics), len(restored))
	assert.Equal(t, metrics[0].ID, restored[0].ID)
	assert.Equal(t, metrics[0].Type, restored[0].Type)
	assert.Equal(t, metrics[0].Delta, restored[0].Delta)
}
