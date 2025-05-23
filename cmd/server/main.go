package main

import (
	"github.com/DimKa163/go-metrics/internal/handlers"
	"github.com/DimKa163/go-metrics/internal/persistence"
	"net/http"
)

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}

func run() error {
	mux := http.NewServeMux()
	gRep := persistence.NewGaugeRepository()
	cRep := persistence.NewCounterRepository()
	updateHandler := handlers.NewUpdateHandler(gRep, cRep)
	mux.HandleFunc("/update/{type}/{name}/{value}", updateHandler.Update)
	return http.ListenAndServe(`:8080`, mux)
}
