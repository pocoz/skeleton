package elasticsearch

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"sync"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/olivere/elastic"
)

// Engine ...
type Engine struct {
	client          *elastic.Client
	ctx             context.Context
	mu              sync.Mutex
	logger          log.Logger
	reconnectAmount int
	cfg             *Config
}

// Config ...
type Config struct {
	Server   string
	User     string
	Password string
	Context  context.Context
	Logger   log.Logger
}

// New ...
func New(cfg *Config) (*Engine, error) {
	engine := &Engine{
		ctx:    cfg.Context,
		mu:     sync.Mutex{},
		logger: cfg.Logger,
		cfg:    cfg,
	}

	err := engine.conn()
	if err != nil {
		return nil, err
	}

	//err = engine.setMappings()
	//if err != nil {
	//	return nil, err
	//}

	level.Info(cfg.Logger).Log("msg", "established elasticsearch connection")
	return engine, nil
}

func (e *Engine) conn() error {
	timeout := time.Second * 360
	httpClient := &http.Client{Timeout: timeout}
	client, err := elastic.NewClient(
		elastic.SetURL(e.cfg.Server),
		elastic.SetBasicAuth(e.cfg.User, e.cfg.Password),
		elastic.SetSniff(false),
		elastic.SetHttpClient(httpClient),
		elastic.SetRetrier(newElasticRetrier(time.Second*3, 10)),
	)
	if err != nil {
		return err
	}

	e.client = client

	return nil
}

func (e *Engine) setMappings() error {
	mappingsBytes, err := ioutil.ReadFile(path.Join("db/elasticsearch/mappings/visenze_index_v0-mapping.json"))
	if err != nil {
		return err
	}

	settingsBytes, err := ioutil.ReadFile(path.Join("db/elasticsearch/mappings/visenze_index_v0-settings.json"))
	if err != nil {
		return err
	}

	exist, err := e.client.IndexExists(yourIndexName).Do(e.ctx)
	if err != nil {
		return err
	}

	if !exist {
		_, err = e.client.CreateIndex(yourIndexName).BodyString(string(settingsBytes)).Do(e.ctx)
		if err != nil {
			return err
		}
	}

	_, err = e.client.PutMapping().
		Type(yourIndexDoctype).
		Index(yourIndexName).
		IncludeTypeName(true).
		BodyString(string(mappingsBytes)).
		Do(e.ctx)
	if err != nil {
		return err
	}

	return nil
}

func (e *Engine) Reconnect() error {
	if e.reconnectAmount > 5 {
		return fmt.Errorf("reconnect amount more than 5, now exit with error")
	}

	err := e.conn()
	if err != nil {
		time.Sleep(time.Minute * 1)
		e.reconnectAmount++
		return e.Reconnect()
	}

	e.reconnectAmount = 0

	return nil
}

// Shutdown close elasticsearch session
func (e *Engine) Shutdown() {
	// close elastic session
	e.mu.Lock()
	if e.client != nil {
		e.client.Stop()
	}
	e.mu.Unlock()

	level.Info(e.logger).Log("msg", "elasticsearch: shutdown complete")
}
