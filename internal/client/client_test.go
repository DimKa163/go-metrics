package client

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/DimKa163/go-metrics/internal/mocks"
	"github.com/DimKa163/go-metrics/internal/models"
)

func TestUpdateGauge_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDoer := mocks.NewMockHttpExecuter(ctrl)
	c := &metricClient{client: mockDoer, addr: "http://localhost"}

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

	err := c.UpdateGauge("Alloc", 123.45)
	assert.NoError(t, err)
}

func TestUpdateCounter_FailStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDoer := mocks.NewMockHttpExecuter(ctrl)
	c := &metricClient{client: mockDoer, addr: "http://localhost"}

	mockDoer.EXPECT().
		Do(gomock.Any()).
		Return(&http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       io.NopCloser(bytes.NewBufferString("bad")),
		}, nil)

	err := c.UpdateCounter("Requests", 10)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected status code")
}

func TestBatchUpdate_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDoer := mocks.NewMockHttpExecuter(ctrl)
	c := &metricClient{client: mockDoer, addr: "http://localhost"}

	metrics := []*models.Metric{
		models.CreateGauge("Alloc", 1.23),
		models.CreateCounter("Requests", 10),
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

	err := c.BatchUpdate(metrics)
	assert.NoError(t, err)
}
