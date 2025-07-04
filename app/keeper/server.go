package keeper

import (
	"context"
	"github.com/DimKa163/go-metrics/internal/files"
	"github.com/DimKa163/go-metrics/internal/logging"
	"github.com/DimKa163/go-metrics/internal/mhttp/controllers"
	"github.com/DimKa163/go-metrics/internal/mhttp/middleware"
	"github.com/DimKa163/go-metrics/internal/persistence"
	"github.com/DimKa163/go-metrics/internal/persistence/mem"
	"github.com/DimKa163/go-metrics/internal/persistence/pg"
	"github.com/DimKa163/go-metrics/internal/tasks"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

type ServiceContainer struct {
	conf             *Config
	filer            *files.Filer
	pg               *pgx.Conn
	repository       persistence.Repository
	metricController controllers.Metrics
	dumpTask         *tasks.DumpTask
}

type Server struct {
	*gin.Engine
	*http.Server
	*ServiceContainer
	useDumpASYNC bool
	useBackup    bool
}

func New(config *Config) (*Server, error) {
	var repository persistence.Repository
	var err error
	var pgConnection *pgx.Conn
	var useDumpASYNC bool
	var useBackup bool
	filer := files.NewFiler(config.Path)

	if config.DatabaseDSN != "" {
		pgConnection, err = pgx.Connect(context.Background(), config.DatabaseDSN)
		if err != nil {
			return nil, err
		}
		repository, err = pg.NewStore(pgConnection)
	} else {
		repository, err = mem.NewStore(filer, mem.StoreOption{
			UseSYNC: config.StoreInterval == 0,
			Restore: config.Restore,
		})

		if err != nil {
			return nil, err
		}
		useDumpASYNC = config.StoreInterval > 0
		useBackup = true
	}
	if err = logging.Initialize(config.LogLevel); err != nil {
		return nil, err
	}
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.GzipMiddleware())
	return &Server{
		ServiceContainer: &ServiceContainer{
			conf:             config,
			pg:               pgConnection,
			filer:            filer,
			repository:       repository,
			metricController: controllers.NewMetricController(repository),
			dumpTask:         tasks.NewDumpTask(repository, filer, time.Duration(config.StoreInterval)*time.Second),
		},
		Server: &http.Server{
			Addr:    config.Addr,
			Handler: router.Handler(),
		},
		useDumpASYNC: useDumpASYNC,
		Engine:       router,
		useBackup:    useBackup,
	}, nil
}

func (s *Server) Map() {
	s.GET("/ping", func(c *gin.Context) {
		if s.pg != nil {
			if err := s.pg.Ping(c); err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
			}
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
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if s.useBackup {
			if err := s.backup(timeoutCtx); err != nil {
				logging.Log.Error("backup failed", zap.Error(err))
			}
		}
		_ = s.Server.Shutdown(timeoutCtx)
	}()
	return s.ListenAndServe()
}

func (s *Server) backup(ctx context.Context) error {
	logging.Log.Debug("start backup before shutdown")
	m, err := s.repository.GetAll(ctx)
	if err != nil {
		return err
	}
	return s.filer.Dump(m)
}
