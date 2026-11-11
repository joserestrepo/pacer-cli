package domain

// Budget is the daily work budget the user wants to respect. pacer compares
// each day's activity against it to decide whether to alert.
type Budget struct {
	// ActiveHours is the maximum number of hours of AI-assisted work per day.
	ActiveHours int
	// MaxParallel is the maximum number of Claude Code sessions active at the
	// same time.
	MaxParallel int
	// Cutoff is the local time ("HH:MM") after which work counts as late.
	Cutoff string
	// MinFocusBlocks is the minimum number of focus blocks the user wants
	// each day.
	MinFocusBlocks int
	// FocusBlockMinMinutes is how long, in minutes, a stretch of single-session
	// work must last to count as a focus block.
	FocusBlockMinMinutes int
}
