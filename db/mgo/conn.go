package mgo

import (
	"context"
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Engine ...
type Engine struct {
	client *mongo.Client
	db     *mongo.Database
	ctx    context.Context
	mu     sync.Mutex
	logger log.Logger
	cfg    *Config
}

// Config ...
type Config struct {
	Server   string
	NameDB   string
	User     string
	Password string
	Context  context.Context
	Logger   log.Logger
}

func New(cfg *Config) (*Engine, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(cfg.Server))
	if err != nil {
		return nil, err
	}

	err = client.Connect(cfg.Context)
	if err != nil {
		return nil, err
	}

	engine := &Engine{
		client: client,
		db:     client.Database(cfg.NameDB),
		ctx:    cfg.Context,
		mu:     sync.Mutex{},
		logger: cfg.Logger,
		cfg:    cfg,
	}

	level.Info(cfg.Logger).Log("msg", "established mongo connection")
	return engine, err
}

// Shutdown close mongo session
func (e *Engine) Shutdown() {
	// close mongo session
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.client != nil {
		err := e.client.Disconnect(e.ctx)
		if err != nil {
			level.Info(e.logger).Log("msg", "mongo: shutdown with error", "err", err)
			return
		}
	}

	level.Info(e.logger).Log("msg", "mongo: shutdown complete")
}
