package main

import (
	"bytes"
	"strings"
	"testing"
)

// TestRunOneShot exercises the one-shot mode (a positional argument present)
// through the io.Reader/io.Writer seams. The interactive mode (no positional
// argument) is covered by interactive_test.go.
//
// Note: the data-load failure exit code (2) is part of the contract but is not
// reachable in a unit test because the dataset is embedded at build time via
// //go:embed; city.Load() therefore cannot fail here.
func TestRunOneShot(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantCode   int
		wantStdout []string // substrings expected in stdout
		wantStderr []string // substrings expected in stderr
		notStdout  []string // substrings that must NOT appear in stdout
	}{
		{
			name:       "matching query lists un-numbered cities",
			args:       []string{"London"},
			wantCode:   0,
			wantStdout: []string{"London", "United Kingdom"},
			notStdout:  []string{"1. "}, // one-shot output is un-numbered (FR-013)
		},
		{
			name:       "no match prints message and succeeds",
			args:       []string{"zzzzzznotacity"},
			wantCode:   0,
			wantStdout: []string{`No cities found matching "zzzzzznotacity".`},
		},
		{
			name:       "empty argument is an invalid-usage error",
			args:       []string{""},
			wantCode:   1,
			wantStderr: []string{"Please provide a city name"},
		},
		{
			name:       "whitespace-only argument is an invalid-usage error",
			args:       []string{"   "},
			wantCode:   1,
			wantStderr: []string{"Please provide a city name"},
		},
		{
			name:     "help flag succeeds",
			args:     []string{"-h"},
			wantCode: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := run(tt.args, strings.NewReader(""), &stdout, &stderr)
			if code != tt.wantCode {
				t.Fatalf("run(%q) exit = %d, want %d (stdout=%q stderr=%q)",
					tt.args, code, tt.wantCode, stdout.String(), stderr.String())
			}
			for _, want := range tt.wantStdout {
				if !strings.Contains(stdout.String(), want) {
					t.Errorf("stdout = %q, want substring %q", stdout.String(), want)
				}
			}
			for _, want := range tt.wantStderr {
				if !strings.Contains(stderr.String(), want) {
					t.Errorf("stderr = %q, want substring %q", stderr.String(), want)
				}
			}
			for _, notWant := range tt.notStdout {
				if strings.Contains(stdout.String(), notWant) {
					t.Errorf("stdout = %q, must not contain %q", stdout.String(), notWant)
				}
			}
		})
	}
}
