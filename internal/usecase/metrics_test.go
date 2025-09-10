package usecase

import (
	"context"
	"github.com/DimKa163/go-metrics/internal/mocks"
	"github.com/DimKa163/go-metrics/internal/models"
	"github.com/DimKa163/go-metrics/internal/persistence"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetShouldReturnMetricWhenMetricExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	mockRepository := mocks.NewMockRepository(ctrl)
	service := NewMetricService(mockRepository)
	metric := getTestCounterMetric(500)
	id := metric.ID
	mockRepository.EXPECT().Find(ctx, id).Return(&metric, nil)

	sut, err := service.Get(ctx, id)

	assert.NoError(t, err, "get metric should be successful")
	assert.Equal(t, metric, sut)
}

func TestGetShouldReturnMetricWhenMetricDoesNotExist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	mockRepository := mocks.NewMockRepository(ctrl)
	service := NewMetricService(mockRepository)
	id := "NotExistsMetric"
	metric := models.Metric{}
	mockRepository.EXPECT().Find(ctx, id).Return(nil, persistence.ErrMetricNotFound)
	sut, err := service.Get(ctx, id)

	assert.ErrorIs(t, err, ErrMetricNotFound)
	assert.Equal(t, metric, sut)
}

func TestGetAllShouldSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	mockRepository := mocks.NewMockRepository(ctrl)
	service := NewMetricService(mockRepository)
	metrics := []models.Metric{
		getTestCounterMetric(5),
		getTestGaugeMetric(23.32),
	}
	mockRepository.EXPECT().GetAll(ctx).Return(metrics, nil)

	sut, err := service.GetAll(ctx)

	assert.NoError(t, err, "get all metrics should be successful")
	assert.Equal(t, metrics, sut)
}

func TestUpdateGaugeWhenMetricExistShouldSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	mockRepository := mocks.NewMockRepository(ctrl)
	service := NewMetricService(mockRepository)
	exitsMetric := getTestGaugeMetric(23.32)
	newMetric := exitsMetric
	value := float64(300.23)
	newMetric.Value = &value
	id := exitsMetric.ID
	mockRepository.EXPECT().Find(ctx, id).Return(&exitsMetric, nil)
	mockRepository.EXPECT().Upsert(ctx, &newMetric).Return(nil)

	sut, err := service.Update(ctx, newMetric)

	assert.NoError(t, err, "update metric should be successful")
	assert.Equal(t, newMetric, sut, "update metric should return true metric")
}

func TestUpdateGaugeWhenMetricDoesNotExistShouldSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	mockRepository := mocks.NewMockRepository(ctrl)
	service := NewMetricService(mockRepository)
	newMetric := getTestGaugeMetric(23.23)
	id := newMetric.ID
	mockRepository.EXPECT().Find(ctx, id).Return(nil, persistence.ErrMetricNotFound)
	mockRepository.EXPECT().Upsert(ctx, &newMetric).Return(nil)

	sut, err := service.Update(ctx, newMetric)

	assert.NoError(t, err, "update metric should be successful")
	assert.Equal(t, newMetric, sut, "update metric should return true metric")
}

func TestUpdateCounterWhenMetricExistShouldSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	mockRepository := mocks.NewMockRepository(ctrl)
	service := NewMetricService(mockRepository)
	exitsMetric := getTestCounterMetric(5)
	expectedMetric := getTestCounterMetric(155)
	newMetric := exitsMetric
	delta := int64(150)
	newMetric.Delta = &delta
	id := exitsMetric.ID

	mockRepository.EXPECT().Find(ctx, id).Return(&exitsMetric, nil)
	mockRepository.EXPECT().Upsert(ctx, &exitsMetric).Return(nil)

	sut, err := service.Update(ctx, newMetric)

	exitsMetric.Update(newMetric)
	assert.NoError(t, err, "update metric should be successful")
	assert.Equal(t, expectedMetric, sut, "update metric should return true metric")
}

func TestUpdateCounterWhenMetricDoesNotExistShouldSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	mockRepository := mocks.NewMockRepository(ctrl)
	service := NewMetricService(mockRepository)
	newMetric := getTestCounterMetric(500)
	id := newMetric.ID
	mockRepository.EXPECT().Find(ctx, id).Return(nil, persistence.ErrMetricNotFound)
	mockRepository.EXPECT().Upsert(ctx, &newMetric).Return(nil)

	sut, err := service.Update(ctx, newMetric)

	assert.NoError(t, err, "update metric should be successful")
	assert.Equal(t, newMetric, sut, "update metric should return true metric")
}

func TestBatchUpdateShouldSuccess(t *testing.T) {

}

func getTestGaugeMetric(value float64) models.Metric {
	metric := models.Metric{
		ID:    "TestGaugeMetric",
		Type:  models.GaugeType,
		Value: &value,
	}
	return metric
}

func getTestCounterMetric(delta int64) models.Metric {
	metric := models.Metric{
		ID:    "TestCounterMetric",
		Type:  models.CounterType,
		Delta: &delta,
	}
	return metric
}
