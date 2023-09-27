package worker

import "polar-sky/task"

func (w *Worker) StartTask(t *task.Task) task.DockerResult {
	cfg := task.NewConfig(t)
	d := task.NewDocker(&cfg)

	result := d.Run()
	if result.Error != nil {
		w.Logger.Errorf("Failed start task %s, error: %v\n", t.Name, result.Error)
		t.State = task.Failed
		w.DB[t.ID] = t
		return result
	}

	t.ContainerID = result.ContainerID
	//TODO: is pointer need to save?
	w.DB[t.ID] = t
	t.State = task.Running

	return result
}
