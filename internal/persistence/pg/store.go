package pg

import (
	"context"
	"database/sql"
	"errors"
	"sync"

	"github.com/cenkalti/backoff/v5"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/DimKa163/go-metrics/internal/models"
	"github.com/DimKa163/go-metrics/internal/persistence"
)

type Store struct {
	*pgxpool.Pool
	mutex    *sync.RWMutex
	attempts []int
}

func NewStore(pgs *pgxpool.Pool, attempts []int) (*Store, error) {
	if err := migrateDB(pgs); err != nil {
		return nil, err
	}
	return &Store{
		Pool:     pgs,
		mutex:    &sync.RWMutex{},
		attempts: attempts,
	}, nil
}

func (s *Store) Find(ctx context.Context, key string) (*models.Metric, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	seconds := s.attempts
	attempt := 0
	query := "SELECT id, type,  delta, value value FROM metrics WHERE id = $1;"
	metric, err := backoff.Retry(ctx, func() (*models.Metric, error) {
		var m models.Metric
		if err := s.QueryRow(ctx, query, key).Scan(&m.ID,
			&m.Type,
			&m.Delta,
			&m.Value); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, backoff.Permanent(err)
			}
			var pgerr *pgconn.PgError
			if errors.As(err, &pgerr) {
				if shouldRetry(pgerr) && attempt < len(seconds) {
					at := attempt
					attempt++
					return nil, backoff.RetryAfter(seconds[at])
				}
			}
			return nil, backoff.Permanent(err)
		}
		return &m, nil
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, persistence.ErrMetricNotFound
		}
		return nil, err
	}
	return metric, nil
}

func (s *Store) GetAll(ctx context.Context) ([]models.Metric, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	seconds := s.attempts
	attempt := 0
	query := "SELECT id, type,  delta, value FROM metrics ORDER BY id ASC;"
	metrics, err := backoff.Retry(ctx, func() ([]models.Metric, error) {
		cursor, err := s.Query(ctx, query)
		if err != nil {
			var pgerr *pgconn.PgError
			if errors.As(err, &pgerr) {
				if shouldRetry(pgerr) && attempt < len(seconds) {
					at := attempt
					attempt++
					return nil, backoff.RetryAfter(seconds[at])
				}
			}
			return nil, backoff.Permanent(err)
		}
		defer cursor.Close()
		var metrics []models.Metric
		for cursor.Next() {
			if err = ctx.Err(); err != nil {
				return nil, backoff.Permanent(err)
			}
			var metric models.Metric
			if err = cursor.Scan(&metric.ID,
				&metric.Type, &metric.Delta, &metric.Value); err != nil {
				var pgerr *pgconn.PgError
				if errors.As(err, &pgerr) {
					if shouldRetry(pgerr) && attempt < len(seconds) {
						at := attempt
						attempt++
						return nil, backoff.RetryAfter(seconds[at])
					}
				}
				return nil, backoff.Permanent(err)
			}
			metrics = append(metrics, metric)
		}
		if err = cursor.Err(); err != nil {
			return nil, backoff.Permanent(err)
		}
		return metrics, nil
	})

	return metrics, err
}

func (s *Store) Upsert(ctx context.Context, metric *models.Metric) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	deleteSQL := "DELETE FROM metrics WHERE id = $1;"
	insertSQL := "INSERT INTO metrics (id, type, delta, value) VALUES ($1, $2, $3, $4);"
	return s.execWithRetry(ctx, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, deleteSQL, metric.ID); err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, insertSQL, metric.ID, metric.Type, metric.Delta, metric.Value); err != nil {
			return err
		}
		return nil
	})
}

func (s *Store) BatchUpsert(ctx context.Context, metrics []models.Metric) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	var err error
	deleteSQL := "DELETE FROM metrics WHERE id = $1;"
	insertSQL := "INSERT INTO metrics (id, type, delta, value) VALUES ($1, $2, $3, $4);"
	return s.execWithRetry(ctx, func(tx pgx.Tx) error {
		for _, metric := range metrics {
			if _, err = tx.Exec(ctx, deleteSQL, metric.ID); err != nil {
				return err
			}
			if _, err = tx.Exec(ctx, insertSQL, metric.ID, metric.Type, metric.Delta, metric.Value); err != nil {
				return err
			}
		}
		return nil
	})
}

func migrateDB(pgx *pgxpool.Pool) error {
	var err error
	db, err := sql.Open("postgres", pgx.Config().ConnString())
	if err != nil {
		return err
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		return err
	}
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}

func (s *Store) execWithRetry(ctx context.Context, txFunc func(pgx.Tx) error) error {
	seconds := s.attempts
	attempt := 0
	_, err := backoff.Retry(ctx, func() (bool, error) {
		tx, err := s.Begin(ctx)
		if err != nil {
			return false, backoff.Permanent(err)
		}
		err = txFunc(tx)
		if err != nil {
			_ = tx.Rollback(ctx)
			var pgerr *pgconn.PgError
			if errors.As(err, &pgerr) {
				if shouldRetry(pgerr) && attempt < len(seconds) {
					at := attempt
					attempt++
					return false, backoff.RetryAfter(seconds[at])
				}
			}
			return false, backoff.Permanent(err)
		}

		if err = tx.Commit(ctx); err != nil {
			var pgerr *pgconn.PgError
			if errors.As(err, &pgerr) {
				if shouldRetry(pgerr) && attempt < len(seconds) {
					at := attempt
					attempt++
					return false, backoff.RetryAfter(seconds[at])
				}
			}
			return false, backoff.Permanent(err)
		}
		return true, nil
	})
	return err
}

func shouldRetry(pgerr *pgconn.PgError) bool {
	switch pgerr.Code {
	case pgerrcode.SerializationFailure:
	case pgerrcode.DeadlockDetected:
	case pgerrcode.TooManyConnections:
	case pgerrcode.LockNotAvailable:
	case pgerrcode.CannotConnectNow:
	case pgerrcode.QueryCanceled:
	case pgerrcode.UniqueViolation:
		return true
	}
	return false
}
