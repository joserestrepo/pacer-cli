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

// The message is for logs and tests, so it is in English and names the key.
func TestFieldErrorMessageNamesTheKey(t *testing.T) {
	err := FieldError{Key: "general.language", Rule: RuleOneOf}

	if msg := err.Error(); !strings.Contains(msg, "general.language") {
		t.Errorf("Error() = %q, want it to contain the key", msg)
	}
}
