package main

import (
	"context"
	"errors"
	"github.com/DimKa163/go-metrics/internal/logging"
	"go.uber.org/zap"
	"net/http"
	"os/signal"
	"syscall"
)

func main() {
	config := getServerConfig()
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	if err := run(ctx, &config); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			logging.Log.Fatal("Failed to run server", zap.Error(err))
		}
	}
}

func run(ctx context.Context, config *ServerConfig) error {
	if err := logging.Initialize(config.LogLevel); err != nil {
		return err
	}
	serviceBuilder := NewServerBuilder(config)
	if err := serviceBuilder.ConfigureServices(); err != nil {
		panic(err)
	}
	srv := serviceBuilder.Build()
	srv.LoadHTMLFiles("views/home.tmpl")
	srv.Route()
	return srv.Run(ctx)
}
