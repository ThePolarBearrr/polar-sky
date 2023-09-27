package worker

import (
	"errors"
	"fmt"
	"polar-sky/log"
	"polar-sky/task"
	"time"
)

func (w *Worker) runTask() task.DockerResult {
	t := w.Queue.Dequeue()
	if t == nil {
		w.Logger.Infof("No task need to run")
		return task.DockerResult{Error: nil}
	}

	taskQueued := t.(task.Task)
	taskPersisted, ok := w.DB[taskQueued.ID]
	if !ok {
		taskPersisted = &taskQueued
		w.DB[taskQueued.ID] = &taskQueued
	}

	var result task.DockerResult
	if task.ValidStateTransition(taskPersisted.State, taskQueued.State) {
		switch taskQueued.State {
		case task.Scheduled:
			result = w.StartTask(&taskQueued)
		case task.Completed:
			result = w.StopTask(&taskQueued)
		default:
			w.Logger.Errorf("Unexpected state %s when running task %s", taskQueued.State, taskQueued.Name)
			result.Error = errors.New("unexpected state when running task")
		}
	} else {
		result.Error = fmt.Errorf("invalid transition from %s to %s", taskPersisted.State, taskQueued.State)
	}

	return result
}

func (w *Worker) RunTasks() {
	for {
		if w.Queue.Len() != 0 {
			result := w.runTask()
			if result.Error != nil {
				log.Logger.Errorf("Failed run task, error: %v\n", result.Error)
			}
		} else {
			log.Logger.Info("No task need to run")
		}
		time.Sleep(time.Second * 10)
	}
}
