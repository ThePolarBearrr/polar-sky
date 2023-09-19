package worker

import (
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"polar-sky/task"
)

type Worker struct {
	Name      string
	Queue     queue.Queue
	TaskCount int
	//TODO: what does Db do
	DB     map[uuid.UUID]task.Task
	Logger zap.Logger
}
