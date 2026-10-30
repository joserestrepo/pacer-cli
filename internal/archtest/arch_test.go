// Package archtest enforces the hexagonal dependency rule: the core
// (internal/core) must not know about adapters or IO libraries. `internal/`
// only blocks imports across modules, so without this test the compiler
// would not catch a violation.
package archtest

import (
	"fmt"
	"os/exec"
	"reflect"
	"sort"
	"strings"
	"testing"
)

const modulePath = "github.com/joserestrepo/pacer-cli"

const (
	corePrefix     = modulePath + "/internal/core"
	adaptersPrefix = modulePath + "/internal/adapters"
	i18nPrefix     = adaptersPrefix + "/i18n"
)

// IO libraries forbidden inside the core.
var ioLibs = []string{
	"modernc.org/sqlite",
	"github.com/godbus/dbus",
	"github.com/fsnotify/fsnotify",
	"github.com/BurntSushi/toml",
	"golang.org/x/term",
	"golang.org/x/text",
	"github.com/nicksnyder/go-i18n",
}

// Pure layers must not depend on orchestration layers.
var pureLayers = []string{
	corePrefix + "/domain",
	corePrefix + "/rhythm",
	corePrefix + "/validation",
}

var orchestrationLayers = []string{
	corePrefix + "/app",
	corePrefix + "/ports",
}

// i18n is presentation text: it may use core/domain (Lang, Anchor) but never
// the orchestration layers.
var i18nForbidden = []string{
	corePrefix + "/ports",
	corePrefix + "/app",
}

// within reports whether pkg is prefix or a subpackage of it. It compares
// with a trailing "/" so "x/term" does not match "x/terminal".
func within(pkg, prefix string) bool {
	return pkg == prefix || strings.HasPrefix(pkg, prefix+"/")
}

func withinAny(pkg string, prefixes []string) bool {
	for _, p := range prefixes {
		if within(pkg, p) {
			return true
		}
	}
	return false
}

// violations is a pure function: it takes a map of package -> dependencies
// (map[string][]string: string keys, slice-of-string values) and returns the
// violations as sorted "package -> dependency" strings. Since it touches no
// disk and runs nothing, it is tested with a table.
func violations(deps map[string][]string) []string {
	var out []string // nil when there are no violations
	for pkg, ds := range deps {
		if within(pkg, i18nPrefix) {
			for _, d := range ds {
				if withinAny(d, i18nForbidden) {
					out = append(out, fmt.Sprintf("%s -> %s", pkg, d))
				}
			}
			continue
		}
		if !within(pkg, corePrefix) {
			continue
		}
		pure := withinAny(pkg, pureLayers)
		for _, d := range ds {
			forbidden := within(d, adaptersPrefix) || withinAny(d, ioLibs) ||
				(pure && withinAny(d, orchestrationLayers))
			if forbidden {
				out = append(out, fmt.Sprintf("%s -> %s", pkg, d))
			}
		}
	}
	// Map iteration order is random in Go, so sort for a deterministic result.
	sort.Strings(out)
	return out
}

// TestCoreDependencies is the real check: it asks the Go toolchain for the
// dependencies of internal/core/... and internal/adapters/i18n/... and applies violations.
func TestCoreDependencies(t *testing.T) {
	// exec.Command builds the command without running it; Output() runs it
	// and returns stdout. -e keeps go list from aborting on broken packages
	// or packages without Go files. {{.Deps}} is the transitive dependency list.
	cmd := exec.Command("go", "list", "-e", "-deps",
		"-f", "{{.ImportPath}} {{join .Deps \" \"}}",
		"./internal/core/...", "./internal/adapters/i18n/...")
	cmd.Dir = "../.." // module root: tests run in the package directory
	raw, err := cmd.Output()
	if err != nil {
		t.Fatalf("go list failed: %v", err)
	}

	deps := map[string][]string{}
	for _, line := range strings.Split(string(raw), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		deps[fields[0]] = fields[1:]
	}

	// internal/core and internal/adapters/i18n may not exist yet or may have no
	// Go packages. That is not a failure, there is just nothing to check.
	if !hasCheckedPackage(deps) {
		t.Skip("internal/core and internal/adapters/i18n have no packages yet")
	}

	if v := violations(deps); len(v) > 0 {
		t.Errorf("architecture violations:\n  %s", strings.Join(v, "\n  "))
	}
}

