package config

import (
	"fmt"
	"slices"
	"time"
)

// Rule identifies which validation rule a setting broke. The CLI uses it to
// pick the message shown to the user in their language.
type Rule int

// Validation rules.
const (
	RuleOneOf      Rule = iota // the value must be one of a fixed list
	RuleRange                  // the value is well formed but outside its allowed limits
	RuleTimeFormat             // the value must be a local time in "HH:MM" form
)

// FieldError reports that one setting has an invalid value. It carries data
// instead of user-facing text, so the caller can translate it.
type FieldError struct {
	// Key is the path of the invalid setting, e.g. "general.language".
	Key string
	// Rule is the rule the value broke.
	Rule Rule
}

// Error returns an English message for logs and tests. It is not shown to
// the user.
func (e FieldError) Error() string {
	switch e.Rule {
	case RuleOneOf:
		return fmt.Sprintf("invalid value for %q: must be one of the allowed values", e.Key)
	case RuleRange:
		return fmt.Sprintf("invalid value for %q: out of the allowed range", e.Key)
	case RuleTimeFormat:
		return fmt.Sprintf("invalid value for %q: must be a time in HH:MM format", e.Key)
	default:
		return fmt.Sprintf("invalid value for %q", e.Key)
	}
}

// languages lists the accepted values of general.language. The match is
// exact: "ES" is rejected.
var languages = []string{
	"auto",
	"en",
	"es",
}

// Validate checks that every setting in c has an allowed value. It returns
// nil when c is valid, or a FieldError for the first invalid setting.
func Validate(c Config) error {
	isValid := slices.Contains(languages, c.General.Language)
	if !isValid {
		return FieldError{Key: "general.language", Rule: RuleOneOf}
	}
	start, err := time.Parse("15:04", c.General.DayStartsAt)
	if err != nil {
		return FieldError{Key: "general.day_starts_at", Rule: RuleTimeFormat}
	}
	startMinutes := (start.Hour() * 60) + start.Minute()
	latestStart := 12 * 60
	if startMinutes > latestStart {
		return FieldError{Key: "general.day_starts_at", Rule: RuleRange}
	}
	return nil
}
