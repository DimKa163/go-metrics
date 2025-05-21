package main

import (
	"net/http"
	"strconv"
	"strings"
)

var memStorage MemStorage

func init() {
	memStorage = MemStorage{
		Gauge:   make(map[string]float64),
		Counter: make(map[string]int64),
	}
}
func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}

func run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/update/{type}/{name}/{value}", update)
	return http.ListenAndServe(`:8080`, mux)
}

func update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	segments := strings.Split(r.URL.Path, "/")[2:]
	t := segments[0]
	name := segments[1]
	value := segments[2]
	switch t {
	case "gauge":
		i, err := strconv.ParseFloat(value, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		updateGauge(name, i)
		break
	case "counter":
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		updateCounter(name, i)
		break
	}
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.Header().Add("Content-Length", "0")
	w.WriteHeader(http.StatusOK)
}

func updateGauge(name string, value float64) {
	memStorage.Gauge[name] = value
}

func updateCounter(name string, value int64) {
	if val, ok := memStorage.Counter[name]; ok {
		memStorage.Counter[name] = val + value
	} else {
		memStorage.Counter[name] = value
	}
}
