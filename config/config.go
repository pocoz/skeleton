package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

const ServiceName = "SKELETON"

// Configuration describe services config
type Configuration struct {
	HTTPPort                 string        `envconfig:"HTTP_PORT"                             required:"false" default:"2524"`
	RateLimitEvery           time.Duration `envconfig:"RATE_LIMIT_EVERY"                      required:"false" default:"1us"`
	RateLimitBurst           int           `envconfig:"RATE_LIMIT_BURST"                      required:"false" default:"100"`
	ElasticServer            string        `envconfig:"ELASTIC_SERVER"                        required:"false" default:""`
	ElasticUser              string        `envconfig:"ELASTIC_USER"                          required:"false" default:""`
	ElasticPassword          string        `envconfig:"ELASTIC_PASSWORD"                      required:"false" default:""`
	SQLServer                string        `envconfig:"SQL_SERVER"                            required:"false" default:""`
	SQLPort                  string        `envconfig:"SQL_PORT"                              required:"false" default:"1433"`
	SQLUser                  string        `envconfig:"SQL_USER"                              required:"false" default:""`
	SQLPassword              string        `envconfig:"SQL_PASSWORD"                          required:"false" default:""`
	SQLDatabase              string        `envconfig:"SQL_DATABASE"                          required:"false" default:""`
	PostgresServer           string        `envconfig:"POSTGRES_SERVER"                                        default:""`
	PostgresPort             string        `envconfig:"POSTGRES_PORT"                                          default:"5432"`
	PostgresUser             string        `envconfig:"POSTGRES_USER"                         required:"true"`
	PostgresPassword         string        `envconfig:"POSTGRES_PASSWORD"                     required:"true"`
	PostgresDatabase         string        `envconfig:"POSTGRES_DATABASE"                                      default:""`
	PostgresMigrationsFolder string        `envconfig:"POSTGRES_MIGRATIONS_FOLDER"                             default:"db/postgres/migrations/skeleton"`
	PostgresMigrationsTable  string        `envconfig:"POSTGRES_MIGRATIONS_TABLE"                              default:"migrations"`
	PostgresMaxButchSize     int           `envconfig:"POSTGRES_MAX_BUTCH_SIZE"                                default:"100"`
	TaskTemplateName         string        `envconfig:"TASK_GET_RECOMMENDATIONS"              required:"false" default:"TaskGetRecommendations"`
	TaskTemplateCronPattern  string        `envconfig:"TASK_GET_RECOMMENDATIONS_CRON_PATTERN" required:"false" default:"0 1 * * * *"`
	TaskTemplateIsActive     bool          `envconfig:"TASK_GET_RECOMMENDATIONS_IS_ACTIVE"    required:"false" default:"false"`
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
