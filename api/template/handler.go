package template

import (
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"golang.org/x/time/rate"

	"github.com/pocoz/skeleton/api/dto"
)

// NewHandler creates a new http handler serving services endpoints.
func NewHandler(srv service, logger log.Logger, limiter *rate.Limiter) http.Handler {
	var (
		svc    = &loggingMiddleware{next: srv, logger: logger}
		router = mux.NewRouter()
	)

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(dto.EncodeError),
	}

	templateEndpoint := makeTemplateEndpoint(svc)
	templateEndpoint = dto.ApplyMiddleware(templateEndpoint, limiter)

	router.Path("/api/internal/v1/templatemicroservice/do/something").Methods(http.MethodPost).Handler(kithttp.NewServer(
		templateEndpoint,
		decodeTemplateRequest,
		encodeTemplateResponse,
		opts...,
	))

	return router
}
