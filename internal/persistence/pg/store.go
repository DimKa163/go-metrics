package pg

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DimKa163/go-metrics/internal/models"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"sync"
)

type Store struct {
	*pgx.Conn
	mutex *sync.RWMutex
}

func NewStore(pgx *pgx.Conn) (*Store, error) {
	if err := migrateDB(pgx); err != nil {
		return nil, err
	}
	return &Store{
		Conn:  pgx,
		mutex: &sync.RWMutex{},
	}, nil
}

func (s *Store) Find(ctx context.Context, key string) (*models.Metric, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	query := "SELECT id, type,  delta, value value FROM metrics WHERE id = $1;"
	var metric models.Metric
	if err := s.QueryRow(ctx, query, key).Scan(
		&metric.ID,
		&metric.Type,
		&metric.Delta,
		&metric.Value,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &metric, nil
}

func (s *Store) GetAll(ctx context.Context) ([]models.Metric, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	var err error
	query := "SELECT id, type,  delta, value FROM metrics ORDER BY id ASC;"

	cursor, err := s.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()
	var metrics []models.Metric
	for cursor.Next() {
		if err = ctx.Err(); err != nil {
			return nil, err
		}
		var metric models.Metric
		if err = cursor.Scan(&metric.ID,
			&metric.Type, &metric.Delta, &metric.Value); err != nil {
			return nil, err
		}
		metrics = append(metrics, metric)
	}
	if err = cursor.Err(); err != nil {
		return nil, err
	}
	return metrics, nil
}

func (s *Store) Upsert(ctx context.Context, metric *models.Metric) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	deleteSQL := "DELETE FROM metrics WHERE id = $1;"
	if _, err := s.Exec(ctx, deleteSQL, metric.ID); err != nil {
		return err
	}
	insertSQL := "INSERT INTO metrics (id, type, delta, value) VALUES ($1, $2, $3, $4);"
	if _, err := s.Exec(ctx, insertSQL, metric.ID, metric.Type, metric.Delta, metric.Value); err != nil {
		return err
	}
	return nil
}

func migrateDB(pgx *pgx.Conn) error {
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
