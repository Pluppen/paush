// ABOUTME: Tests for the breathing animation frame generator.
// ABOUTME: Verifies frame dimensions, growth with breath phase, seeded randomness, and labels.
package main

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func countInk(frame string) int {
	n := 0
	for _, r := range frame {
		if r != ' ' && r != '\n' {
			n++
		}
	}
	return n
}

func TestFrameDimensions(t *testing.T) {
	a := NewAnimation(1, 44, 15)
	frame := a.Frame(1.0)
	lines := strings.Split(frame, "\n")
	if len(lines) != 15 {
		t.Fatalf("expected 15 lines, got %d", len(lines))
	}
	for i, line := range lines {
		if w := utf8.RuneCountInString(line); w != 44 {
			t.Errorf("line %d: expected width 44, got %d", i, w)
		}
	}
}

func newMotionAnimation(m motionFunc) *Animation {
	return &Animation{
		width:   44,
		height:  15,
		shape:   shapes[0],
		charset: charsets[0],
		motion:  m,
		noise:   make([]float64, 44*15),
	}
}

func TestEachMotionVariesWithPhase(t *testing.T) {
	for i, m := range motions {
		a := newMotionAnimation(m)
		if a.Frame(0.2) == a.Frame(0.8) {
			t.Errorf("motion %d renders identical frames at phases 0.2 and 0.8", i)
		}
	}
}

func TestEachMotionDrawsMidBreath(t *testing.T) {
	for i, m := range motions {
		a := newMotionAnimation(m)
		for _, phase := range []float64{0.5, 1.0} {
			if countInk(a.Frame(phase)) == 0 {
				t.Errorf("motion %d renders an empty frame at phase %v", i, phase)
			}
		}
	}
}

func TestMotionsAreDistinct(t *testing.T) {
	frames := make([]string, len(motions))
	for i, m := range motions {
		frames[i] = newMotionAnimation(m).Frame(0.5)
	}
	for i := 0; i < len(frames); i++ {
		for j := i + 1; j < len(frames); j++ {
			if frames[i] == frames[j] {
				t.Errorf("motions %d and %d render identical frames at phase 0.5", i, j)
			}
		}
	}
}

func TestSameSeedSameFrames(t *testing.T) {
	a := NewAnimation(3, 44, 15)
	b := NewAnimation(3, 44, 15)
	for _, phase := range []float64{0.0, 0.5, 1.0} {
		if a.Frame(phase) != b.Frame(phase) {
			t.Errorf("same seed produced different frames at phase %v", phase)
		}
	}
}

func TestDifferentSeedsDiffer(t *testing.T) {
	first := NewAnimation(0, 44, 15).Frame(1.0)
	for seed := int64(1); seed <= 9; seed++ {
		if NewAnimation(seed, 44, 15).Frame(1.0) != first {
			return
		}
	}
	t.Error("seeds 0-9 all produced identical frames")
}

func TestLabel(t *testing.T) {
	if got := Label(true); !strings.Contains(got, "breathe in") {
		t.Errorf("inhale label = %q, want it to contain \"breathe in\"", got)
	}
	if got := Label(false); !strings.Contains(got, "breathe out") {
		t.Errorf("exhale label = %q, want it to contain \"breathe out\"", got)
	}
}
