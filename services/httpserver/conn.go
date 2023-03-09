package httpserver

import (
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"golang.org/x/time/rate"

	"github.com/pocoz/skeleton/api/root"
	"github.com/pocoz/skeleton/api/template"
	"github.com/pocoz/skeleton/db/elasticsearch"
	"github.com/pocoz/skeleton/db/mssql"
)

// ServerHTTP is a services http server
type ServerHTTP struct {
	logger log.Logger
	srv    *http.Server
}

// Config is a http server configuration
type Config struct {
	Port        string
	Logger      log.Logger
	RateLimiter *rate.Limiter
	DBES        *elasticsearch.Engine
	DBMsSQL     *mssql.Engine
}

// New creates a new http server
func New(cfg *Config) (*ServerHTTP, error) {
	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:         "0.0.0.0:" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	server := &ServerHTTP{
		logger: cfg.Logger,
		srv:    srv,
	}

	svcRoot := &root.BasicService{}
	handlerRoot := root.NewHandler(svcRoot, cfg.Logger, cfg.RateLimiter)
	mux.Handle("/", accessControl(handlerRoot))

	svcTemplate := &template.BasicService{
		Logger:  cfg.Logger,
		DBES:    cfg.DBES,
		DBMsSQL: cfg.DBMsSQL,
	}
	handler := template.NewHandler(svcTemplate, cfg.Logger, cfg.RateLimiter)
	mux.Handle("/api/internal/v1/templatemicroservice/", accessControl(handler))

	return server, nil
}

// CORS headers
func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

// Run starts the server.
func (s *ServerHTTP) Run() error {
	err := s.srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Shutdown stopped the http server.
func (s *ServerHTTP) Shutdown() {
	err := s.srv.Close()
	if err != nil {
		level.Info(s.logger).Log("msg", "http server: shutdown has err", "err:", err)
		return
	}
	level.Info(s.logger).Log("msg", "http server: shutdown complete")
}
