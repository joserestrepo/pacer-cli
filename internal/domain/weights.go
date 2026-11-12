package domain

// Weights sets how much each budget signal counts toward the daily score.
// A weight of 1 is the default; 0 removes the signal from the score.
type Weights struct {
	// Hours weighs going over the budget's active hours.
	Hours float64
	// Parallel weighs running more sessions at once than the budget allows.
	Parallel float64
	// Late weighs working after the budget's cutoff.
	Late float64
	// Focus weighs having fewer focus blocks than the budget asks for.
	Focus float64
}
