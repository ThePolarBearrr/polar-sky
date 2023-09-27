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
	DB        map[uuid.UUID]*task.Task
	Logger    zap.SugaredLogger
	Stats     *Stats
}

func (w *Worker) AddTask(t task.Task) {
	w.Queue.Enqueue(t)
}

func (w *Worker) GetTasks() []task.Task {
	tasks := make([]task.Task, 0)
	for _, v := range w.DB {
		tasks = append(tasks, *v)
	}
	return tasks
}
