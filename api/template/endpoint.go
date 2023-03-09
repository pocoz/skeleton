package template

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"github.com/pocoz/skeleton/models"
)

func makeTemplateEndpoint(svc service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(templateRequest)
		res, err := svc.template(req.Options)
		return templateResponse{TemplateLaunchResult: res, Err: err}, nil
	}
}

type templateRequest struct {
	Options *models.TemplateRequest `json:"options"`
}

// Обязательно в респонсе отмечаем в json omitempty
type templateResponse struct {
	TemplateLaunchResult *models.TemplateResponse `json:"template_launch_result,omitempty"`
	Err                  error                    `json:"err,omitempty"`
}
