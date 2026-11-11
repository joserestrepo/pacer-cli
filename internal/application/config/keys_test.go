package config

import (
	"reflect"
	"testing"
	"time"
)

// Every field of every section in Config must have exactly one Key.
// If someone adds a field to Config and forgets its Key, this test fails.
func TestKeysCoverEveryConfigField(t *testing.T) {
	want := 0
	for _, section := range reflect.VisibleFields(reflect.TypeOf(Config{})) {
		want += section.Type.NumField()
	}

	if len(Keys) != want {
		t.Errorf("len(Keys) = %d, want %d (one per Config field)", len(Keys), want)
	}
}

func TestKeyPathsAreUnique(t *testing.T) {
	seen := map[string]bool{}
	for _, k := range Keys {
		if seen[k.Path] {
			t.Errorf("duplicate key %q", k.Path)
		}
		seen[k.Path] = true
	}
}

// The Kind of a key must match the Go type its Get returns, so the TOML
// adapter can render and parse the value correctly.
func TestKeyKindMatchesValueType(t *testing.T) {
	typeOf := map[Kind]reflect.Type{
		KindBool:     reflect.TypeOf(false),
		KindInt:      reflect.TypeOf(0),
		KindString:   reflect.TypeOf(""),
		KindDuration: reflect.TypeOf(time.Duration(0)),
	}

	c := Defaults()
	for _, k := range Keys {
		t.Run(k.Path, func(t *testing.T) {
			got := reflect.TypeOf(k.Get(c))
			if got != typeOf[k.Kind] {
				t.Errorf("Get returns %v, but Kind expects %v", got, typeOf[k.Kind])
			}
		})
	}
}

// Spot-check that paths point to the right field: Get on Defaults() must
// return that field's default.
func TestKeyGetReturnsDefault(t *testing.T) {
	tests := []struct {
		path string
		want any
	}{
		{"general.language", "auto"},
		{"budget.active_hours", 6},
		{"budget.cutoff", "19:00"},
		{"checkin.reminder", true},
		{"alerts.cooldown", 60 * time.Minute},
	}

	byPath := map[string]Key{}
	for _, k := range Keys {
		byPath[k.Path] = k
	}

	c := Defaults()
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			k, ok := byPath[tt.path]
			if !ok {
				t.Fatalf("key %q not found in Keys", tt.path)
			}
			if got := k.Get(c); got != tt.want {
				t.Errorf("Get(Defaults()) = %v, want %v", got, tt.want)
			}
		})
	}
}
