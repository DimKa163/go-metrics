package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/DimKa163/go-metrics/internal/models"
	"net/http"
	"time"
)

type MetricClient interface {
	UpdateGauge(name string, value float64) error
	UpdateCounter(name string, value int64) error

	BatchUpdate(metrics []*models.Metric) error
}

type metricClient struct {
	client http.Client
	addr   string
}

func NewClient(addr string, transports []func(transport http.RoundTripper) http.RoundTripper) MetricClient {
	var transport http.RoundTripper
	defaultTransport := &http.Transport{}
	transport = defaultTransport
	for _, t := range transports {
		transport = t(transport)
	}
	return &metricClient{
		client: http.Client{
			Transport: transport,
			Timeout:   30 * time.Second,
		},
		addr: addr,
	}
}

func (c *metricClient) UpdateGauge(name string, value float64) error {
	metric := models.CreateGauge(name, value)

	req, err := c.createRequest(metric)
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

func (c *metricClient) UpdateCounter(name string, value int64) error {
	metric := models.CreateCounter(name, value)

	req, err := c.createRequest(metric)
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

func (c *metricClient) BatchUpdate(metrics []*models.Metric) error {
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

func (c *metricClient) createRequest(metric *models.Metric) (*http.Request, error) {
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

func (c *metricClient) createBatchRequest(metrics []*models.Metric) (*http.Request, error) {
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
