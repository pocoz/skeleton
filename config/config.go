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
	ElasticServer            string        `envconfig:"ELASTIC_SERVER"                        required:"false" default:"http://elastic-nginx.test.cloud.croc-comp.goods.local:19200"`
	ElasticUser              string        `envconfig:"ELASTIC_USER"                          required:"false" default:""`
	ElasticPassword          string        `envconfig:"ELASTIC_PASSWORD"                      required:"false" default:""`
	SQLServer                string        `envconfig:"SQL_SERVER"                            required:"false" default:"mlb-sdb-002.corp.goods.ru"`
	SQLPort                  string        `envconfig:"SQL_PORT"                              required:"false" default:"1433"`
	SQLUser                  string        `envconfig:"SQL_USER"                              required:"false" default:"malibu_app"`
	SQLPassword              string        `envconfig:"SQL_PASSWORD"                          required:"false" default:"712Hgk9a99XU"`
	SQLDatabase              string        `envconfig:"SQL_DATABASE"                          required:"false" default:"Malibu"`
	PostgresServer           string        `envconfig:"POSTGRES_SERVER"                                        default:"advertiser-postgres-haproxy.test.cloud.croc-comp.goods.local"`
	PostgresPort             string        `envconfig:"POSTGRES_PORT"                                          default:"6432"`
	PostgresUser             string        `envconfig:"POSTGRES_USER"                         required:"true"`
	PostgresPassword         string        `envconfig:"POSTGRES_PASSWORD"                     required:"true"`
	PostgresDatabase         string        `envconfig:"POSTGRES_DATABASE"                                      default:"advertiser"`
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
