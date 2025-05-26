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
	client http.Client
	addr   string
}

func NemClient(addr string) MetricClient {
	return &metricClient{
		client: http.Client{},
		addr:   addr,
	}
}

func (api *metricClient) UpdateGauge(name string, value float64) error {
	fullAddr := fmt.Sprintf("%s/update/%s/%s/%f", api.addr, GaugeType, name, value)
	req, err := http.NewRequest(http.MethodPost, fullAddr, nil)
	if err != nil {
		return err
	}
	res, err := api.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	return nil
}

func (api *metricClient) UpdateCounter(name string, value int64) error {
	fullAddr := fmt.Sprintf("%s/update/%s/%s/%d", api.addr, CounterType, name, value)
	req, err := http.NewRequest(http.MethodPost, fullAddr, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/plain")
	res, err := api.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	return nil
}
