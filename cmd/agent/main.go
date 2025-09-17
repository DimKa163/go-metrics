package main

import (
	"log"

	"github.com/DimKa163/go-metrics/app/collector"
)

func main() {
	var config collector.Config
	ParseFlags(&config)
	app := collector.NewCollector(&config)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
