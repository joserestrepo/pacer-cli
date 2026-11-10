// Package archtest enforces the hexagonal dependency rule: the domain
// (internal/domain) and the application (internal/application) must not know
// about infrastructure or IO libraries. `internal/` only blocks imports
// across modules, so without this test the compiler would not catch a
// violation.
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
	domainPrefix         = modulePath + "/internal/domain"
	applicationPrefix    = modulePath + "/internal/application"
	infrastructurePrefix = modulePath + "/internal/infrastructure"
	i18nPrefix           = infrastructurePrefix + "/i18n"
)

// IO libraries forbidden inside the domain and the application.
var ioLibs = []string{
	"modernc.org/sqlite",
	"github.com/godbus/dbus",
	"github.com/fsnotify/fsnotify",
	"github.com/BurntSushi/toml",
	"golang.org/x/term",
	"golang.org/x/text",
	"github.com/nicksnyder/go-i18n",
}

// The domain is the innermost layer: it must not depend on the application.
var domainForbidden = []string{
	applicationPrefix,
	infrastructurePrefix,
}

// i18n is presentation text: it may use the domain (Lang, Anchor) but never
// the application's ports or use cases.
var i18nForbidden = []string{
	applicationPrefix + "/ports",
	applicationPrefix + "/usecase",
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
		for _, d := range ds {
			if forbidden(pkg, d) {
				out = append(out, fmt.Sprintf("%s -> %s", pkg, d))
			}
		}
	}
	// Map iteration order is random in Go, so sort for a deterministic result.
	sort.Strings(out)
	return out
}

// forbidden reports whether pkg is not allowed to depend on dep. Packages
// outside domain, application and infrastructure/i18n are not checked.
func forbidden(pkg, dep string) bool {
	switch {
	case within(pkg, domainPrefix):
		return withinAny(dep, domainForbidden) || withinAny(dep, ioLibs)
	case within(pkg, applicationPrefix):
		return within(dep, infrastructurePrefix) || withinAny(dep, ioLibs)
	case within(pkg, i18nPrefix):
		return withinAny(dep, i18nForbidden)
	}
	return false
}

// TestDependencies is the real check: it asks the Go toolchain for the
// dependencies of the checked packages and applies violations.
func TestDependencies(t *testing.T) {
	// exec.Command builds the command without running it; Output() runs it
	// and returns stdout. -e keeps go list from aborting on broken packages
	// or packages without Go files. {{.Deps}} is the transitive dependency list.
	cmd := exec.Command("go", "list", "-e", "-deps",
		"-f", "{{.ImportPath}} {{join .Deps \" \"}}",
		"./internal/domain/...", "./internal/application/...",
		"./internal/infrastructure/i18n/...")
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

	// The checked packages may not exist yet or may have no Go packages. That is not a failure, there is just nothing to check.
	if !hasCheckedPackage(deps) {
		t.Skip("domain, application and infrastructure/i18n have no packages yet")
	}

	if v := violations(deps); len(v) > 0 {
		t.Errorf("architecture violations:\n  %s", strings.Join(v, "\n  "))
	}
}

