// ABOUTME: Tests for terminal-size handling in the CLI entry point.
// ABOUTME: Verifies frame dimensions derived from terminal size, including fallbacks.
package main

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestBreatheRestoresTerminal(t *testing.T) {
	var buf bytes.Buffer
	breathe(&buf, 20*time.Millisecond, 30, 10)
	out := buf.String()
	if !strings.HasPrefix(out, enterAltScreen) {
		t.Error("output does not start by entering the alternate screen")
	}
	if !strings.HasSuffix(out, leaveAltScreen+showCursor) {
		t.Error("output does not end by leaving the alternate screen and restoring the cursor")
	}
	for _, want := range []string{"breathe in", "breathe out", "carry on, mindfully"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q", want)
		}
	}
}

func TestResolveDuration(t *testing.T) {
	cases := []struct {
		name     string
		flagVal  string
		envVal   string
		want     time.Duration
		wantErr  bool
	}{
		{"default", "", "", 2 * time.Second, false},
		{"flag only", "8", "", 8 * time.Second, false},
		{"env only", "", "4", 4 * time.Second, false},
		{"flag beats env", "8", "4", 8 * time.Second, false},
		{"fractional seconds", "2.5", "", 2500 * time.Millisecond, false},
		{"flag not a number", "soon", "", 0, true},
		{"env not a number", "", "later", 0, true},
		{"zero", "0", "", 0, true},
		{"negative", "-3", "", 0, true},
		{"over ten minutes", "601", "", 0, true},
	}
	for _, c := range cases {
		got, err := resolveDuration(c.flagVal, c.envVal)
		if c.wantErr {
			if err == nil {
				t.Errorf("%s: expected error, got %v", c.name, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("%s: unexpected error: %v", c.name, err)
		} else if got != c.want {
			t.Errorf("%s: resolveDuration(%q, %q) = %v, want %v", c.name, c.flagVal, c.envVal, got, c.want)
		}
	}
}

func TestPacingFillsDuration(t *testing.T) {
	for _, d := range []time.Duration{time.Second, 2 * time.Second, 8 * time.Second, 30 * time.Second} {
		steps, delay := pacing(d)
		if steps < 8 {
			t.Errorf("pacing(%v): %d steps is too choppy", d, steps)
		}
		total := 2 * time.Duration(steps+1) * delay
		if diff := (total - d).Abs(); diff > 5*time.Millisecond {
			t.Errorf("pacing(%v): frames span %v, off by %v", d, total, diff)
		}
	}
}

func TestFrameSize(t *testing.T) {
	cases := []struct {
		name           string
		termW, termH   int
		wantW, wantH   int
	}{
		{"standard terminal", 80, 24, 80, 21},
		{"large terminal", 200, 50, 200, 47},
		{"unknown size", 0, 0, 44, 15},
		{"too narrow", 10, 50, 44, 15},
		{"too short", 80, 6, 44, 15},
	}
	for _, c := range cases {
		w, h := frameSize(c.termW, c.termH)
		if w != c.wantW || h != c.wantH {
			t.Errorf("%s: frameSize(%d, %d) = (%d, %d), want (%d, %d)",
				c.name, c.termW, c.termH, w, h, c.wantW, c.wantH)
		}
	}
}
