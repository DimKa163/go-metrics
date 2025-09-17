package main

import (
	"errors"
	"net/http"
	_ "net/http/pprof"

	"go.uber.org/zap"

	"github.com/DimKa163/go-metrics/app/keeper"
	"github.com/DimKa163/go-metrics/internal/logging"
)

func main() {
	var config keeper.Config
	ParseFlags(&config)

	app, err := keeper.New(&config)
	if err != nil {
		logging.Log.Fatal("Failed to configure keeper", zap.Error(err))
	}
	app.Map()
	app.LoadHTMLFiles("views/home.tmpl")
	if err := app.Run(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			logging.Log.Fatal("Failed to run keeper", zap.Error(err))
		}
	}
}
