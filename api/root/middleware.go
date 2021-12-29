package root

import (
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// loggingMiddleware wraps Service and logs request information to the provided Logger.
type loggingMiddleware struct {
	next   service
	logger log.Logger
}

func (m *loggingMiddleware) healthy() {
	begin := time.Now()
	m.next.healthy()
	level.Info(m.logger).Log(
		"method", "healthy",
		"elapsed", time.Since(begin),
	)
}
