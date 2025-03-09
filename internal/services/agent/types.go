package agent

import (
	"github.com/jahvon/flow/types/executable"
)

type State string

const (
	StateUnknown State = "unknown"
	StateRunning State = "running"
	StateStopped State = "stopped"
)

type ExecutableTask struct {
	Ref     executable.Ref
	Args    []string
	EnvVars map[string]string
	ID      string
	Status  TaskStatus
}

type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusFailed    TaskStatus = "failed"
)

type ScheduledTask struct {
	Task     ExecutableTask
	Schedule string
	Enabled  bool
}

type Options struct {
	WorkDir     string
	LogFile     string
	DisplayName string
	Description string
}
