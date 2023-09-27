package manager

func (m *Manager) SelectWorker() string {
	var newWorker int
	if m.LastWorker+1 < len(m.Workers) {
		newWorker = m.LastWorker + 1
		m.LastWorker = newWorker
	} else {
		newWorker = 0
		m.LastWorker = newWorker
	}
	return m.Workers[newWorker]
}
