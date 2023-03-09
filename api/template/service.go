package template

import (
	"fmt"

	"github.com/go-kit/kit/log"

	"github.com/pocoz/skeleton/db/elasticsearch"
	"github.com/pocoz/skeleton/db/mssql"
	"github.com/pocoz/skeleton/models"
)

type BasicService struct {
	Logger  log.Logger
	DBES    *elasticsearch.Engine
	DBMsSQL *mssql.Engine
}

type service interface {
	template(req *models.TemplateRequest) (*models.TemplateResponse, error)
}

// transfer launch transporting data with given options
func (s *BasicService) template(options *models.TemplateRequest) (*models.TemplateResponse, error) {
	fmt.Println(fmt.Sprintf("%#v", options))

	resp := &models.TemplateResponse{
		Msg:                 "nice!",
		IndicesSpent:        nil,
		IndicesInProcessing: nil,
		IndicesForProcess:   nil,
	}

	return resp, nil
}
