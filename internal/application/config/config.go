package config

import (
	"time"

	"github.com/joserestrepo/pacer-cli/internal/domain"
)

// Config is the configuration for the application.
type Config struct {
	General General
	Budget  domain.Budget
	Checkin Checkin
	Alerts  Alerts
}

// Checkin is the configuration for the check-in feature.
type Checkin struct {
	Reminder      bool
	ReminderAt    string
	StartReminder bool
	BadDayMin     int
}

// Alerts is the configuration for the alerts feature.
type Alerts struct {
	Enabled  bool
	Cooldown time.Duration
}

// General is the additional configuration for the application.
type General struct {
	DayStartsAt string // DayStartsAt is the time ("HH:MM") pacer starts a new day; work before it counts toward the previous day.
	Language    string // Language is the language the user wants to use for the application.
}

// Defaults returns the values for the config by default. It is used when the user sets nothing in the config file.
func Defaults() Config {
	return Config{
		General: General{
			DayStartsAt: "04:00",
			Language:    "auto", // "auto" means "use the system language"
		},
		Budget: domain.Budget{
			ActiveHours:          6,
			MaxParallel:          2,
			Cutoff:               "19:00",
			MinFocusBlocks:       2,
			FocusBlockMinMinutes: 25,
		},
		Checkin: Checkin{
			Reminder:      true,
			ReminderAt:    "", // empty means "use the budget cutoff time"
			StartReminder: true,
			BadDayMin:     4,
		},
		Alerts: Alerts{
			Enabled:  true,
			Cooldown: 60 * time.Minute,
		},
	}
}
