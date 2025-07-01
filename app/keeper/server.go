package keeper

import (
	"context"
	"github.com/DimKa163/go-metrics/internal/files"
	"github.com/DimKa163/go-metrics/internal/logging"
	"github.com/DimKa163/go-metrics/internal/mhttp/controllers"
	"github.com/DimKa163/go-metrics/internal/mhttp/middleware"
	"github.com/DimKa163/go-metrics/internal/persistence"
	"github.com/DimKa163/go-metrics/internal/tasks"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

type ServiceContainer struct {
	conf             *Config
	filer            *files.Filer
	pg               *pgxpool.Pool
	repository       persistence.Repository
	metricController controllers.Metrics
	dumpTask         *tasks.DumpTask
}

type Server struct {
	*gin.Engine
	*http.Server
	*ServiceContainer
	useDumpASYNC bool
}

func New(config *Config) (*Server, error) {
	conn, err := pgxpool.New(context.Background(), config.Database)
	if err != nil {
		return nil, err
	}
	filer := files.NewFiler(config.Path)
	store, err := persistence.NewStore(conn, filer, persistence.StoreOption{
		UseSYNC: config.StoreInterval == 0,
		Restore: config.Restore,
	})
	if err != nil {
		return nil, err
	}
	if err := logging.Initialize(config.LogLevel); err != nil {
		return nil, err
	}
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.GzipMiddleware())
	return &Server{
		ServiceContainer: &ServiceContainer{
			conf:             config,
			filer:            filer,
			repository:       store,
			metricController: controllers.NewMetricController(store),
			dumpTask:         tasks.NewDumpTask(store, filer, time.Duration(config.StoreInterval)*time.Second),
		},
		Server: &http.Server{
			Addr:    config.Addr,
			Handler: router.Handler(),
		},
		useDumpASYNC: config.StoreInterval > 0,
		Engine:       router,
	}, nil
}

func (s *Server) Map() {
	s.GET("/ping", func(c *gin.Context) {
		if err := s.repository.Ping(c.Request.Context()); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		c.String(http.StatusOK, "pong")
	})
	s.GET("/", s.metricController.Home)
	s.GET("/value/:type/:name", s.metricController.Get)
	s.POST("/value", s.metricController.GetJSON)
	s.POST("/update/:type/:name/:value", s.metricController.Update)
	s.POST("/update", s.metricController.UpdateJSON)
}

func (s *Server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	if s.useDumpASYNC {
		s.dumpTask.Start(ctx)
	}
	go func() {
		<-ctx.Done()
		if err := s.backup(); err != nil {
			logging.Log.Error("backup failed", zap.Error(err))
		}
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		_ = s.Server.Shutdown(timeoutCtx)
	}()
	return s.ListenAndServe()
}

func (s *Server) backup() error {
	logging.Log.Debug("start backup before shutdown")
	return s.filer.Dump(s.repository.GetAll())
}
