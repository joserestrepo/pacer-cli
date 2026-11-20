package config

import (
	"errors"
	"strings"
	"testing"
)

// The defaults must always be a valid configuration.
func TestValidateDefaults(t *testing.T) {
	if err := Validate(Defaults()); err != nil {
		t.Errorf("Validate(Defaults()) = %v, want nil", err)
	}
}

// Spec: config-management, "Requirement: Language", scenario "Invalid language".
func TestValidateLanguage(t *testing.T) {
	tests := []struct {
		name     string
		language string
		wantErr  bool
	}{
		{"auto", "auto", false},
		{"english", "en", false},
		{"spanish", "es", false},
		{"unsupported language", "fr", true},
		{"uppercase is rejected", "ES", true},
		{"empty is rejected", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Defaults()
			c.General.Language = tt.language

			err := Validate(c)

			if !tt.wantErr {
				if err != nil {
					t.Errorf("Validate() = %v, want nil", err)
				}
				return
			}

			var fieldErr FieldError
			if !errors.As(err, &fieldErr) {
				t.Fatalf("Validate() = %v, want a FieldError", err)
			}
			if fieldErr.Key != "general.language" {
				t.Errorf("Key = %q, want %q", fieldErr.Key, "general.language")
			}
			if fieldErr.Rule != RuleOneOf {
				t.Errorf("Rule = %v, want RuleOneOf", fieldErr.Rule)
			}
		})
	}
}

// Spec: config-management, day_starts_at must be "HH:MM" between 00:00 and 12:00.
func TestValidateDayStartsAt(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		wantErr  bool
		wantRule Rule
	}{
		{"default", "04:00", false, 0},
		{"midnight is the lower bound", "00:00", false, 0},
		{"noon is the upper bound", "12:00", false, 0},
		{"after noon is out of range", "12:01", true, RuleRange},
		{"afternoon is out of range", "18:00", true, RuleRange},
		{"not HH:MM", "4am", true, RuleTimeFormat},
		{"hour 25 does not exist", "25:00", true, RuleTimeFormat},
		{"empty is rejected", "", true, RuleTimeFormat},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Defaults()
			c.General.DayStartsAt = tt.value

			err := Validate(c)

			if !tt.wantErr {
				if err != nil {
					t.Errorf("Validate() = %v, want nil", err)
				}
				return
			}

			var fieldErr FieldError
			if !errors.As(err, &fieldErr) {
				t.Fatalf("Validate() = %v, want a FieldError", err)
			}
			if fieldErr.Key != "general.day_starts_at" {
				t.Errorf("Key = %q, want %q", fieldErr.Key, "general.day_starts_at")
			}
			if fieldErr.Rule != tt.wantRule {
				t.Errorf("Rule = %v, want %v", fieldErr.Rule, tt.wantRule)
			}
		})
	}
}

// The message is for logs and tests, so it is in English, names the key and
// explains which rule was broken: each rule has its own message.
func TestFieldErrorMessage(t *testing.T) {
	rules := []Rule{RuleOneOf, RuleTimeFormat, RuleRange}
	seen := map[string]Rule{}

	for _, rule := range rules {
		msg := FieldError{Key: "general.language", Rule: rule}.Error()

		if !strings.Contains(msg, "general.language") {
			t.Errorf("rule %v: Error() = %q, want it to contain the key", rule, msg)
		}
		if other, dup := seen[msg]; dup {
			t.Errorf("rules %v and %v share the message %q", other, rule, msg)
		}
		seen[msg] = rule
	}
}
