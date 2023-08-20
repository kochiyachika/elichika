package enum

const (
	LiveFinishStatusCleared  = 1
	LiveFinishStatusFailure  = 11
	LiveFinishStatusUserExit = 12
)

var (
	// beginner, intermediate, advanced, expert (advanced+), challenge
	// 40 is also used in network (not in db), but not sure what it represent
	LiveDifficultyTypes = [5]int{10, 20, 30, 35, 37}
)
