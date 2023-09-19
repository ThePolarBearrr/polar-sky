package task

import (
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	"time"
)

type State int

const (
	Pending State = iota
	Scheduled
	Running
	Completed
	Failed
)

type Task struct {
	ID            uuid.UUID
	Name          string
	State         State
	Image         string
	Memory        int
	Disk          int
	ExposedPorts  nat.PortSet
	PortBindings  map[string]string
	RestartPolicy string
	StartTime     time.Time
	FinishTime    time.Time
}

// TODO: change state define to express diff taskevent
type TaskEvent struct {
	ID        uuid.UUID
	State     State // which state the task should transition to
	TimeStamp time.Time
	Task      Task
}

type Config struct {
	Name          string
	Env           []string
	RestartPolicy string
	Cmd           string
	Image         string
	Memory        int64
	Disk          int64
	AttachStdin   bool
	AttachStdout  bool
	AttachStderr  bool
}
