package root

import (
	"net/http"

	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"golang.org/x/time/rate"

	"github.com/dieqnt/skeleton/api/dto"
)

// NewHandler creates a new http handler serving services endpoints.
func NewHandler(srv service, logger log.Logger, limiter *rate.Limiter) http.Handler {
	var (
		svc    = &loggingMiddleware{next: srv, logger: logger}
		router = mux.NewRouter()
	)

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(dto.EncodeError),
	}

	healthCheckEndpoint := makeHealthyEndpoint(svc)
	healthCheckEndpoint = dto.ApplyMiddleware(healthCheckEndpoint, limiter)

	router.Path("/healthy").Methods(http.MethodGet).Handler(kithttp.NewServer(
		healthCheckEndpoint,
		decodeHealthyRequest,
		encodeHealthyResponse,
		opts...,
	))

	return router
}
