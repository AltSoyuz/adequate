package store

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/AltSoyuz/adequate/internal/store/dal"
	"github.com/AltSoyuz/adequate/lib/db"
	"github.com/AltSoyuz/adequate/lib/logger"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations/*.sql
var migFS embed.FS

type Store struct {
	DB      *sql.DB
	Queries *dal.Queries
}

func Init(ctx context.Context, path string) *Store {
	// Ensure the directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		logger.Fatal("store.init.mkdir", "err", err)
	}

	tr := &db.Tracer{
		SlowThresh:  50 * time.Millisecond,
		MaskArgs:    true,
		SampleEvery: 1,
	}
	db.RegisterTracedDriver(tr)

	dsn := fmt.Sprintf("file:%s?_foreign_keys=on&_busy_timeout=5000&_journal_mode=WAL", path)

	sqlDb, err := sql.Open(db.TracedDriverName, dsn)
	if err != nil {
		logger.Fatal("store.init", "err", err)
	}
	sqlDb.SetMaxOpenConns(1)
	sqlDb.SetMaxIdleConns(1)

	if err := sqlDb.PingContext(ctx); err != nil {
		_ = sqlDb.Close()
		logger.Fatal("store.init.ping", "err", err)
	}

	err = db.Migrate(ctx, sqlDb, migFS)
	if err != nil {
		_ = sqlDb.Close()
		logger.Fatal("store.migrate", "err", err)
	}

	q := dal.New(sqlDb)
	return &Store{DB: sqlDb, Queries: q}
}

func (s *Store) Close() {
	if err := s.DB.Close(); err != nil {
		logger.Error("store.close", "err", err)
	}
}

func WithTx(ctx context.Context, db *sql.DB, fn func(ctx context.Context, q *dal.Queries) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			logger.Error("store.tx.rollback", "error", err)
		}
	}()

	q := dal.New(tx)
	if err := fn(ctx, q); err != nil {
		return err
	}
	return tx.Commit()
}
