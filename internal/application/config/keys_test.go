package config

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

// Every leaf field of Config must have exactly one Key. A leaf is a field
// that is not itself a section, so nested sections such as Scoring.Weights
// are counted field by field. If someone adds a field to Config and forgets
// its Key, this test fails.
func TestKeysCoverEveryConfigField(t *testing.T) {
	want := countLeaves(reflect.TypeOf(Config{}))

	if len(Keys) != want {
		t.Errorf("len(Keys) = %d, want %d (one per Config leaf field)", len(Keys), want)
	}
}

// countLeaves counts the non-struct fields of t, descending into struct
// fields. It calls itself for each nested section: that is recursion.
func countLeaves(t reflect.Type) int {
	n := 0
	for _, f := range reflect.VisibleFields(t) {
		if f.Type.Kind() == reflect.Struct {
			n += countLeaves(f.Type)
		} else {
			n++
		}
	}
	return n
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
		KindFloat:    reflect.TypeOf(0.0),
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

// Every Path must name its field: "checkin.bad_day_min" must read
// Config.Checkin.BadDayMin, and "scoring.weights.hours" must read
// Config.Scoring.Weights.Hours. This catches typos such as "checkin.bad_day".
func TestKeyPathNamesItsField(t *testing.T) {
	c := reflect.ValueOf(Defaults())
	for _, k := range Keys {
		t.Run(k.Path, func(t *testing.T) {
			parts := strings.Split(k.Path, ".")
			if len(parts) < 2 {
				t.Fatalf("path %q is not \"section.name\"", k.Path)
			}
			field := c
			for _, p := range parts {
				field = field.FieldByName(snakeToCamel(p))
				if !field.IsValid() {
					t.Fatalf("no Config field for %q (stopped at %q)", k.Path, p)
				}
			}
			if got := k.Get(Defaults()); got != field.Interface() {
				t.Errorf("Get returns %v, but the field is %v", got, field.Interface())
			}
		})
	}
}

// snakeToCamel turns "bad_day_min" into "BadDayMin".
func snakeToCamel(s string) string {
	parts := strings.Split(s, "_")
	for i, p := range parts {
		if p != "" {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, "")
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
