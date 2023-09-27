package task

var stateTransitionMap = map[State][]State{
	Pending:   {Scheduled},
	Scheduled: {Running, Scheduled, Failed},
	Running:   {Failed, Completed},
	Completed: {},
	Failed:    {},
}

func Contains(states []State, state State) bool {
	for _, s := range states {
		if s == state {
			return true
		}
	}
	return false
}

func ValidStateTransition(src State, dst State) bool {
	return Contains(stateTransitionMap[src], dst)
}
