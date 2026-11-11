package config

import (
	"time"

	"github.com/joserestrepo/pacer-cli/internal/domain"
)

// Config is the configuration for the application.
type Config struct {
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

// Defaults returns the values for the config by default. It is used when the user sets nothing in the config file.
func Defaults() Config {
	return Config{
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
