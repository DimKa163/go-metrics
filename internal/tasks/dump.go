package tasks

import (
	"context"
	"github.com/DimKa163/go-metrics/internal/files"
	"github.com/DimKa163/go-metrics/internal/logging"
	"github.com/DimKa163/go-metrics/internal/persistence"
	"go.uber.org/zap"
	"time"
)

type DumpTask struct {
	repository persistence.Repository
	filer      *files.Filer
	interval   time.Duration
}

func NewDumpTask(repository persistence.Repository, filer *files.Filer, interval time.Duration) *DumpTask {
	return &DumpTask{
		repository: repository,
		filer:      filer,
		interval:   interval,
	}
}

func (task *DumpTask) Start(ctx context.Context) {
	go func() {
		if err := task.run(ctx); err != nil {
			logging.Log.Error("dump task cancelled", zap.Error(err))
		}
	}()
}

func (task *DumpTask) run(ctx context.Context) error {
	storeTicker := time.NewTicker(task.interval)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-storeTicker.C:
			logging.Log.Info("Storing metrics...")
			startTime := time.Now()
			metrics, err := task.repository.GetAll(ctx)
			if err != nil {
				logging.Log.Error("dump task cancelled", zap.Error(err))
			}
			if err := task.filer.Dump(metrics); err != nil {
				logging.Log.Error("Dump with error", zap.Error(err))
			}
			elapsed := time.Since(startTime)
			logging.Log.Info("Storing metrics... done", zap.Duration("elapsed", elapsed))
		}
	}
}
