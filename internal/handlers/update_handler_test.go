package handlers

import (
	"github.com/DimKa163/go-metrics/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdate(t *testing.T) {
	cases := []struct {
		name               string
		method             string
		url                string
		data               map[models.MetricType]map[string]*models.Metric
		expectedStatusCode int
	}{
		{
			name:               "success gauge",
			method:             http.MethodPost,
			url:                "/update/gauge/TestMetric/100",
			data:               make(map[models.MetricType]map[string]*models.Metric),
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "success counter",
			method:             http.MethodPost,
			url:                "/update/counter/TestCounter/2",
			data:               make(map[models.MetricType]map[string]*models.Metric),
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "wrong type",
			method:             http.MethodPost,
			url:                "/update/otherType/TestMetric/100",
			data:               make(map[models.MetricType]map[string]*models.Metric),
			expectedStatusCode: http.StatusBadRequest,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			router := gin.Default()
			updateHandler := Update(&mockGaugeRepository{
				data: c.data,
			})
			router.POST("/update/:type/:name/:value", updateHandler)
			res := httptest.NewRecorder()
			router.ServeHTTP(res, httptest.NewRequest(c.method, c.url, nil))
			assert.Equal(t, c.expectedStatusCode, res.Code)
		})
	}
}

type mockGaugeRepository struct {
	data map[models.MetricType]map[string]*models.Metric
}

func (m *mockGaugeRepository) Find(metricType models.MetricType, key string) *models.Metric {
	return nil
}
func (m *mockGaugeRepository) Get(metricType models.MetricType, key string) (*models.Metric, error) {
	return nil, nil
}
func (m *mockGaugeRepository) GetAll() []models.Metric {
	return make([]models.Metric, 0)
}
func (m *mockGaugeRepository) Create(metric *models.Metric) error {
	return nil
}
