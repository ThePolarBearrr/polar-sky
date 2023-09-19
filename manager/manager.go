package manager

import (
	"github.com/google/uuid"
	"go.uber.org/zap"
	"polar-sky/task"
	"polar-sky/worker"
)

type Manager struct {
	//TODO: what use
	TaskDB      map[string][]task.Task
	TaskEventDB map[string][]task.TaskEvent

	Workers       []string
	WorkerTaskMap map[string][]uuid.UUID
	TaskWorkerMap map[uuid.UUID]worker.Worker

	Logger zap.Logger
}
