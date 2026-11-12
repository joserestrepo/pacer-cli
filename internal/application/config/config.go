// Package config defines pacer's configuration contract: every setting, its
// default value and its key. Reading and writing the TOML file lives in
// infrastructure/configfile.
package config

import (
	"time"

	"github.com/joserestrepo/pacer-cli/internal/domain"
)

// Config holds every user-configurable setting of pacer. Each field is one
// section of the TOML file, in the same order.
type Config struct {
	General    General
	Budget     domain.Budget
	Alerts     Alerts
	Checkin    Checkin
	Ingest     Ingest
	Scoring    Scoring
	Validation Validation
}

// General holds settings that affect the whole application rather than a
// single feature.
type General struct {
	// DayStartsAt is the local time ("HH:MM") at which pacer starts a new day;
	// work before it counts toward the previous day.
	DayStartsAt string
	// Language is the language of user-facing text: "auto", "en" or "es".
	// "auto" picks it from LC_ALL, LC_MESSAGES or LANG.
	Language string
}

// Alerts controls when and how pacer warns the user about their work rhythm.
type Alerts struct {
	// Enabled turns alert evaluation on or off. Activity is still recorded
	// when it is off.
	Enabled bool
	// Cooldown is the minimum time before the same alert rule fires again.
	Cooldown time.Duration
	// Mode is how pacer intervenes. Only "notify" (a desktop notification)
	// exists today; the opt-in soft brake arrives with pacer-integrations.
	Mode string
	// ScoreThreshold is the daily score at or above which pacer alerts, even
	// if no single budget limit was crossed.
	ScoreThreshold float64
}

// Checkin controls the daily fatigue check-ins (1 fresh to 5 drained) that
// calibrate the score.
type Checkin struct {
	// Reminder turns the end-of-day check-in reminder on or off.
	Reminder bool
	// ReminderAt is the local time ("HH:MM") of the end-of-day reminder.
	// Empty means "use the budget cutoff".
	ReminderAt string
	// StartReminder turns the start-of-day reminder on or off. It fires on
	// the first Claude Code activity of the day.
	StartReminder bool
	// BadDayMin is the lowest end-of-day answer that counts as a bad day.
	BadDayMin int
	// Actions shows answer buttons on the reminder notification, so the user
	// can check in without opening a terminal.
	Actions bool
}

// Ingest controls how pacer reads Claude Code's session logs.
type Ingest struct {
	// PollInterval is how often the log files are checked for new lines, as
	// a fallback when file-change events are missed.
	PollInterval time.Duration
	// IdleGap is the longest pause that still counts as continuous work.
	IdleGap time.Duration
	// CountSubagentsAsParallel makes subagent sessions count toward the
	// number of parallel sessions.
	CountSubagentsAsParallel bool
}

// Scoring controls how pacer turns the day's activity into a score.
type Scoring struct {
	// EvaluateEvery is how often the daemon recomputes the score and decides
	// whether to alert.
	EvaluateEvery time.Duration
	// BaselineDays is how many past days form the personal baseline that
	// today is compared against.
	BaselineDays int
	// MinBaselineDays is the minimum number of active days in that window
	// before the baseline influences the score; below it, only the budget
	// decides.
	MinBaselineDays int
	// Weights sets how much each budget signal counts toward the score.
	Weights domain.Weights
}

// Validation controls the metrics that check whether pacer's score matches
// how tired the user reports feeling.
type Validation struct {
	// MinPairs is the minimum number of days with both a score and a
	// check-in needed before a correlation is reported.
	MinPairs int
}

// Defaults returns the configuration pacer uses for every key the user does
// not set: the TOML file is decoded on top of it.
func Defaults() Config {
	return Config{
		General: General{
			DayStartsAt: "04:00",
			Language:    "auto",
		},
		Budget: domain.Budget{
			ActiveHours:          6,
			MaxParallel:          2,
			Cutoff:               "19:00",
			MinFocusBlocks:       2,
			FocusBlockMinMinutes: 25,
		},
		Alerts: Alerts{
			Enabled:        true,
			Cooldown:       60 * time.Minute,
			Mode:           "notify",
			ScoreThreshold: 0.7,
		},
		Checkin: Checkin{
			Reminder:      true,
			ReminderAt:    "", // empty means "use the budget cutoff"
			StartReminder: true,
			BadDayMin:     4,
			Actions:       true,
		},
		Ingest: Ingest{
			PollInterval:             10 * time.Second,
			IdleGap:                  5 * time.Minute,
			CountSubagentsAsParallel: false,
		},
		Scoring: Scoring{
			EvaluateEvery:   60 * time.Second,
			BaselineDays:    28,
			MinBaselineDays: 7,
			Weights: domain.Weights{
				Hours:    1.0,
				Parallel: 1.0,
				Late:     1.0,
				Focus:    1.0,
			},
		},
		Validation: Validation{
			MinPairs: 10,
		},
	}
}
