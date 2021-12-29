package root

import (
	"context"
	"net/http"
)

// Service.Healthy encoder/decoder
func decodeHealthyRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return nil, nil
}

func encodeHealthyResponse(_ context.Context, w http.ResponseWriter, _ interface{}) error {
	w.WriteHeader(http.StatusOK)

	return nil
}
