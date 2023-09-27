package worker

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"net/http"
	"polar-sky/log"
	"polar-sky/task"
)

func (a *API) StartTaskHandler(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	te := task.TaskEvent{}
	err := d.Decode(&te)
	if err != nil {
		msg := fmt.Sprintf("Failed decode request body, error: %v\n", err)
		log.Logger.Errorf(msg)
		w.WriteHeader(400)
		e := ErrResponse{
			HTTPStatusCode: 400,
			Message:        msg,
		}
		json.NewEncoder(w).Encode(e)
		return
	}

	a.Worker.AddTask(te.Task)
	log.Logger.Infof("Added task %v\n", te.Task.Name)
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(te.Task)
	return
}

func (a *API) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(a.Worker.GetTasks())
}

func (a *API) StopTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskID")
	if taskID == "" {
		log.Logger.Errorf("No taskID passed in request")
		w.WriteHeader(400)
	}
	uID, _ := uuid.Parse(taskID)
	tsk, ok := a.Worker.DB[uID]
	if !ok {
		log.Logger.Errorf("No task %s in worker", taskID)
		w.WriteHeader(404)
	}

	taskCopy := *tsk
	taskCopy.State = task.Completed
	a.Worker.AddTask(taskCopy)
	log.Logger.Infof("Added task %s to stop container %s", taskCopy.ID, taskCopy.ContainerID)
	w.WriteHeader(204)
}

func (a *API) GetStatsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(a.Worker.Stats)
}
