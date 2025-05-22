package main

import (
	"errors"
	"github.com/DimKa163/go-metrics/internal/handlers"
	"net/http"
	"strings"
)

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
	err := handlers.Update(t, name, value)
	if err != nil {
		if errors.Is(err, handlers.ErrBadRequest) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if errors.Is(err, handlers.ErrTypeNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.Header().Add("Content-Length", "0")
	w.WriteHeader(http.StatusOK)
}
