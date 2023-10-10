package store

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"github.com/K3das/bigcord/scraping/store/db"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

//go:embed schema.sql
var ddl string

type Store struct {
	log *zap.SugaredLogger

	*db.Queries
}

func NewStore(ctx context.Context, rawLog *zap.Logger, dsn string) (*Store, error) {
	s := &Store{}
	s.log = rawLog.Sugar().With("source", "store")

	conn, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening sqlite: %w", err)
	}

	if _, err := conn.ExecContext(ctx, ddl); err != nil {
		return nil, fmt.Errorf("error running ddl: %w", err)
	}

	s.Queries = db.New(conn)

	s.log.Infof("created store: %s", dsn)

	return s, nil
}
