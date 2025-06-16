package mhttp

import (
	"github.com/DimKa163/go-metrics/internal/mhttp/handlers"
	"github.com/DimKa163/go-metrics/internal/mhttp/middleware"
	"github.com/DimKa163/go-metrics/internal/services"
	"github.com/gin-gonic/gin"
)

func NewRouter(container *services.ServiceContainer) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.GzipMiddleware())
	router.GET("/", handlers.HomeHandler(container))
	router.GET("/value/:type/:name", handlers.GetHandler(container))
	router.POST("/value", handlers.GetHandlerJSON(container))
	router.POST("/update/:type/:name/:value", handlers.Update(container))
	router.POST("/update", handlers.UpdateJSON(container))
	return router
}
