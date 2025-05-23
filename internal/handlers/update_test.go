package handlers

import (
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
			expectedStatusCode: http.StatusOK},
		{
			name:               "wrong method",
			method:             http.MethodGet,
			url:                "/update/gauge/TestMetric/100",
			expectedStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name:               "wrong path",
			method:             http.MethodPost,
			url:                "/update/TestMetric/100",
			expectedStatusCode: http.StatusBadRequest,
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
			req := httptest.NewRequest(c.method, c.url, nil)
			res := httptest.NewRecorder()
			updateHandler := NewUpdateHandler(&mockGaugeRepository{},
				&mockCounterRepository{})
			updateHandler.Update(res, req)
			assert.Equal(t, c.expectedStatusCode, res.Code)
		})
	}
}

type mockGaugeRepository struct {
}

func (m *mockGaugeRepository) Update(_ string, _ float64) {

}

type mockCounterRepository struct {
}

func (m *mockCounterRepository) Increment(_ string, _ int64) {

}
