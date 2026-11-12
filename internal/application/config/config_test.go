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
		{"alerts.mode", c.Alerts.Mode, "notify"},
		{"alerts.score_threshold", c.Alerts.ScoreThreshold, 0.7},
		{"checkin.actions", c.Checkin.Actions, true},
		{"ingest.poll_interval", c.Ingest.PollInterval, 10 * time.Second},
		{"ingest.idle_gap", c.Ingest.IdleGap, 5 * time.Minute},
		{"ingest.count_subagents_as_parallel", c.Ingest.CountSubagentsAsParallel, false},
		{"validation.min_pairs", c.Validation.MinPairs, 10},
		{"scoring.evaluate_every", c.Scoring.EvaluateEvery, 60 * time.Second},
		{"scoring.baseline_days", c.Scoring.BaselineDays, 28},
		{"scoring.min_baseline_days", c.Scoring.MinBaselineDays, 7},
		{"scoring.weights.hours", c.Scoring.Weights.Hours, 1.0},
		{"scoring.weights.parallel", c.Scoring.Weights.Parallel, 1.0},
		{"scoring.weights.late", c.Scoring.Weights.Late, 1.0},
		{"scoring.weights.focus", c.Scoring.Weights.Focus, 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v (%T), want %v (%T)", tt.got, tt.got, tt.want, tt.want)
			}
		})
	}
}
