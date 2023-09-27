package tests

import (
	"fmt"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"polar-sky/manager"
	"polar-sky/node"
	"polar-sky/task"
	"polar-sky/worker"
	"testing"
)

func TestSkeleton(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	testLogger := logger.Sugar()

	tsk := &task.Task{
		ID:     uuid.New(),
		Name:   "task-0",
		State:  task.Pending,
		Image:  "test-image",
		Memory: 16,
		Disk:   1024,
	}

	taskEvent := &task.TaskEvent{
		ID:    uuid.New(),
		State: task.Running,
		Task:  *tsk,
	}

	fmt.Printf("test task: %v\n", *tsk)
	fmt.Printf("test task event: %v\n", *taskEvent)

	wker := &worker.Worker{
		Name:      "test-worker",
		Queue:     *queue.New(),
		TaskCount: 1,
		DB:        make(map[uuid.UUID]*task.Task),
		Logger:    *testLogger,
	}

	wker.RunTasks()
	wker.StartTask(tsk)
	wker.CollectStats()
	wker.StopTask(tsk)

	m := &manager.Manager{
		TaskDB:      make(map[uuid.UUID]*task.Task),
		TaskEventDB: make(map[uuid.UUID]*task.TaskEvent),
		Workers:     []string{wker.Name},
	}

	m.SelectWorker()
	m.SendTask()
	m.UpdateTasks()

	nd := &node.Node{
		IP:              "192.168.0.0",
		Name:            "test-node",
		MemoryCapacity:  32,
		MemoryAllocated: 16,
		DiskCapacity:    2048,
		DiskAllocated:   1024,
		TaskCount:       1,
	}
	fmt.Printf("node: %v\n", *nd)
}
