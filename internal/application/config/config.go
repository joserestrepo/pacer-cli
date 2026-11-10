package config

import (
	"time"

	"github.com/joserestrepo/pacer-cli/internal/domain"
)

type Checkin struct {
	Reminder      bool
	ReminderAt    string
	StartReminder bool
	BadDayMin     int
}

type Alerts struct {
	Enabled  bool
	Cooldown time.Duration
}

type Config struct {
	Budget  domain.Budget
	Checkin Checkin
	Alerts  Alerts
}

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
			ReminderAt:    "",
			StartReminder: true,
			BadDayMin:     4,
		},
		Alerts: Alerts{
			Enabled:  true,
			Cooldown: 60 * time.Minute,
		},
	}
}
