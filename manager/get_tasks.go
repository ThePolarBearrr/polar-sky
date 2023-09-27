package manager

import "polar-sky/task"

func (m *Manager) GetTasks() []task.Task {
	var tasks []task.Task
	for _, tsk := range m.TaskDB {
		tasks = append(tasks, *tsk)
	}
	return tasks
}
