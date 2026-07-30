package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// The coverage/ directory tracks, topic by topic, which parts of the Microsoft
// Writing Style Guide this style implements. Each file mirrors one top-level
// section of the guide; each key is a subtopic (or an A-Z word list entry) set
// to true or false, optionally followed by a comment naming the rules that
// implement it.
//
// This test reports the resulting metric and enforces the invariants that keep
// it honest: values must be canonical booleans, and every rule named in a
// comment must still exist. A rule that gets renamed or merged away otherwise
// leaves the manifest silently claiming coverage it no longer has.

// wordListFile holds the A-Z word list, whose ~850 entries would otherwise
// dominate any single headline percentage, so it's reported separately.
const wordListFile = "a-z.yml"

var ruleRef = regexp.MustCompile(`([A-Za-z]+)\.yml`)

type tally struct{ covered, total int }

func (t tally) pct() float64 {
	if t.total == 0 {
		return 0
	}
	return 100 * float64(t.covered) / float64(t.total)
}

func parseManifest(t *testing.T, path string) tally {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}

	var got tally
	for i, line := range strings.Split(string(data), "\n") {
		lineNo := i + 1

		// A comment names the implementing rules; check they still exist.
		if idx := strings.Index(line, "#"); idx >= 0 {
			for _, m := range ruleRef.FindAllStringSubmatch(line[idx:], -1) {
				rule := filepath.Join("Google", m[1]+".yml")
				if _, err := os.Stat(rule); err != nil {
					t.Errorf("%s:%d: names %s.yml, which doesn't exist", path, lineNo, m[1])
				}
			}
			line = line[:idx]
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		key, value, found := strings.Cut(line, ":")
		if !found {
			t.Errorf("%s:%d: not a key/value pair: %q", path, lineNo, line)
			continue
		}

		switch strings.TrimSpace(value) {
		case "true":
			got.covered++
			got.total++
		case "false":
			got.total++
		default:
			t.Errorf("%s:%d: %q must be exactly true or false, got %q",
				path, lineNo, strings.TrimSpace(key), strings.TrimSpace(value))
		}
	}

	return got
}

func TestCoverage(t *testing.T) {
	paths, err := filepath.Glob(filepath.Join("coverage", "*.yml"))
	if err != nil || len(paths) == 0 {
		t.Fatalf("no coverage manifests found: %v", err)
	}
	sort.Strings(paths)

	var guidelines, wordList tally
	var report []string

	width := 0
	for _, p := range paths {
		if n := len(filepath.Base(p)); n > width {
			width = n
		}
	}

	for _, p := range paths {
		got := parseManifest(t, p)
		report = append(report, fmt.Sprintf("  %-*s  %3d/%-3d  %5.1f%%",
			width, filepath.Base(p), got.covered, got.total, got.pct()))

		if filepath.Base(p) == wordListFile {
			wordList = got
		} else {
			guidelines.covered += got.covered
			guidelines.total += got.total
		}
	}

	summary := fmt.Sprintf("\n\n  guidelines:     %d/%d (%.1f%%)",
		guidelines.covered, guidelines.total, guidelines.pct())
	// Only styles that track the A-Z word list report it; skip the line otherwise
	// rather than printing a meaningless 0/0.
	if wordList.total > 0 {
		summary += fmt.Sprintf("\n  A-Z word list:  %d/%d (%.1f%%)",
			wordList.covered, wordList.total, wordList.pct())
	}

	t.Log("style guide coverage\n" + strings.Join(report, "\n") + summary)
}
