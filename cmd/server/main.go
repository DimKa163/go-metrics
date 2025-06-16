package main

import (
	"github.com/DimKa163/go-metrics/internal/handlers"
	"github.com/DimKa163/go-metrics/internal/logging"
	"github.com/DimKa163/go-metrics/internal/middleware"
	"github.com/DimKa163/go-metrics/internal/persistence"
	"github.com/gin-gonic/gin"
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
	router := setup()
	router.LoadHTMLFiles("views/home.tmpl")
	store := persistence.NewMemStorage()
	router.GET("/", handlers.HomeHandler(store))
	router.GET("/value/:type/:name", handlers.GetHandler(store))
	router.POST("/update/:type/:name/:value", handlers.Update(store))
	return router.Run(addr)
}

func setup() *gin.Engine {
	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(middleware.WithLogging())
	return router
}
