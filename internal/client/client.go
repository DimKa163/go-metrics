package client

import (
	"fmt"
	"net/http"
)

const (
	GaugeType   = "gauge"
	CounterType = "counter"
)

type MetricClient interface {
	UpdateGauge(name string, value float64) error
	UpdateCounter(name string, value int64) error
}

type metricClient struct {
	addr string
}

func NemClient(addr string) MetricClient {
	return &metricClient{
		addr: addr,
	}
}

func (client *metricClient) UpdateGauge(name string, value float64) error {
	fullAddr := fmt.Sprintf("%supdate/%s/%s/%f", client.addr, GaugeType, name, value)
	req, err := http.NewRequest(http.MethodPost, fullAddr, nil)
	if err != nil {
		return err
	}
	httpClient := &http.Client{}
	res, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	return nil
}

func (client *metricClient) UpdateCounter(name string, value int64) error {
	fullAddr := fmt.Sprintf("%supdate/%s/%s/%d", client.addr, CounterType, name, value)
	req, err := http.NewRequest(http.MethodPost, fullAddr, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/plain")
	httpClient := &http.Client{}
	res, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	return nil
}
