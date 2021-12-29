package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

const ServiceName = "TEMPLATE_MICROSERVICE"

// Configuration describe services config
type Configuration struct {
	HTTPPort                string        `envconfig:"HTTP_PORT"                             required:"false" default:"2524"`
	RateLimitEvery          time.Duration `envconfig:"RATE_LIMIT_EVERY"                      required:"false" default:"1us"`
	RateLimitBurst          int           `envconfig:"RATE_LIMIT_BURST"                      required:"false" default:"100"`
	ElasticServer           string        `envconfig:"ELASTIC_SERVER"                        required:"false" default:"http://elastic-nginx.test.cloud.croc-comp.goods.local:19200"`
	ElasticUser             string        `envconfig:"ELASTIC_USER"                          required:"false" default:""`
	ElasticPassword         string        `envconfig:"ELASTIC_PASSWORD"                      required:"false" default:""`
	SQLServer               string        `envconfig:"SQL_SERVER"                            required:"false" default:"mlb-sdb-002.corp.goods.ru"`
	SQLPort                 string        `envconfig:"SQL_PORT"                              required:"false" default:"1433"`
	SQLUser                 string        `envconfig:"SQL_USER"                              required:"false" default:"malibu_app"`
	SQLPassword             string        `envconfig:"SQL_PASSWORD"                          required:"false" default:"712Hgk9a99XU"`
	SQLDatabase             string        `envconfig:"SQL_DATABASE"                          required:"false" default:"Malibu"`
	TaskTemplateName        string        `envconfig:"TASK_GET_RECOMMENDATIONS"              required:"false" default:"TaskGetRecommendations"`
	TaskTemplateCronPattern string        `envconfig:"TASK_GET_RECOMMENDATIONS_CRON_PATTERN" required:"false" default:"0 1 * * * *"`
	TaskTemplateIsActive    bool          `envconfig:"TASK_GET_RECOMMENDATIONS_IS_ACTIVE"    required:"false" default:"false"`
}

// New initialize configuration
func New() (*Configuration, error) {
	var cfg Configuration
	err := envconfig.Process(ServiceName, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