func hasCheckedPackage(deps map[string][]string) bool {
	for pkg := range deps {
		if within(pkg, corePrefix) || within(pkg, i18nPrefix) {
			return true
		}
	}
	return false
}

// Table-driven test: each case is a struct with a name, input and expected
// output. t.Run creates one subtest per case, so a failure names the exact case.
func TestViolations(t *testing.T) {
	const m = modulePath

	tests := []struct {
		name string
		deps map[string][]string
		want []string
	}{
		{
			name: "no packages, no violations",
			deps: map[string][]string{},
			want: nil,
		},
		{
			name: "clean core",
			deps: map[string][]string{
				m + "/internal/core/domain": {"time", "fmt"},
				m + "/internal/core/app":    {m + "/internal/core/domain", m + "/internal/core/ports"},
			},
			want: nil,
		},
		{
			name: "core depends on adapters",
			deps: map[string][]string{
				m + "/internal/core/app": {m + "/internal/adapters/sqlite"},
			},
			want: []string{m + "/internal/core/app -> " + m + "/internal/adapters/sqlite"},
		},
		{
			name: "core depends on IO libs",
			deps: map[string][]string{
				m + "/internal/core/app": {
					"modernc.org/sqlite",
					"github.com/godbus/dbus/v5",
					"github.com/fsnotify/fsnotify",
					"github.com/BurntSushi/toml",
					"golang.org/x/term",
				},
			},
			want: []string{
				m + "/internal/core/app -> github.com/BurntSushi/toml",
				m + "/internal/core/app -> github.com/fsnotify/fsnotify",
				m + "/internal/core/app -> github.com/godbus/dbus/v5",
				m + "/internal/core/app -> golang.org/x/term",
				m + "/internal/core/app -> modernc.org/sqlite",
			},
		},
		{
			name: "domain depends on app",
			deps: map[string][]string{
				m + "/internal/core/domain": {m + "/internal/core/app"},
			},
			want: []string{m + "/internal/core/domain -> " + m + "/internal/core/app"},
		},
		{
			name: "rhythm depends on ports",
			deps: map[string][]string{
				m + "/internal/core/rhythm": {m + "/internal/core/ports"},
			},
			want: []string{m + "/internal/core/rhythm -> " + m + "/internal/core/ports"},
		},
		{
			name: "validation depends on app",
			deps: map[string][]string{
				m + "/internal/core/validation": {m + "/internal/core/app"},
			},
			want: []string{m + "/internal/core/validation -> " + m + "/internal/core/app"},
		},
		{
			name: "app may depend on ports",
			deps: map[string][]string{
				m + "/internal/core/app": {m + "/internal/core/ports"},
			},
			want: nil,
		},
		{
			name: "core depends on adapters/i18n",
			deps: map[string][]string{
				m + "/internal/core/app": {m + "/internal/adapters/i18n"},
			},
			want: []string{m + "/internal/core/app -> " + m + "/internal/adapters/i18n"},
		},
		{
			name: "core depends on i18n libs",
			deps: map[string][]string{
				m + "/internal/core/app": {
					"golang.org/x/text/message",
					"github.com/nicksnyder/go-i18n/v2/i18n",
				},
			},
			want: []string{
				m + "/internal/core/app -> github.com/nicksnyder/go-i18n/v2/i18n",
				m + "/internal/core/app -> golang.org/x/text/message",
			},
		},
		{
			name: "i18n depends on ports",
			deps: map[string][]string{
				m + "/internal/adapters/i18n": {m + "/internal/core/ports"},
			},
			want: []string{m + "/internal/adapters/i18n -> " + m + "/internal/core/ports"},
		},
		{
			name: "i18n depends on app",
			deps: map[string][]string{
				m + "/internal/adapters/i18n": {m + "/internal/core/app"},
			},
			want: []string{m + "/internal/adapters/i18n -> " + m + "/internal/core/app"},
		},
		{
			name: "i18n may depend on domain",
			deps: map[string][]string{
				m + "/internal/adapters/i18n": {m + "/internal/core/domain", "fmt"},
			},
			want: nil,
		},
		{
			name: "rules do not apply outside core",
			deps: map[string][]string{
				m + "/internal/adapters/sqlite": {"modernc.org/sqlite"},
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := violations(tt.deps)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("violations() = %v, want %v", got, tt.want)
			}
		})
	}
}
