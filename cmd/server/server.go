package main

import (
	"context"
	"github.com/DimKa163/go-metrics/internal/files"
	"github.com/DimKa163/go-metrics/internal/logging"
	"github.com/DimKa163/go-metrics/internal/mhttp/handlers"
	"github.com/DimKa163/go-metrics/internal/mhttp/middleware"
	"github.com/DimKa163/go-metrics/internal/persistence"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type ServerBuilder interface {
	ConfigureServices() error
	Build() Server
}

type Services struct {
	conf       *ServerConfig
	filer      files.Filer
	repository persistence.Repository
}

func NewServerBuilder(config *ServerConfig) ServerBuilder {
	return &Services{
		conf: config,
	}
}

func (s *Services) ConfigureServices() error {
	s.filer = files.NewFiler(s.conf.Path)
	store, err := persistence.ConfigureStore(s.filer, persistence.StoreOption{
		UseSYNC: s.conf.StoreInterval == 0,
		Restore: s.conf.Restore,
	})
	if err != nil {
		return err
	}
	s.repository = store
	return nil
}

func (s *Services) Build() Server {
	router := gin.New()
	return &server{
		Services: s,
		Engine:   router,
	}
}

type Server interface {
	Route()

	LoadHTMLFiles(files ...string)
	Router() *gin.Engine

	ListenAndServe() error

	Run(ctx context.Context) error

	Shutdown(context context.Context) error
}

type server struct {
	*gin.Engine
	*http.Server
	*Services
}

func (s *server) Route() {
	s.Use(gin.Recovery())
	s.Use(middleware.LoggingMiddleware())
	s.Use(middleware.GzipMiddleware())
	s.GET("/", handlers.HomeHandler(s.repository))
	s.GET("/value/:type/:name", handlers.GetHandler(s.repository))
	s.POST("/value", handlers.GetHandlerJSON(s.repository))
	s.POST("/update/:type/:name/:value", handlers.Update(s.repository))
	s.POST("/update", handlers.UpdateJSON(s.repository))
}

func (s *server) Router() *gin.Engine {
	return s.Engine
}

func (s *server) Run(ctx context.Context) error {
	s.Server = &http.Server{
		Addr:    s.conf.Addr,
		Handler: s.Router().Handler(),
	}
	if s.conf.StoreInterval != 0 {
		go func() {
			err := asyncStore(ctx, s.repository, s.filer, time.Duration(s.conf.StoreInterval)*time.Second)
			if err != nil {
				logging.Log.Info("async store failed", zap.Error(err))
			}
		}()
	}
	go func() {
		<-ctx.Done()
		logging.Log.Info("shutting down server...")
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := backup(s.repository, s.filer); err != nil {
			logging.Log.Info("backup failed", zap.Error(err))
		}
		_ = s.Server.Shutdown(timeoutCtx)
	}()
	return s.Server.ListenAndServe()
}
