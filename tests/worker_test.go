package tests

import (
	"github.com/docker/docker/client"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"polar-sky/log"
	"polar-sky/task"
	"polar-sky/worker"
	"testing"
)

func TestWorker(t *testing.T) {
	testLogger, _ := zap.NewDevelopment()
	logger := testLogger.Sugar()
	log.Logger = *logger
	dc, _ := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.41"))
	task.Client = *dc

	w := worker.Worker{
		Queue:  *queue.New(),
		DB:     make(map[uuid.UUID]*task.Task),
		Logger: *logger,
	}

	tsk := task.Task{
		ID:    uuid.New(),
		Name:  "test-task-1",
		Image: "strm/helloworld-http",
		State: task.Scheduled,
	}
	w.AddTask(tsk)

	//result := w.RunTask()
	//if result.Error != nil {
	//	panic(result.Error)
	//}
	//tsk.ContainerID = result.ContainerID
	//
	//log.Logger.Infof("task %s is running in container %s", tsk.Name, result.ContainerID)
	//time.Sleep(time.Second * 30)
	//
	//tsk.State = task.Completed
	//w.AddTask(tsk)
	//result = w.RunTask()
	//if result.Error != nil {
	//	panic(result.Error)
	//}
}
