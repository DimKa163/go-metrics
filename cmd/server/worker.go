package main

import (
	"context"
	"github.com/DimKa163/go-metrics/internal/files"
	"github.com/DimKa163/go-metrics/internal/logging"
	"github.com/DimKa163/go-metrics/internal/persistence"
	"go.uber.org/zap"
	"time"
)

func asyncStore(ctx context.Context, repository persistence.Repository, filer files.Filer, interval time.Duration) error {
	storeTicker := time.NewTicker(interval)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-storeTicker.C:
			logging.Log.Info("Storing metrics...")
			startTime := time.Now()
			metrics := repository.GetAll()
			if err := filer.Dump(metrics); err != nil {
				logging.Log.Error("Dump with error", zap.Error(err))
			}
			elapsed := time.Since(startTime)
			logging.Log.Info("Storing metrics... done", zap.Duration("elapsed", elapsed))
		}
	}
}

func backup(repository persistence.Repository, filer files.Filer) error {
	logging.Log.Info("Back up metrics...")
	metrics := repository.GetAll()
	if err := filer.Dump(metrics); err != nil {
		logging.Log.Error("Dump with error", zap.Error(err))
		return err
	}
	return nil
}
