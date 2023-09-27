package manager

import (
	"encoding/json"
	"fmt"
	"net/http"
	"polar-sky/log"
	"polar-sky/task"
	"time"
)

func (m *Manager) updateTasks() {
	for _, worker := range m.Workers {
		log.Logger.Infof("Start to update tasks from worker %s", worker)
		url := fmt.Sprintf("http://%s/tasks", worker)
		resp, err := http.Get(url)
		if err != nil {
			log.Logger.Errorf("Send http request to worker %s failed, error: %v\n", worker, err)
			continue
		}
		d := json.NewDecoder(resp.Body)
		var tasks []*task.Task
		err = d.Decode(&tasks)
		if err != nil {
			log.Logger.Errorf("Decode resp body to tasks failed, error: %v\n", err)
			continue
		}

		for _, t := range tasks {
			log.Logger.Infof("Attempting to update task %s", t.ID)

			tsk, ok := m.TaskDB[t.ID]
			if !ok {
				log.Logger.Errorf("Task %s not found in manager db", t.ID)
				continue
			}

			if tsk.State != t.State {
				tsk.State = t.State
			}
			tsk.StartTime = t.StartTime
			tsk.FinishTime = t.FinishTime
			tsk.ContainerID = t.ContainerID
		}
	}
}

func (m *Manager) UpdateTasks() {
	for {
		log.Logger.Infof("Checking for task updates from workers")
		m.updateTasks()
		time.Sleep(time.Second * 15)
	}
}
