package main

import (
	"log"

	"github.com/DimKa163/go-metrics/app/collector"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	var config collector.Config
	ParseFlags(&config)
	app, err := collector.NewCollector(&config)
	if err != nil {
		log.Fatal(err)
	}
	if err = app.Run(buildVersion, buildDate, buildCommit); err != nil {
		log.Fatal(err)
	}
}
