package manager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"polar-sky/log"
	"polar-sky/task"
	"polar-sky/worker"
	"time"
)

func (m *Manager) SendTask() {
	if m.Pending.Len() > 0 {
		w := m.SelectWorker()

		t := m.Pending.Dequeue()
		te := t.(task.TaskEvent)
		tsk := te.Task
		log.Logger.Infof("Send task %s from pending queue", tsk.ID)
		tsk.State = task.Scheduled

		m.TaskEventDB[te.ID] = &te
		m.TaskDB[tsk.ID] = &tsk
		m.TaskWorkerMap[tsk.ID] = w
		m.WorkerTaskMap[w] = append(m.WorkerTaskMap[w], tsk.ID)

		data, err := json.Marshal(te)
		if err != nil {
			log.Logger.Errorf("Failed marshal tsk %v to json, error: %v", tsk, err)
			m.Pending.Enqueue(tsk)
			return
		}

		url := fmt.Sprintf("http://%s/tasks", w)
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
		if err != nil {
			log.Logger.Errorf("Send http request to %s failed, error: %v\n", url, err)
			m.Pending.Enqueue(tsk)
			return
		}

		d := json.NewDecoder(resp.Body)
		if resp.StatusCode != http.StatusCreated {
			e := worker.ErrResponse{}
			err = d.Decode(&e)
			if err != nil {
				log.Logger.Errorf("Decode json body failed, error: %v\n", err)
			}
			log.Logger.Errorf("Response error (%d): %s\n", e.HTTPStatusCode, e.Message)
			return
		}

		tsk = task.Task{}
		err = d.Decode(&tsk)
		if err != nil {
			log.Logger.Errorf("Decode json body failed, error: %v\n", err)
		}
		log.Logger.Infof("Succeed send task %v to worker %s", tsk, w)
	} else {
		log.Logger.Infof("No task in pending queue")
	}
}

func (m *Manager) ProcessTasks() {
	for {
		log.Logger.Infof("Processing any tasks in the queue")
		m.SendTask()
		time.Sleep(time.Second * 10)
	}
}
