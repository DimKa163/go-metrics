package main

import (
	"github.com/DimKa163/go-metrics/app/collector"
	"log"
)

func main() {
	var config collector.Config
	ParseFlags(&config)
	app := collector.NewScrapper(&config)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
