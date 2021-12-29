package mssql

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// Engine ...
type Engine struct {
	db       *sql.DB
	ctx      context.Context
	mu       sync.Mutex
	logger   log.Logger
	queryMap map[uint8]options
}

// Config ...
type Config struct {
	Server   string
	Port     string
	User     string
	Password string
	Database string
	Context  context.Context
	Logger   log.Logger
}

func New(cfg *Config) (*Engine, error) {
	// Build connection string
	dataSourceName := fmt.Sprintf(
		"server=%s;user id=%s;password=%s;port=%s;database=%s;",
		cfg.Server,
		cfg.User,
		cfg.Password,
		cfg.Port,
		cfg.Database,
	)

	// Create connection pool
	db, err := sql.Open("sqlserver", dataSourceName)
	if err != nil {
		return nil, err
	}
	err = db.PingContext(cfg.Context)
	if err != nil {
		return nil, err
	}

	storage := &Engine{
		db:       db,
		ctx:      cfg.Context,
		mu:       sync.Mutex{},
		logger:   cfg.Logger,
		queryMap: make(map[uint8]options),
	}
	storage.queryInit()

	level.Info(cfg.Logger).Log("msg", "established mssql connection")
	return storage, nil
}

// Shutdown close mssql session
func (e *Engine) Shutdown() {
	// close mssql session
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.db != nil {
		err := e.db.Close()
		if err != nil {
			level.Error(e.logger).Log("msg", "mssql: shutdown with error", "err", err)
			return
		}
	}

	level.Info(e.logger).Log("msg", "mssql: shutdown complete")
}
