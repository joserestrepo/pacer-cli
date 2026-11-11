package config

import (
	"testing"
	"time"
)

// Spec: config-management, "Requirement: Defaults without file".
// Scenarios: "Bool default true" and "New defaults" (budget, checkin and alerts part).
func TestDefaults(t *testing.T) {
	c := Defaults()

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"budget.active_hours", c.Budget.ActiveHours, 6},
		{"budget.max_parallel", c.Budget.MaxParallel, 2},
		{"budget.cutoff", c.Budget.Cutoff, "19:00"},
		{"budget.min_focus_blocks", c.Budget.MinFocusBlocks, 2},
		{"budget.focus_block_min_minutes", c.Budget.FocusBlockMinMinutes, 25},
		{"checkin.reminder", c.Checkin.Reminder, true},
		{"checkin.reminder_at", c.Checkin.ReminderAt, ""},
		{"checkin.start_reminder", c.Checkin.StartReminder, true},
		{"checkin.bad_day_min", c.Checkin.BadDayMin, 4},
		{"alerts.enabled", c.Alerts.Enabled, true},
		{"alerts.cooldown", c.Alerts.Cooldown, 60 * time.Minute},
		{"general.day_starts_at", c.General.DayStartsAt, "04:00"},
		{"general.language", c.General.Language, "auto"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v (%T), want %v (%T)", tt.got, tt.got, tt.want, tt.want)
			}
		})
	}
}
