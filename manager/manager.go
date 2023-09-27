package manager

import (
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"polar-sky/task"
)

type Manager struct {
	LastWorker int
	Pending    queue.Queue

	TaskDB      map[uuid.UUID]*task.Task
	TaskEventDB map[uuid.UUID]*task.TaskEvent

	Workers       []string
	WorkerTaskMap map[string][]uuid.UUID
	TaskWorkerMap map[uuid.UUID]string
}

func (m *Manager) AddTaskEvent(te task.TaskEvent) {
	m.Pending.Enqueue(te)
}

func New(workers []string) *Manager {
	taskDB := make(map[uuid.UUID]*task.Task)
	taskEventDB := make(map[uuid.UUID]*task.TaskEvent)
	workerTaskMap := make(map[string][]uuid.UUID)
	taskWorkerMap := make(map[uuid.UUID]string)
	for _, worker := range workers {
		workerTaskMap[worker] = []uuid.UUID{}
	}
	return &Manager{
		TaskDB:        taskDB,
		TaskEventDB:   taskEventDB,
		Workers:       workers,
		WorkerTaskMap: workerTaskMap,
		TaskWorkerMap: taskWorkerMap,
	}
}
