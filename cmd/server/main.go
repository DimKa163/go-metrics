package main

import (
	"github.com/DimKa163/go-metrics/internal/logging"
	"github.com/DimKa163/go-metrics/internal/mhttp"
	"github.com/DimKa163/go-metrics/internal/persistence"
	"github.com/DimKa163/go-metrics/internal/services"
)

func main() {
	parseFlag()
	err := run()
	if err != nil {
		panic(err)
	}
}

func run() error {
	if err := logging.Initialize(logLevel); err != nil {
		return err
	}
	router := mhttp.NewRouter(&services.ServiceContainer{
		Repository: persistence.NewMemStorage(),
	})
	router.LoadHTMLFiles("views/home.tmpl")
	return router.Run(addr)
}
