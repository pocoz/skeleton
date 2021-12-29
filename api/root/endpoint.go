package root

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

func makeHealthyEndpoint(svc service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		svc.healthy()

		return healthyResponse{}, nil
	}
}

type healthyResponse struct{}
