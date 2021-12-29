package scrollsvc

import (
	"context"
	"sync"

	"github.com/go-kit/kit/log"

	"github.com/dieqnt/skeleton/db/elasticsearch"
	"github.com/dieqnt/skeleton/models"
)

type Scroller struct {
	logger  log.Logger
	dbES    *elasticsearch.Engine
	ctx     context.Context
	eofChan chan bool
	errChan chan error
	wg      *sync.WaitGroup
	Tasks   []*models.Task
}

// Config is a transporter configuration
type Config struct {
	Logger      log.Logger
	DBES        *elasticsearch.Engine
	Context     context.Context
	ConfigTasks *ConfigTasks
}

// ConfigTasks ...
type ConfigTasks struct {
	TaskTemplateName        string
	TaskTemplateCronPattern string
	TaskTemplateIsActive    bool
}

func New(cfg *Config) (*Scroller, error) {
	t := &Scroller{
		logger:  cfg.Logger,
		dbES:    cfg.DBES,
		ctx:     cfg.Context,
		errChan: make(chan error, 1),
		eofChan: make(chan bool, 1),
		wg:      &sync.WaitGroup{},
		Tasks:   make([]*models.Task, 0),
	}

	taskTemplate := &models.Task{
		Name:        cfg.ConfigTasks.TaskTemplateName,
		CronPattern: cfg.ConfigTasks.TaskTemplateCronPattern,
		IsActive:    cfg.ConfigTasks.TaskTemplateIsActive,
		Func:        t.Run,
	}
	t.Tasks = append(t.Tasks, taskTemplate)

	return t, nil
}

func (s *Scroller) Run() error {
	s.wg.Add(1)
	go s.processBuffer()

	s.wg.Wait()

	return nil
}
