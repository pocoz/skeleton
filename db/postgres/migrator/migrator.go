package migrator

import (
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/pocoz/skeleton/models"
)

type Engine struct {
	ctx     context.Context
	logger  log.Logger
	options *options
}

type Config struct {
	Context       context.Context
	Logger        log.Logger
	CredentialsDB *models.CredentialsDB
	Folder        string
	Table         string
}

type options struct {
	connString string
	folder     string
}

func New(cfg *Config) *Engine {
	folder := fmt.Sprintf("file://%s", cfg.Folder)
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?x-migrations-table=%s&sslmode=disable",
		cfg.CredentialsDB.User,
		cfg.CredentialsDB.Password,
		cfg.CredentialsDB.Server,
		cfg.CredentialsDB.Port,
		cfg.CredentialsDB.Database,
		cfg.Table,
	)

	engine := &Engine{
		ctx:    cfg.Context,
		logger: cfg.Logger,
		options: &options{
			connString: connString,
			folder:     folder,
		},
	}

	return engine
}

func (e *Engine) Run() error {
	migration, err := migrate.New(e.options.folder, e.options.connString)
	if err != nil {
		return err
	}
	defer func() {
		srcErr, dbErr := migration.Close()
		if srcErr != nil {
			level.Warn(e.logger).Log("msg", "source migrations close was failure", "err", srcErr)
		}
		if dbErr != nil {
			level.Warn(e.logger).Log("msg", "database migrations close was failure", "err", dbErr)
		}
	}()

	version, _, _ := migration.Version()
	level.Info(e.logger).Log("msg", "start migrations", "version", version)

	err = migration.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			level.Info(e.logger).Log("msg", "migrations no changes detected")

			return nil
		}

		return err
	}

	version, _, _ = migration.Version()
	level.Info(e.logger).Log("msg", "migrations successfully completed", "version", version)

	return nil
}
