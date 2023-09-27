package tests

import (
	"fmt"
	"github.com/docker/docker/client"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"polar-sky/log"
	"polar-sky/manager"
	"polar-sky/task"
	"polar-sky/worker"
	"testing"
)

func TestAPI(t *testing.T) {
	testLogger, _ := zap.NewDevelopment()
	logger := testLogger.Sugar()
	log.Logger = *logger
	//host := os.Getenv("SKY_HOST")
	//port, _ := strconv.Atoi(os.Getenv("SKY_PORT"))
	host := "127.0.0.1"
	port := 34567
	mHost := "127.0.0.1"
	mPort := 34568
	dc, _ := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.41"))
	task.Client = *dc

	w := worker.Worker{
		DB:     make(map[uuid.UUID]*task.Task),
		Queue:  *queue.New(),
		Logger: *logger,
	}
	api := worker.API{
		Address: host,
		Port:    port,
		Worker:  &w,
	}

	go w.RunTasks()
	//go w.CollectStats()
	go api.Start()

	workers := []string{fmt.Sprintf("%s:%d", host, port)}
	m := manager.New(workers)

	//for i := 0; i < 3; i++ {
	//	tsk := task.Task{
	//		ID:    uuid.New(),
	//		Name:  fmt.Sprintf("test-task-%d", i),
	//		Image: "strm/helloworld-http",
	//		State: task.Scheduled,
	//	}
	//	te := task.TaskEvent{
	//		ID:    uuid.New(),
	//		State: task.Running,
	//		Task:  tsk,
	//	}
	//	m.AddTaskEvent(te)
	//	m.SendTask()
	//}

	//go func() {
	//	for {
	//		log.Logger.Infof("[Manager] Updating tasks from %d workers\n", len(m.Workers))
	//		m.UpdateTasks()
	//		time.Sleep(time.Second * 15)
	//	}
	//}()
	//
	//for {
	//	for _, tsk := range m.TaskDB {
	//		log.Logger.Infof("[Manager] TaskID: %s, state: %d", tsk.ID, tsk.State)
	//		time.Sleep(time.Second * 15)
	//	}
	//}

	go m.ProcessTasks()
	go m.UpdateTasks()
	mapi := manager.API{Address: mHost, Port: mPort, Manager: m}
	mapi.Start()
}
