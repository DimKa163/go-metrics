package keeper

import (
	"context"
	docs "github.com/DimKa163/go-metrics/docs"
	"github.com/DimKa163/go-metrics/internal/logging"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

type HTTPServer struct {
	services *ServiceContainer
	Engine   *gin.Engine
	*http.Server
}

func NewHTTPServer(services *ServiceContainer, address string, engine *gin.Engine) *HTTPServer {
	return &HTTPServer{
		services: services,
		Engine:   engine,
		Server: &http.Server{
			Addr:    address,
			Handler: engine.Handler(),
		},
	}
}

func (hs *HTTPServer) ListenAndServe() error {
	logging.Log.Info("Starting HTTP server")
	return hs.Server.ListenAndServe()
}

func (hs *HTTPServer) Map() {
	docs.SwaggerInfo.BasePath = ""
	hs.Engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	hs.Engine.GET("/ping", func(c *gin.Context) {
		if hs.services.Pg != nil {
			if err := hs.services.Pg.Ping(c); err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}
		c.String(http.StatusOK, "pong")
	})
	hs.services.MetricController.Map(hs.Engine)
}

func (hs *HTTPServer) Shutdown(ctx context.Context) error {
	return hs.Server.Shutdown(ctx)
}
