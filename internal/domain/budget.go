package domain

// Budget is the daily work budget the user wants to respect.
type Budget struct {
	ActiveHours          int    // ActiveHours is the number of hours the user wants to work each day.
	MaxParallel          int    // MaxParallel is the number of sessions the user can run in parallel.
	Cutoff               string // Cutoff is the time the user wants to stop working in format "HH:MM".
	MinFocusBlocks       int    // MinFocusBlocks is the number of focus blocks the user wants to have each day.
	FocusBlockMinMinutes int    // FocusBlockMinMinutes is the minimum number of minutes the user wants to spend on a focus block.
}
