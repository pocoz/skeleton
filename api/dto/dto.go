package dto

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/ratelimit"
	"golang.org/x/time/rate"
)

type basicRequest struct {
	Data interface{} `json:"data"`
	Meta interface{} `json:"meta"`
}

type BasicResponse struct {
	Success int         `json:"success"`
	Data    interface{} `json:"data"`
	Meta    interface{} `json:"meta"`
}

type BasicReject struct {
	Success int         `json:"success"`
	Data    DataReject  `json:"data"`
	Meta    interface{} `json:"meta"`
}

type DataReject struct {
	Err interface{} `json:"err"`
}

func NewBasicResponse(isError bool, data interface{}, meta interface{}) interface{} {
	if !isError {
		response := new(BasicResponse)
		response.Success = 1
		response.Data = data
		if meta == nil {
			meta = &struct{}{}
		}
		response.Meta = meta

		return response
	} else {
		response := new(BasicReject)
		response.Success = 0
		response.Data.Err = data
		if meta == nil {
			meta = &struct{}{}
		}
		response.Meta = meta

		return response
	}
}

func DecodeBasicRequest(r *http.Request) ([]byte, error) {
	request := basicRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(request.Data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// EncodeError writes a services error to the given http.ResponseWriter
func EncodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(NewBasicResponse(true, err.Error(), nil))
}

func ApplyMiddleware(e endpoint.Endpoint, limiter *rate.Limiter) endpoint.Endpoint {
	return ratelimit.NewErroringLimiter(limiter)(e)
}
