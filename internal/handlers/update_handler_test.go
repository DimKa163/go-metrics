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
		data               map[string]map[string]*models.Metric
		expectedStatusCode int
	}{
		{
			name:               "success gauge",
			method:             http.MethodPost,
			url:                "/update/gauge/TestMetric/100",
			data:               make(map[string]map[string]*models.Metric),
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "success counter",
			method:             http.MethodPost,
			url:                "/update/counter/TestCounter/2",
			data:               make(map[string]map[string]*models.Metric),
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "wrong type",
			method:             http.MethodPost,
			url:                "/update/otherType/TestMetric/100",
			data:               make(map[string]map[string]*models.Metric),
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

func TestUpdateJSON(t *testing.T) {
	cases := []struct {
		name               string
		method             string
		url                string
		body               string
		expectedBody       int
		data               map[string]map[string]*models.Metric
		expectedStatusCode int
	}{
		{
			name:               "success gauge",
			method:             http.MethodPost,
			url:                "/update",
			data:               make(map[string]map[string]*models.Metric),
			expectedStatusCode: http.StatusOK,
			body:               "{\n    \"id\": \"TestMetric\",\n    \"type\": \"gauge\",\n    \"value\": 100.0\n}",
		},
		{
			name:               "success counter",
			method:             http.MethodPost,
			url:                "/update",
			data:               make(map[string]map[string]*models.Metric),
			body:               "{\n    \"id\": \"TestMetric\",\n    \"type\": \"counter\",\n    \"delta\": 100\n}",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "wrong type",
			method:             http.MethodPost,
			url:                "/update",
			data:               make(map[string]map[string]*models.Metric),
			body:               "{\n    \"id\": \"TestMetric\",\n    \"type\": \"otherType\",\n    \"delta\": 100\n}",
			expectedStatusCode: http.StatusBadRequest,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			router := gin.Default()
			updateHandler := UpdateJSON(&mockGaugeRepository{
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
	data map[string]map[string]*models.Metric
}

func (m *mockGaugeRepository) Find(metricType string, key string) *models.Metric {
	return nil
}
func (m *mockGaugeRepository) Get(metricType string, key string) (*models.Metric, error) {
	return nil, nil
}
func (m *mockGaugeRepository) GetAll() []models.Metric {
	return make([]models.Metric, 0)
}
func (m *mockGaugeRepository) Create(metric *models.Metric) error {
	return nil
}
