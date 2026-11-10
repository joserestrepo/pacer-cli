package domain

type Budget struct {
	ActiveHours          int
	MaxParallel          int
	Cutoff               string
	MinFocusBlocks       int
	FocusBlockMinMinutes int
}
