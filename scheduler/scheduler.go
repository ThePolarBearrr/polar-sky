package scheduler

type Scheduler interface {
	SelectCandidateNodes()
	Pick()
	Score()
}
