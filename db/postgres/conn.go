package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	_ "github.com/lib/pq"

	"github.com/pocoz/skeleton/models"
)

// Engine ...
type Engine struct {
	ctx                context.Context
	mu                 sync.Mutex
	logger             log.Logger
	db                 *sql.DB
	queryMapAdvertiser map[uint8]models.QueryOptions
	options            *options
}

// Config ...
type Config struct {
	Context       context.Context
	Logger        log.Logger
	MaxButchSize  int
	CredentialsDB *models.CredentialsDB
}

type options struct {
	driverName   string
	connString   string
	maxButchSize int
}

func New(cfg *Config) (*Engine, error) {
	// Build connection string for malibu
	connString := fmt.Sprintf(
		"host=%s user=%s password=%s port=%s dbname=%s sslmode=disable",
		cfg.CredentialsDB.Server,
		cfg.CredentialsDB.User,
		cfg.CredentialsDB.Password,
		cfg.CredentialsDB.Port,
		cfg.CredentialsDB.Database,
	)
	driverName := "postgres"

	// Create connection pool malibu
	db, err := sql.Open(driverName, connString)
	if err != nil {
		return nil, err
	}
	err = db.PingContext(cfg.Context)
	if err != nil {
		return nil, err
	}

	storage := &Engine{
		ctx:                cfg.Context,
		logger:             cfg.Logger,
		mu:                 sync.Mutex{},
		db:                 db,
		queryMapAdvertiser: make(map[uint8]models.QueryOptions),
		options: &options{
			driverName:   driverName,
			connString:   connString,
			maxButchSize: cfg.MaxButchSize,
		},
	}
	storage.initQueryAdvertister()

	level.Info(cfg.Logger).Log("msg", "established postgres connections")

	return storage, nil
}

func (e *Engine) PingReconnect(attempt uint, err error) {
	level.Error(e.logger).Log("msg", "PingReconnect: upstairs method failed", "err", err)
	level.Info(e.logger).Log("msg", "PingReconnect processing", "attempt", attempt, "total_attempts", models.RetryFibonacciAmount)

	err = e.db.Ping()
	if err == nil {
		return
	}

	// Create connection pool malibu
	db, err := sql.Open(e.options.driverName, e.options.connString)
	if err != nil {
		level.Error(e.logger).Log("msg", "PingReconnect: create new connection failure", "err", err)

		return
	}

	e.db = db
}

// Shutdown close mssql session
func (e *Engine) Shutdown() {
	// close mssql session
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.db != nil {
		err := e.db.Close()
		if err != nil {
			level.Error(e.logger).Log("msg", "mssql: malibu shutdown with error", "err", err)

			return
		}
	}

	level.Info(e.logger).Log("msg", "mssql: shutdown complete")
}
