package template

import (
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/pocoz/skeleton/models"
)

// loggingMiddleware wraps Service and logs request information to the provided Logger.
type loggingMiddleware struct {
	next   service
	logger log.Logger
}

func (m *loggingMiddleware) template(options *models.TemplateRequest) (*models.TemplateResponse, error) {
	begin := time.Now()
	res, err := m.next.template(options)
	level.Info(m.logger).Log(
		"method", "template",
		"elapsed", time.Since(begin),
		"err", err,
	)
	return res, err
}
