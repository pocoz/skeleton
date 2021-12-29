package template

import (
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/transport"
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
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(dto.EncodeError),
	}

	templateEndpoint := makeTemplateEndpoint(svc)
	templateEndpoint = applyMiddleware(templateEndpoint, limiter)

	router.Path("/api/internal/v1/templatemicroservice/do/something").Methods(http.MethodPost).Handler(kithttp.NewServer(
		templateEndpoint,
		decodeTemplateRequest,
		encodeTemplateResponse,
		opts...,
	))

	return router
}

func applyMiddleware(e endpoint.Endpoint, limiter *rate.Limiter) endpoint.Endpoint {
	return ratelimit.NewErroringLimiter(limiter)(e)
}
