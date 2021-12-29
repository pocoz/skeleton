package memlogger

import (
	"runtime"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type Config struct {
	Logger log.Logger
}

type Engine struct {
	logger log.Logger
	ticker *time.Ticker
}

func New(cfg *Config) *Engine {
	return &Engine{
		logger: cfg.Logger,
		ticker: time.NewTicker(1 * time.Minute),
	}
}

func (e *Engine) Run() {
	for {
		<-e.ticker.C
		e.show()
	}
}

func (e *Engine) TickerReset(duration time.Duration) {
	e.ticker.Reset(duration)
}

// show outputs the current, total and OS memory being used.
// As well as the number of garage collection cycles completed.
func (e *Engine) show() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	level.Info(e.logger).Log(
		"msg", "memory use info",
		"alloc", BytesToHumanReadableForm(int(m.Alloc), 0),
		"total_alloc", BytesToHumanReadableForm(int(m.TotalAlloc), 0),
		"sys", BytesToHumanReadableForm(int(m.Sys), 0),
		"num_gc", m.NumGC,
	)
}
