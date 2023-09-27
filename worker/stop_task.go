package worker

import (
	"polar-sky/task"
	"time"
)

func (w *Worker) StopTask(t *task.Task) task.DockerResult {
	w.Logger.Infof("Starting to stop task %s", t.Name)

	cfg := task.NewConfig(t)
	d := task.NewDocker(&cfg)

	result := d.Stop(t.ContainerID)
	if result.Error != nil {
		w.Logger.Errorf("Falied stop container %s, error: %v\n", t.ContainerID, result.Error)
		return result
	}

	t.FinishTime = time.Now().UTC()
	t.State = task.Completed
	w.DB[t.ID] = t
	w.Logger.Infof("Stopped and removed container %s for task %s\n", t.ContainerID, t.Name)

	return result
}
