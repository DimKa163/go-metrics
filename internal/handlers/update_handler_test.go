package handlers

import (
	"github.com/DimKa163/go-metrics/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUpdate(t *testing.T) {
	cases := []struct {
		name               string
		method             string
		url                string
		body               string
		expectedBody       int
		data               map[models.MetricType]map[string]*models.Metric
		expectedStatusCode int
	}{
		{
			name:               "success gauge",
			method:             http.MethodPost,
			url:                "/update",
			data:               make(map[models.MetricType]map[string]*models.Metric),
			expectedStatusCode: http.StatusOK,
			body:               "{\n    \"id\": \"TestMetric\",\n    \"type\": \"gauge\",\n    \"value\": 100.0\n}",
		},
		{
			name:               "success counter",
			method:             http.MethodPost,
			url:                "/update",
			data:               make(map[models.MetricType]map[string]*models.Metric),
			body:               "{\n    \"id\": \"TestMetric\",\n    \"type\": \"counter\",\n    \"delta\": 100\n}",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "wrong type",
			method:             http.MethodPost,
			url:                "/update",
			data:               make(map[models.MetricType]map[string]*models.Metric),
			body:               "{\n    \"id\": \"TestMetric\",\n    \"type\": \"otherType\",\n    \"delta\": 100\n}",
			expectedStatusCode: http.StatusBadRequest,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			router := gin.Default()
			updateHandler := Update(&mockGaugeRepository{
				data: c.data,
			})
			router.POST("/update", updateHandler)
			res := httptest.NewRecorder()
			router.ServeHTTP(res, httptest.NewRequest(c.method, c.url, strings.NewReader(c.body)))
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
