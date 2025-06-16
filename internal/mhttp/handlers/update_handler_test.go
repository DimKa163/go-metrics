package handlers

import (
	"bytes"
	"compress/gzip"
	"github.com/DimKa163/go-metrics/internal/mhttp/middleware"
	"github.com/DimKa163/go-metrics/internal/models"
	"github.com/DimKa163/go-metrics/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
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
		container          *services.ServiceContainer
		expectedStatusCode int
	}{
		{
			name:   "success gauge",
			method: http.MethodPost,
			url:    "/update/gauge/TestMetric/100",
			container: &services.ServiceContainer{
				Repository: &mockGaugeRepository{
					data: make(map[string]map[string]*models.Metric),
				},
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:   "success counter",
			method: http.MethodPost,
			url:    "/update/counter/TestCounter/2",
			container: &services.ServiceContainer{
				Repository: &mockGaugeRepository{
					data: make(map[string]map[string]*models.Metric),
				},
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:   "wrong type",
			method: http.MethodPost,
			url:    "/update/otherType/TestMetric/100",
			container: &services.ServiceContainer{
				Repository: &mockGaugeRepository{
					data: make(map[string]map[string]*models.Metric),
				},
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			router := gin.Default()
			updateHandler := Update(c.container)
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
		container          *services.ServiceContainer
		expectedStatusCode int
	}{
		{
			name:   "success gauge",
			method: http.MethodPost,
			url:    "/update",
			container: &services.ServiceContainer{
				Repository: &mockGaugeRepository{
					data: make(map[string]map[string]*models.Metric),
				},
			},
			expectedStatusCode: http.StatusOK,
			body:               "{\n    \"id\": \"TestMetric\",\n    \"type\": \"gauge\",\n    \"value\": 100.0\n}",
		},
		{
			name:   "success counter",
			method: http.MethodPost,
			url:    "/update",
			container: &services.ServiceContainer{
				Repository: &mockGaugeRepository{
					data: make(map[string]map[string]*models.Metric),
				},
			},
			body:               "{\n    \"id\": \"TestMetric\",\n    \"type\": \"counter\",\n    \"delta\": 100\n}",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:   "wrong type",
			method: http.MethodPost,
			url:    "/update",
			container: &services.ServiceContainer{
				Repository: &mockGaugeRepository{
					data: make(map[string]map[string]*models.Metric),
				},
			},
			body:               "{\n    \"id\": \"TestMetric\",\n    \"type\": \"otherType\",\n    \"delta\": 100\n}",
			expectedStatusCode: http.StatusBadRequest,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			router := gin.Default()
			updateHandler := UpdateJSON(c.container)
			router.POST("/update", updateHandler)
			res := httptest.NewRecorder()
			router.ServeHTTP(res, httptest.NewRequest(c.method, c.url, strings.NewReader(c.body)))
			assert.Equal(t, c.expectedStatusCode, res.Code)
		})
	}
}

func TestUpdateGzip(t *testing.T) {
	container := &services.ServiceContainer{
		Repository: &mockGaugeRepository{
			data: make(map[string]map[string]*models.Metric),
		},
	}
	router := gin.Default()
	router.Use(middleware.GzipMiddleware())
	updateHandler := UpdateJSON(container)
	router.Use(middleware.GzipMiddleware())
	router.POST("/update", updateHandler)
	requestBody := `{
		"id": "Alloc",
		"type": "gauge",
    	"value": 1001.233
	}`
	successBody := `{
		"id": "Alloc",
		"type": "gauge",
    	"value": 1001.233
	}`
	t.Run("success gauge", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(buf)
		_, err := zb.Write([]byte(requestBody))
		require.NoError(t, err)
		err = zb.Close()
		require.NoError(t, err)
		r := httptest.NewRequest("POST", "/update", buf)
		r.Header.Set("Content-Encoding", "gzip")
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Accept-Encoding", "")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, r)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.Code)

		b, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.JSONEq(t, successBody, string(b))
	})
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
