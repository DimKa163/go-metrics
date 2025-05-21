package handlers

import (
	"errors"
	"github.com/DimKa163/go-metrics/internal/persistence"
	"strconv"
)

var BadRequestError = errors.New("value is not a number")

var TypeNotFoundError = errors.New("type not found")

func Update(t string, name string, value string) error {
	switch t {
	case "gauge":
		i, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return BadRequestError
		}
		updateGauge(name, i)
	case "counter":
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return BadRequestError
		}
		updateCounter(name, i)
	default:
		return TypeNotFoundError
	}
	return nil
}

func updateGauge(name string, value float64) {
	persistence.Store.Gauge[name] = value
}

func updateCounter(name string, value int64) {
	if val, ok := persistence.Store.Counter[name]; ok {
		persistence.Store.Counter[name] = val + value
	} else {
		persistence.Store.Counter[name] = value
	}
}
