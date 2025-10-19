// Package keeper application for storing runtime metric
package keeper

import (
	"context"
	"fmt"
	"github.com/DimKa163/go-metrics/internal/gc/gserver"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"net"
	"os/signal"
	"syscall"
	"time"

	"github.com/DimKa163/go-metrics/internal/crypto"
	"github.com/DimKa163/go-metrics/internal/files"
	"github.com/DimKa163/go-metrics/internal/logging"
	"github.com/DimKa163/go-metrics/internal/mhttp/controllers"
	"github.com/DimKa163/go-metrics/internal/mhttp/middleware"
	"github.com/DimKa163/go-metrics/internal/persistence"
	"github.com/DimKa163/go-metrics/internal/persistence/mem"
	"github.com/DimKa163/go-metrics/internal/persistence/pg"
	"github.com/DimKa163/go-metrics/internal/tasks"
	"github.com/DimKa163/go-metrics/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ServerImpl interface {
	ListenAndServe() error
	Map()
	Shutdown(ctx context.Context) error
}

type ServiceContainer struct {
	Conf             *Config
	Filer            *files.Filer
	Pg               *pgxpool.Pool
	Repository       persistence.Repository
	MetricController controllers.Metrics
	DumpTask         *tasks.DumpTask
	Crypto           *crypto.Decrypter
	GrpcService      *gserver.MetricServer
}

type Server struct {
	ServerImpl
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
	var decrypter *crypto.Decrypter
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

	if config.PrivateKeyFilePath != "" {
		decrypter, err = crypto.NewDecrypter(config.PrivateKeyFilePath)
		if err != nil {
			return nil, err
		}
	}
	if err = logging.Initialize(config.LogLevel); err != nil {
		return nil, err
	}
	server := &Server{
		ServiceContainer: &ServiceContainer{
			Conf:             config,
			Pg:               pgConnection,
			Filer:            filer,
			Repository:       repository,
			MetricController: controllers.NewMetricController(usecase.NewMetricService(repository)),
			DumpTask:         tasks.NewDumpTask(repository, filer, time.Duration(config.StoreInterval)*time.Second),
			Crypto:           decrypter,
			GrpcService:      gserver.NewMetricServer(usecase.NewMetricService(repository)),
		},
		useDumpASYNC: useDumpASYNC,
		useBackup:    useBackup,
	}
	server.ServerImpl, err = CreateServerImpl(server.ServiceContainer, config)
	if err != nil {
		return nil, err
	}
	return server, nil
}

// Map routes
func (s *Server) Map() {
	s.ServerImpl.Map()
}

// Run app
func (s *Server) Run(buildVersion string, buildDate string, buildCommit string) error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()
	if s.useDumpASYNC {
		s.DumpTask.Start(ctx)
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
		_ = s.ServerImpl.Shutdown(timeoutCtx)
	}()
	printBuildInfo(buildVersion, buildDate, buildCommit)
	return s.ListenAndServe()
}

func printBuildInfo(buildVersion string, buildDate string, buildCommit string) {
	fmt.Printf("Build version: %s\n", ifNan(buildVersion))
	fmt.Printf("Build date: %s\n", ifNan(buildDate))
	fmt.Printf("Build commit: %s\n", ifNan(buildCommit))
}

func ifNan(value string) string {
	if value == "" {
		return "N/A"
	}
	return value
}

func (s *Server) backup(ctx context.Context) error {
	logging.Log.Info("start backup before shutdown")
	m, err := s.Repository.GetAll(ctx)
	if err != nil {
		return err
	}
	return s.Filer.Dump(m)
}

func CreateServerImpl(services *ServiceContainer, config *Config) (ServerImpl, error) {
	var httpServer *HTTPServer
	var grpcServer *GRPCServer
	var err error
	if config.Addr != "" {
		router := gin.New()
		router.Use(gin.Recovery())
		router.Use(middleware.LoggingMiddleware())
		router.Use(middleware.GzipMiddleware())
		if config.TrustedSubnet != "" {
			_, ipNet, err := net.ParseCIDR(config.TrustedSubnet)
			if err != nil {
				return nil, err
			}
			router.Use(middleware.IdentifyMiddleware(ipNet))
		}
		if services.Crypto != nil {
			router.Use(middleware.CryptoMiddleware(services.Crypto))
		}
		if config.Key != "" {
			router.Use(middleware.Hash(config.Key))
		}
		httpServer, err = NewHTTPServer(services, config.Addr, router), nil
		if err != nil {
			return nil, err
		}
	}
	if config.GRPCAddr != "" {
		listener, err := net.Listen("tcp", config.GRPCAddr)
		if err != nil {
			return nil, err
		}
		chain := make([]grpc.UnaryServerInterceptor, 0)
		chain = append(chain, gserver.UnaryLoggingInterceptor())
		if config.TrustedSubnet != "" {
			_, ipNet, err := net.ParseCIDR(config.TrustedSubnet)
			if err != nil {
				return nil, err
			}
			chain = append(chain, gserver.UnaryIdentifyInterceptor(ipNet))
		}

		serv := grpc.NewServer(grpc.ChainUnaryInterceptor(chain...))
		grpcServer, err = NewGRPCServer(listener, serv, services), nil
		if err != nil {
			return nil, err
		}
	}
	if httpServer != nil && grpcServer != nil {
		return NewCompositeServer(grpcServer, httpServer), nil
	} else if httpServer != nil {
		return httpServer, nil
	} else {
		return grpcServer, nil
	}
}

type CompositeServer struct {
	grpcServer *GRPCServer
	httpServer *HTTPServer
}

func NewCompositeServer(grpcServer *GRPCServer, httpServer *HTTPServer) *CompositeServer {
	return &CompositeServer{
		grpcServer: grpcServer,
		httpServer: httpServer,
	}
}

func (cs *CompositeServer) Map() {
	cs.grpcServer.Map()
	cs.httpServer.Map()
}

func (cs *CompositeServer) Shutdown(ctx context.Context) error {
	errGroup, ctx := errgroup.WithContext(ctx)
	errGroup.Go(func() error {
		return cs.grpcServer.Shutdown(ctx)
	})
	errGroup.Go(func() error {
		return cs.httpServer.Shutdown(ctx)
	})
	return errGroup.Wait()
}

func (cs *CompositeServer) ListenAndServe() error {
	errGroup, _ := errgroup.WithContext(context.Background())
	errGroup.Go(cs.grpcServer.ListenAndServe)
	errGroup.Go(cs.httpServer.ListenAndServe)
	return errGroup.Wait()
}
