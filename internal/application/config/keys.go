package config

// Kind is the type of value a Key holds. The TOML adapter uses it to
// parse and render the value.
type Kind int

// Supported value kinds.
const (
	KindBool     Kind = iota // true or false
	KindInt                  // whole number, e.g. 6
	KindString               // text, e.g. "19:00"
	KindDuration             // Go duration, e.g. "60m" or "1h30m"
	KindFloat                // decimal number, e.g. 0.7
)

// Key describes one setting the user can read or change by its path, as in
// `pacer config set budget.cutoff 18:00`. Get(Defaults()) is its default value.
type Key struct {
	// Path is the setting's location as "section.name" in snake_case,
	// matching the TOML file (e.g. "budget.active_hours").
	Path string
	// Kind is the type of the setting's value.
	Kind Kind
	// Get reads this setting's value from a Config.
	Get func(Config) any
}

// Keys lists every user-facing setting, in the same order as the TOML file.
// It is the single place that maps a path to a Config field.
var Keys = []Key{
	{Path: "general.day_starts_at", Kind: KindString, Get: func(c Config) any { return c.General.DayStartsAt }},
	{Path: "general.language", Kind: KindString, Get: func(c Config) any { return c.General.Language }},
	{Path: "budget.active_hours", Kind: KindInt, Get: func(c Config) any { return c.Budget.ActiveHours }},
	{Path: "budget.max_parallel", Kind: KindInt, Get: func(c Config) any { return c.Budget.MaxParallel }},
	{Path: "budget.cutoff", Kind: KindString, Get: func(c Config) any { return c.Budget.Cutoff }},
	{Path: "budget.min_focus_blocks", Kind: KindInt, Get: func(c Config) any { return c.Budget.MinFocusBlocks }},
	{Path: "budget.focus_block_min_minutes", Kind: KindInt, Get: func(c Config) any { return c.Budget.FocusBlockMinMinutes }},
	{Path: "alerts.enabled", Kind: KindBool, Get: func(c Config) any { return c.Alerts.Enabled }},
	{Path: "alerts.cooldown", Kind: KindDuration, Get: func(c Config) any { return c.Alerts.Cooldown }},
	{Path: "alerts.mode", Kind: KindString, Get: func(c Config) any { return c.Alerts.Mode }},
	{Path: "alerts.score_threshold", Kind: KindFloat, Get: func(c Config) any { return c.Alerts.ScoreThreshold }},
	{Path: "checkin.reminder", Kind: KindBool, Get: func(c Config) any { return c.Checkin.Reminder }},
	{Path: "checkin.reminder_at", Kind: KindString, Get: func(c Config) any { return c.Checkin.ReminderAt }},
	{Path: "checkin.start_reminder", Kind: KindBool, Get: func(c Config) any { return c.Checkin.StartReminder }},
	{Path: "checkin.bad_day_min", Kind: KindInt, Get: func(c Config) any { return c.Checkin.BadDayMin }},
	{Path: "checkin.actions", Kind: KindBool, Get: func(c Config) any { return c.Checkin.Actions }},
	{Path: "ingest.poll_interval", Kind: KindDuration, Get: func(c Config) any { return c.Ingest.PollInterval }},
	{Path: "ingest.idle_gap", Kind: KindDuration, Get: func(c Config) any { return c.Ingest.IdleGap }},
	{Path: "ingest.count_subagents_as_parallel", Kind: KindBool, Get: func(c Config) any { return c.Ingest.CountSubagentsAsParallel }},
	{Path: "validation.min_pairs", Kind: KindInt, Get: func(c Config) any { return c.Validation.MinPairs }},
}
