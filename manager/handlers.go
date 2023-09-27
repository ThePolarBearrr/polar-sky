package manager

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"net/http"
	"polar-sky/log"
	"polar-sky/task"
	"time"
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

	a.Manager.AddTaskEvent(te)
	log.Logger.Infof("Added task event %v to pending queue\n", te.ID)
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(te.Task)
	return
}

func (a *API) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(a.Manager.GetTasks())
}

func (a *API) StopTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskID")
	if taskID == "" {
		log.Logger.Errorf("No taskID passed in request")
		w.WriteHeader(400)
	}
	uID, _ := uuid.Parse(taskID)
	tsk, ok := a.Manager.TaskDB[uID]
	if !ok {
		log.Logger.Errorf("No task %s in manager task db", taskID)
		w.WriteHeader(404)
	}

	te := task.TaskEvent{
		ID:        uuid.New(),
		State:     task.Completed,
		TimeStamp: time.Now(),
	}

	taskCopy := *tsk
	taskCopy.State = task.Completed
	te.Task = taskCopy
	a.Manager.AddTaskEvent(te)
	log.Logger.Infof("Added task event %s to pending queue", te.ID)
	w.WriteHeader(204)
}
