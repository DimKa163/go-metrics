package controllers

import (
	"bytes"
	"compress/gzip"
	"context"
	"github.com/DimKa163/go-metrics/internal/files"
	"github.com/DimKa163/go-metrics/internal/mhttp/middleware"
	"github.com/DimKa163/go-metrics/internal/models"
	"github.com/DimKa163/go-metrics/internal/persistence/mem"
	"github.com/DimKa163/go-metrics/internal/usecase"
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
		expectedStatusCode int
	}{
		{
			name:               "success gauge",
			method:             http.MethodPost,
			url:                "/update/gauge/TestMetric/100",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "success counter",
			method:             http.MethodPost,
			url:                "/update/counter/TestCounter/2",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "wrong type",
			method:             http.MethodPost,
			url:                "/update/otherType/TestMetric/100",
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			router := gin.Default()
			sut := NewMetricController(configureService())
			sut.Map(router)
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
		expectedStatusCode int
	}{
		{
			name:               "success gauge",
			method:             http.MethodPost,
			url:                "/update/",
			expectedStatusCode: http.StatusOK,
			body:               "{\n    \"id\": \"TestMetric\",\n    \"type\": \"gauge\",\n    \"value\": 100.0\n}",
		},
		{
			name:               "success counter",
			method:             http.MethodPost,
			url:                "/update/",
			body:               "{\n    \"id\": \"TestMetric\",\n    \"type\": \"counter\",\n    \"delta\": 100\n}",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "wrong type",
			method:             http.MethodPost,
			url:                "/update/",
			body:               "{\n    \"id\": \"TestMetric\",\n    \"type\": \"otherType\",\n    \"delta\": 100\n}",
			expectedStatusCode: http.StatusBadRequest,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			router := gin.Default()
			sut := NewMetricController(configureService())
			sut.Map(router)
			res := httptest.NewRecorder()
			router.ServeHTTP(res, httptest.NewRequest(c.method, c.url, strings.NewReader(c.body)))
			assert.Equal(t, c.expectedStatusCode, res.Code)
		})
	}
}

func TestUpdatesJSON(t *testing.T) {
	cases := []struct {
		name               string
		method             string
		body               string
		url                string
		expectedStatusCode int
	}{
		{
			name:               "success batch update",
			method:             http.MethodPost,
			url:                "/updates",
			body:               `[{"id": "TestCounterMetric", "type": "counter", "delta": 100}, {"id": "TestGaugeMetric", "type": "gauge", "value": 200.23}]`,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "bad batch update",
			method:             http.MethodPost,
			url:                "/updates",
			body:               `[{"id": "TestMetric", "type": "counte", "delta": 100}, {"id": "TestMetric", "type": "gauge", "value": 200.23}]`,
			expectedStatusCode: http.StatusBadRequest,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			router := gin.Default()
			sut := NewMetricController(configureService())
			sut.Map(router)
			res := httptest.NewRecorder()
			router.ServeHTTP(res, httptest.NewRequest(c.method, c.url, strings.NewReader(c.body)))
			assert.Equal(t, c.expectedStatusCode, res.Code)
		})
	}
}

func TestGetJSON(t *testing.T) {
	cases := []struct {
		name               string
		method             string
		url                string
		body               string
		expectedStatusCode int
	}{
		{
			name:               "success get gauge",
			method:             http.MethodPost,
			url:                "/value/",
			body:               `{"type": "gauge", "id": "FoundedGaugeMetric"}`,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "success get counter",
			method:             http.MethodPost,
			url:                "/value/",
			body:               `{"type": "counter", "id": "FoundedCounterMetric"}`,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "not found get gauge",
			method:             http.MethodPost,
			url:                "/value/",
			body:               `{"type": "gauge", "id": "NotFoundedGaugeMetric"}`,
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:               "not found get counter",
			method:             http.MethodPost,
			url:                "/value/",
			body:               `{"type": "counter", "id": "NotFoundedCounterMetric"}`,
			expectedStatusCode: http.StatusNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			router := gin.Default()
			sut := NewMetricController(configureService())
			sut.Map(router)
			res := httptest.NewRecorder()
			router.ServeHTTP(res, httptest.NewRequest(c.method, c.url, strings.NewReader(c.body)))
			assert.Equal(t, c.expectedStatusCode, res.Code)
		})
	}
}

func TestGet(t *testing.T) {
	cases := []struct {
		name               string
		method             string
		url                string
		body               string
		expectedStatusCode int
	}{
		{
			name:               "success get gauge",
			method:             http.MethodGet,
			url:                "/value/gauge/FoundedGaugeMetric",
			body:               ``,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "success get counter",
			method:             http.MethodGet,
			url:                "/value/counter/FoundedCounterMetric",
			body:               ``,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "not found get gauge",
			method:             http.MethodGet,
			url:                "/value/gauge/NotFoundedGaugeMetric",
			body:               ``,
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:               "not found get counter",
			method:             http.MethodGet,
			url:                "/value/counter/NotFoundedCounterMetric",
			body:               ``,
			expectedStatusCode: http.StatusNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			router := gin.Default()
			sut := NewMetricController(configureService())
			sut.Map(router)
			res := httptest.NewRecorder()
			router.ServeHTTP(res, httptest.NewRequest(c.method, c.url, nil))
			assert.Equal(t, c.expectedStatusCode, res.Code)
		})
	}
}

func TestUpdateGzip(t *testing.T) {
	router := gin.Default()
	router.Use(middleware.GzipMiddleware())
	metrics := NewMetricController(configureService())
	router.POST("/update", metrics.UpdateJSON)
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

func configureService() *usecase.MetricService {
	return usecase.NewMetricService(configureFileRepository())
}
func configureFileRepository() *mem.MemoryStore {
	attempts := []int{1, 3, 5}
	filer := files.NewFiler("test_dump", attempts)
	repository, _ := mem.NewStore(filer, mem.StoreOption{
		UseSYNC: false,
		Restore: false,
	})
	delta := int64(29)
	value := float64(29.9)
	repository.Upsert(context.Background(), &models.Metric{
		ID:    "FoundedCounterMetric",
		Type:  models.CounterType,
		Delta: &delta,
		Value: nil,
	})
	repository.Upsert(context.Background(), &models.Metric{
		ID:    "FoundedGaugeMetric",
		Type:  models.CounterType,
		Delta: nil,
		Value: &value,
	})
	return repository
}

type mockGaugeRepository struct {
	data map[string]map[string]*models.Metric
}

func (m *mockGaugeRepository) Ping(_ context.Context) error {
	return nil
}
func (m *mockGaugeRepository) Find(_ context.Context, _ string) (*models.Metric, error) {
	return nil, nil
}
func (m *mockGaugeRepository) GetAll(_ context.Context) ([]models.Metric, error) {
	return make([]models.Metric, 0), nil
}
func (m *mockGaugeRepository) Upsert(_ context.Context, _ *models.Metric) error {
	return nil
}

func (m *mockGaugeRepository) BatchUpsert(_ context.Context, _ []models.Metric) error {
	return nil
}