func hasCheckedPackage(deps map[string][]string) bool {
	for pkg := range deps {
		if within(pkg, domainPrefix) || within(pkg, applicationPrefix) || within(pkg, i18nPrefix) {
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
			name: "clean layers",
			deps: map[string][]string{
				m + "/internal/domain":              {"time", "fmt"},
				m + "/internal/domain/rhythm":       {m + "/internal/domain"},
				m + "/internal/application/config":  {m + "/internal/domain"},
				m + "/internal/application/usecase": {m + "/internal/domain", m + "/internal/application/ports"},
			},
			want: nil,
		},
		{
			name: "application may depend on domain",
			deps: map[string][]string{
				m + "/internal/application/config": {m + "/internal/domain"},
			},
			want: nil,
		},
		{
			name: "usecase may depend on ports",
			deps: map[string][]string{
				m + "/internal/application/usecase": {m + "/internal/application/ports"},
			},
			want: nil,
		},
		{
			name: "rhythm may depend on domain",
			deps: map[string][]string{
				m + "/internal/domain/rhythm": {m + "/internal/domain"},
			},
			want: nil,
		},
		{
			name: "domain depends on application",
			deps: map[string][]string{
				m + "/internal/domain": {m + "/internal/application/usecase"},
			},
			want: []string{m + "/internal/domain -> " + m + "/internal/application/usecase"},
		},
		{
			name: "rhythm depends on ports",
			deps: map[string][]string{
				m + "/internal/domain/rhythm": {m + "/internal/application/ports"},
			},
			want: []string{m + "/internal/domain/rhythm -> " + m + "/internal/application/ports"},
		},
		{
			name: "validation depends on config",
			deps: map[string][]string{
				m + "/internal/domain/validation": {m + "/internal/application/config"},
			},
			want: []string{m + "/internal/domain/validation -> " + m + "/internal/application/config"},
		},
		{
			name: "domain depends on infrastructure",
			deps: map[string][]string{
				m + "/internal/domain": {m + "/internal/infrastructure/sqlite"},
			},
			want: []string{m + "/internal/domain -> " + m + "/internal/infrastructure/sqlite"},
		},
		{
			name: "application depends on infrastructure",
			deps: map[string][]string{
				m + "/internal/application/usecase": {m + "/internal/infrastructure/sqlite"},
			},
			want: []string{m + "/internal/application/usecase -> " + m + "/internal/infrastructure/sqlite"},
		},
		{
			name: "application depends on infrastructure/i18n",
			deps: map[string][]string{
				m + "/internal/application/usecase": {m + "/internal/infrastructure/i18n"},
			},
			want: []string{m + "/internal/application/usecase -> " + m + "/internal/infrastructure/i18n"},
		},
		{
			name: "application depends on IO libs",
			deps: map[string][]string{
				m + "/internal/application/usecase": {
					"modernc.org/sqlite",
					"github.com/godbus/dbus/v5",
					"github.com/fsnotify/fsnotify",
					"github.com/BurntSushi/toml",
					"golang.org/x/term",
				},
			},
			want: []string{
				m + "/internal/application/usecase -> github.com/BurntSushi/toml",
				m + "/internal/application/usecase -> github.com/fsnotify/fsnotify",
				m + "/internal/application/usecase -> github.com/godbus/dbus/v5",
				m + "/internal/application/usecase -> golang.org/x/term",
				m + "/internal/application/usecase -> modernc.org/sqlite",
			},
		},
		{
			name: "domain depends on IO libs",
			deps: map[string][]string{
				m + "/internal/domain": {"github.com/BurntSushi/toml"},
			},
			want: []string{m + "/internal/domain -> github.com/BurntSushi/toml"},
		},
		{
			name: "application depends on i18n libs",
			deps: map[string][]string{
				m + "/internal/application/usecase": {
					"golang.org/x/text/message",
					"github.com/nicksnyder/go-i18n/v2/i18n",
				},
			},
			want: []string{
				m + "/internal/application/usecase -> github.com/nicksnyder/go-i18n/v2/i18n",
				m + "/internal/application/usecase -> golang.org/x/text/message",
			},
		},
		{
			name: "i18n depends on ports",
			deps: map[string][]string{
				m + "/internal/infrastructure/i18n": {m + "/internal/application/ports"},
			},
			want: []string{m + "/internal/infrastructure/i18n -> " + m + "/internal/application/ports"},
		},
		{
			name: "i18n depends on usecase",
			deps: map[string][]string{
				m + "/internal/infrastructure/i18n": {m + "/internal/application/usecase"},
			},
			want: []string{m + "/internal/infrastructure/i18n -> " + m + "/internal/application/usecase"},
		},
		{
			name: "i18n may depend on domain",
			deps: map[string][]string{
				m + "/internal/infrastructure/i18n": {m + "/internal/domain", "fmt"},
			},
			want: nil,
		},
		{
			name: "rules do not apply to other infrastructure",
			deps: map[string][]string{
				m + "/internal/infrastructure/sqlite": {"modernc.org/sqlite", m + "/internal/application/ports"},
			},
			want: nil,
		},
		{
			name: "rules do not apply to cmd",
			deps: map[string][]string{
				m + "/cmd/pacer": {m + "/internal/infrastructure/sqlite", m + "/internal/application/usecase"},
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
