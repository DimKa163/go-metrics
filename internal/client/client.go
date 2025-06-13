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
}

type metricClient struct {
	client http.Client
	addr   string
}

func NewClient(addr string) MetricClient {
	return &metricClient{
		client: http.Client{
			Timeout: 30 * time.Second,
		},
		addr: addr,
	}
}

func (c *metricClient) UpdateGauge(name string, value float64) error {
	fullAddr := fmt.Sprintf("%s/update", c.addr)
	metric := models.CreateGauge(name, value)
	data, err := json.Marshal(metric)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, fullAddr, bytes.NewReader(data))
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
	fullAddr := fmt.Sprintf("%s/update", c.addr)
	metric := models.CreateCounter(name, value)
	data, err := json.Marshal(metric)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, fullAddr, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
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
