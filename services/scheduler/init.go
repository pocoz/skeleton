package scheduler

import (
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/robfig/cron"

	"github.com/dieqnt/skeleton/models"
)

// Engine ...
type Engine struct {
	logger      log.Logger
	activeTasks map[string]bool
	mu          sync.RWMutex
	cron        *cron.Cron
}

// New ...
func New(logger log.Logger) (*Engine, error) {
	engine := &Engine{
		logger:      logger,
		activeTasks: make(map[string]bool),
		mu:          sync.RWMutex{},
		cron:        cron.New(),
	}

	return engine, nil
}

func (s *Engine) createSafeCronTask(taskName string, targetFunc func() error) func() {
	return func() {
		s.mu.Lock()
		defer s.mu.Unlock()

		if _, isActive := s.activeTasks[taskName]; isActive {
			level.Warn(s.logger).Log("msg", "task is already running, the new one will not be launched", "task", taskName)
			return
		}

		level.Info(s.logger).Log("msg", "task start", "task", taskName)
		s.activeTasks[taskName] = true

		err := targetFunc()
		delete(s.activeTasks, taskName)
		if err != nil {
			level.Error(s.logger).Log("msg", "ask completed with error", "task", taskName, "err", err)
			return
		}

		level.Info(s.logger).Log("msg", "task completed successfully", "task", taskName)
		return
	}
}

// Run ...
func (s *Engine) Run(tasks []*models.Task) error {
	for _, task := range tasks {
		if task.IsActive {
			err := s.cron.AddFunc(task.CronPattern, s.createSafeCronTask(task.Name, task.Func))
			if err != nil {
				return err
			}

			level.Info(s.logger).Log("msg", "task added successfully", "name", task.Name, "pattern", task.CronPattern)
		}
	}

	s.cron.Start()
	return nil
}

// Shutdown ...
func (s *Engine) Shutdown() {
	s.cron.Stop()
}
