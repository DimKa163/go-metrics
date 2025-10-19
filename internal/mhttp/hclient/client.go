package hclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/DimKa163/go-metrics/internal/mhttp/contracts"
	"github.com/DimKa163/go-metrics/internal/models"
	"net/http"
	"time"
)

const (
	HTTP  = "http"
	HTTPS = "https"
)

type HTTPExecuter interface {
	Do(req *http.Request) (*http.Response, error)
}

type HTTPMetricClient struct {
	client HTTPExecuter
	addr   string
}

type RequestHandlerFactory func(transport http.RoundTripper) http.RoundTripper

func NewClient(addr string, transports ...RequestHandlerFactory) *HTTPMetricClient {
	var transport http.RoundTripper
	defaultTransport := &http.Transport{}
	transport = defaultTransport
	for _, t := range transports {
		transport = t(transport)
	}
	return &HTTPMetricClient{
		client: &http.Client{
			Transport: transport,
			Timeout:   30 * time.Second,
		},
		addr: addr,
	}
}

// UpdateGauge create/update gauge metric
func (c *HTTPMetricClient) UpdateGauge(_ context.Context, name string, value float64) error {
	metric := models.CreateGauge(name, value)
	return c.send(metric)
}

// UpdateCounter create/update gauge metric
func (c *HTTPMetricClient) UpdateCounter(_ context.Context, name string, value int64) error {
	metric := models.CreateCounter(name, value)
	return c.send(metric)
}

// BatchUpdate create/update many metrics
func (c *HTTPMetricClient) BatchUpdate(_ context.Context, metrics []*contracts.Metric) error {
	req, err := c.createBatchRequest(metrics)
	if err != nil {
		return err
	}
	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	return nil
}

func (c *HTTPMetricClient) send(metric *models.Metric) error {
	req, err := c.createRequest(metric)
	if err != nil {
		return err
	}
	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	return nil
}

func (c *HTTPMetricClient) createRequest(metric *models.Metric) (*http.Request, error) {
	fullAddr := fmt.Sprintf("%s/update", c.addr)

	data, err := json.Marshal(metric)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(data)

	req, err := http.NewRequest(http.MethodPost, fullAddr, buffer)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Encoding", "gzip")

	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

func (c *HTTPMetricClient) createBatchRequest(metrics []*contracts.Metric) (*http.Request, error) {
	fullAddr := fmt.Sprintf("%s/updates", c.addr)
	data, err := json.Marshal(metrics)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, fullAddr, bytes.NewBuffer(data))

	if err != nil {
		return nil, err
	}

	return req, nil
}
