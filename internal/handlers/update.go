package handlers

import (
	"github.com/DimKa163/go-metrics/internal/persistence"
	"net/http"
	"strconv"
	"strings"
)

const (
	GaugeType   = "gauge"
	CounterType = "counter"
)

type UpdateHandler interface {
	Update(w http.ResponseWriter, r *http.Request)
}

type updateHandler struct {
	gaugeRepository   persistence.GaugeRepository
	counterRepository persistence.CounterRepository
}

func NewUpdateHandler(gaugeRepository persistence.GaugeRepository, counterRepository persistence.CounterRepository) UpdateHandler {
	return &updateHandler{
		gaugeRepository:   gaugeRepository,
		counterRepository: counterRepository,
	}
}

func (h *updateHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	segments := strings.Split(r.URL.Path, "/")[2:]
	t := segments[0]
	switch t {
	case GaugeType:
		name := segments[1]
		value, err := strconv.ParseFloat(segments[2], 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		h.gaugeRepository.Update(name, value)
	case CounterType:
		name := segments[1]
		value, err := strconv.ParseInt(segments[2], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		h.counterRepository.Increment(name, value)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}
