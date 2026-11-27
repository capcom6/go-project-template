package db

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/capcom6/go-project-template/internal/db/migrations"
	"github.com/go-core-fx/goosefx"
	"github.com/go-core-fx/logger"
	"github.com/go-core-fx/sqlfx"
	"github.com/pressly/goose/v3/database"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/schema"
	"go.uber.org/fx"

	_ "github.com/go-sql-driver/mysql" // required
)

func Module() fx.Option {
	return fx.Module(
		"db",
		logger.WithNamedLogger("db"),
		fx.Invoke(ensureSQLiteDataDir),
		fx.Provide(func() database.Dialect {
			return database.DialectMySQL
		}),
		fx.Provide(func() schema.Dialect {
			return mysqldialect.New()
		}),
		fx.Provide(func() goosefx.Storage {
			return goosefx.Storage(migrations.FS)
		}),
	)
}

func ensureSQLiteDataDir(cfg sqlfx.Config) error {
	u, err := url.Parse(cfg.URL)
	if err != nil {
		return fmt.Errorf("parse database url: %w", err)
	}

	if u.Scheme != "sqlite" && u.Scheme != "sqlite3" {
		return nil
	}

	filename := u.Path
	if u.Host != "" && u.Host != "localhost" && u.Host != ":memory:" {
		filename = u.Host + u.Path
	}

	dir := filepath.Dir(filename)
	if dir == "." {
		return nil
	}

	if mkdirErr := os.MkdirAll(dir, 0o750); mkdirErr != nil {
		return fmt.Errorf("create database directory: %w", mkdirErr)
	}

	return nil
}
