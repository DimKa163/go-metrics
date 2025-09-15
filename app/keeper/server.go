package keeper

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/DimKa163/go-metrics/internal/files"
	"github.com/DimKa163/go-metrics/internal/logging"
	"github.com/DimKa163/go-metrics/internal/mhttp/controllers"
	"github.com/DimKa163/go-metrics/internal/mhttp/middleware"
	"github.com/DimKa163/go-metrics/internal/persistence"
	"github.com/DimKa163/go-metrics/internal/persistence/mem"
	"github.com/DimKa163/go-metrics/internal/persistence/pg"
	"github.com/DimKa163/go-metrics/internal/tasks"
	"github.com/DimKa163/go-metrics/internal/usecase"
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
	useBackup    bool
}

func New(config *Config) (*Server, error) {
	var repository persistence.Repository
	var err error
	var pgConnection *pgxpool.Pool
	var useDumpASYNC bool
	var useBackup bool
	attempts := []int{1, 3, 5}
	filer := files.NewFiler(config.Path, attempts)

	if config.DatabaseDSN != "" {
		pgConnection, err = pgxpool.New(context.Background(), config.DatabaseDSN)
		if err != nil {
			return nil, err
		}
		repository, err = pg.NewStore(pgConnection, attempts)
		if err != nil {
			return nil, err
		}
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
	if config.Key != "" {
		router.Use(middleware.Hash(config.Key))
	}
	return &Server{
		ServiceContainer: &ServiceContainer{
			conf:             config,
			pg:               pgConnection,
			filer:            filer,
			repository:       repository,
			metricController: controllers.NewMetricController(usecase.NewMetricService(repository)),
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
	pprof.Register(s.Engine)
	s.GET("/ping", func(c *gin.Context) {
		if s.pg != nil {
			if err := s.pg.Ping(c); err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}
		c.String(http.StatusOK, "pong")
	})
	s.metricController.Map(s.Engine)
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
	logging.Log.Info("start backup before shutdown")
	m, err := s.repository.GetAll(ctx)
	if err != nil {
		return err
	}
	return s.filer.Dump(m)
}
