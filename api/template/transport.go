package template

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/dieqnt/skeleton/api/dto"
)

// Service.Template encoder/decoder
func decodeTemplateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	data, err := dto.DecodeBasicRequest(r)
	if err != nil {
		return nil, err
	}

	req := templateRequest{}
	err = json.Unmarshal(data, &req)
	return req, err
}

func encodeTemplateResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	res := response.(templateResponse)
	if res.Err != nil {
		return json.NewEncoder(w).Encode(dto.NewBasicResponse(true, res, nil))
	}

	return json.NewEncoder(w).Encode(dto.NewBasicResponse(false, res, nil))
}
