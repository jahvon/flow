package agent

import (
	"sync"

	"github.com/jahvon/tuikit/io"
	"github.com/kardianos/service"

	"github.com/jahvon/flow/internal/cache"
)

type Agent interface {
	Run()
	Start(s service.Service) error
	Stop(s service.Service) error
	Status() (State, error)
	Enqueue(task ExecutableTask) error
	Schedule(task ExecutableTask, schedule string) error
	Unschedule(taskID string) error
	ListScheduled() ([]ScheduledTask, error)
	GetScheduled(taskID string) (*ScheduledTask, error)
	EnableScheduled(taskID string) error
	DisableScheduled(taskID string) error
}

type agent struct {
	service   service.Service
	options   *Options
	logger    io.Logger
	execCache cache.ExecutableCache

	queue chan string
	mu    sync.RWMutex
}

func NewAgent(logger io.Logger, execCache cache.ExecutableCache, options *Options) (Agent, error) {
	svcCfg := &service.Config{Name: "flow"}
	a := &agent{
		options:   options,
		logger:    logger,
		execCache: execCache,
		queue:     make(chan string),
		mu:        sync.RWMutex{},
	}
	svc, err := service.New(a, svcCfg)
	if err != nil {
		return nil, err
	}
	a.service = svc
	return a, nil
}

func (a *agent) Start(s service.Service) error {
	// This is called by kardianos/service when the agent is started
	// Implement the initialization logic here
	go func() {
		for {
			select {
			case task := <-a.queue:
				a.logger.Infof("Hello task %s", task)
			}
		}
	}()
	return nil
}

func (a *agent) Run() {
	_ = a.Start(a.service)
}

func (a *agent) Stop(s service.Service) error {
	// This is called by kardianos/service when the agent is stopped
	// Implement the cleanup logic here
	return nil
}

func (a *agent) Status() (State, error) {
	// Implement status checking logic
	return StateUnknown, nil
}

func (a *agent) Enqueue(task ExecutableTask) error {
	// Implement task queuing logic

	return nil
}

func (a *agent) Schedule(task ExecutableTask, schedule string) error {
	// Implement task scheduling logic
	return nil
}

func (a *agent) Unschedule(taskID string) error {
	// Implement task unscheduling logic
	return nil
}

func (a *agent) ListScheduled() ([]ScheduledTask, error) {
	// Implement listing scheduled tasks
	return nil, nil
}

func (a *agent) GetScheduled(taskID string) (*ScheduledTask, error) {
	// Implement getting a specific scheduled task
	return nil, nil
}

func (a *agent) EnableScheduled(taskID string) error {
	// Implement enabling a scheduled task
	return nil
}

func (a *agent) DisableScheduled(taskID string) error {
	// Implement disabling a scheduled task
	return nil
}
