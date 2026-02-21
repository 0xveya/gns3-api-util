package db

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/stefanistkuhl/gns3util/pkg/cluster/db/sqlc"
	"github.com/stefanistkuhl/gns3util/pkg/utils"
	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var Schema string

type Store struct {
	*sqlc.Queries
	DB *sql.DB
}

func openDB(dbPath string) (*sql.DB, error) {
	dsn := fmt.Sprintf("file:%s?_foreign_keys=on&_busy_timeout=5000&_journal_mode=WAL", dbPath)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}
	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}

func Init() (*Store, error) {
	dir, err := utils.GetGNS3Dir()
	if err != nil {
		return nil, fmt.Errorf("get dir: %w", err)
	}
	dbPath := filepath.Join(dir, "clusterData.db")

	_, statErr := os.Stat(dbPath)
	dbExists := !os.IsNotExist(statErr)

	db, err := openDB(dbPath)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	if !dbExists {
		if _, err := db.Exec(Schema); err != nil {
			_ = db.Close()
			return nil, fmt.Errorf("apply initial schema: %w", err)
		}
	}

	return &Store{
		Queries: sqlc.New(db),
		DB:      db,
	}, nil
}

func (s *Store) ReadOnly(ctx context.Context, fn func(*sqlc.Queries) error) error {
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{
		ReadOnly: true,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := s.WithTx(tx)

	if err := fn(qtx); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Store) ReadOnlyTx(ctx context.Context, fn func(*sqlc.Queries) error) error {
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	qtx := s.WithTx(tx)

	if err := fn(qtx); err != nil {
		return err
	}

	return tx.Commit()
}
