package hclient

import (
	"bytes"
	"context"
	"github.com/DimKa163/go-metrics/internal/mhttp/contracts"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/DimKa163/go-metrics/internal/mocks"
)

func TestUpdateGauge_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	mockDoer := mocks.NewMockHttpExecuter(ctrl)
	c := &HTTPMetricClient{client: mockDoer, addr: "http://localhost"}

	mockDoer.EXPECT().
		Do(gomock.Any()).
		DoAndReturn(func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, "http://localhost/update", req.URL.String())
			body, _ := io.ReadAll(req.Body)
			assert.Contains(t, string(body), `"type":"gauge"`)
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString("ok")),
			}, nil
		})

	err := c.UpdateGauge(ctx, "Alloc", 123.45)
	assert.NoError(t, err)
}

func TestUpdateCounter_FailStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	mockDoer := mocks.NewMockHttpExecuter(ctrl)
	c := &HTTPMetricClient{client: mockDoer, addr: "http://localhost"}

	mockDoer.EXPECT().
		Do(gomock.Any()).
		Return(&http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       io.NopCloser(bytes.NewBufferString("bad")),
		}, nil)

	err := c.UpdateCounter(ctx, "Requests", 10)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected status code")
}

func TestBatchUpdate_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockDoer := mocks.NewMockHttpExecuter(ctrl)
	c := &HTTPMetricClient{client: mockDoer, addr: "http://localhost"}

	metrics := []*contracts.Metric{
		contracts.CreateGauge("Alloc", 1.23),
		contracts.CreateCounter("Requests", 10),
	}

	mockDoer.EXPECT().
		Do(gomock.Any()).
		DoAndReturn(func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, "http://localhost/updates", req.URL.String())
			body, _ := io.ReadAll(req.Body)
			assert.Contains(t, string(body), `"Alloc"`)
			assert.Contains(t, string(body), `"Requests"`)
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString("ok")),
			}, nil
		})

	err := c.BatchUpdate(ctx, metrics)
	assert.NoError(t, err)
}

func ExampleNewClient() {
	// поднимаем тестовый сервер, чтобы не ходить наружу
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	ctx := context.Background()
	// создаём клиента
	c := NewClient(server.URL, nil)

	// обновляем gauge
	_ = c.UpdateGauge(ctx, "cpu", 0.95)

	// обновляем counter
	_ = c.UpdateCounter(ctx, "requests", 10)

	// batch update
	_ = c.BatchUpdate(ctx, []*contracts.Metric{
		contracts.CreateGauge("memory", 128.0),
		contracts.CreateCounter("hits", 42),
	})
}
