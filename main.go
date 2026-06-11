// ABOUTME: Entry point for paush, a short mindfulness break in the terminal.
// ABOUTME: Plays a randomized ASCII breathing animation; duration and typo aliases are configurable.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"golang.org/x/term"
)

const (
	defaultWidth    = 44
	defaultHeight   = 15
	defaultDuration = 2 * time.Second
	maxDuration     = 600 * time.Second
	targetFrameRate = 30
	outroHold       = 400 * time.Millisecond
	reset           = "\x1b[0m"
	hideCursor      = "\x1b[?25l"
	showCursor      = "\x1b[?25h"
	// The animation plays on the alternate screen buffer so that leaving
	// it restores the terminal exactly as it was before the run.
	enterAltScreen = "\x1b[?1049h\x1b[H"
	leaveAltScreen = "\x1b[?1049l"
)

// resolveDuration picks the animation length from the -t flag value or the
// PAUSH_DURATION environment variable, both in seconds, flag first.
func resolveDuration(flagVal, envVal string) (time.Duration, error) {
	src, name := flagVal, "-t"
	if src == "" {
		src, name = envVal, "PAUSH_DURATION"
	}
	if src == "" {
		return defaultDuration, nil
	}
	secs, err := strconv.ParseFloat(src, 64)
	if err != nil || secs <= 0 || time.Duration(secs*float64(time.Second)) > maxDuration {
		return 0, fmt.Errorf("%s must be a number of seconds between 0 and 600, got %q", name, src)
	}
	return time.Duration(secs * float64(time.Second)), nil
}

// pacing splits half a breath into frames at roughly the target frame rate.
// Each half renders steps+1 frames (phases 0 through 1 inclusive), spaced by
// delay, so one full breath spans the whole duration.
func pacing(d time.Duration) (steps int, delay time.Duration) {
	half := d / 2
	steps = int(half * targetFrameRate / time.Second)
	if steps < 8 {
		steps = 8
	}
	return steps, half / time.Duration(steps+1)
}

// frameSize converts a terminal size into animation dimensions, reserving
// rows for the breathing label. Falls back to the defaults when the size is
// unknown or too small to animate.
func frameSize(termW, termH int) (int, int) {
	w, h := termW, termH-3
	if w < 20 || h < 8 {
		return defaultWidth, defaultHeight
	}
	return w, h
}

func center(s string, width int) string {
	pad := (width - utf8.RuneCountInString(s)) / 2
	if pad < 0 {
		pad = 0
	}
	return strings.Repeat(" ", pad) + s
}

func render(w io.Writer, anim *Animation, t float64, inhale bool) {
	fmt.Fprint(w, anim.Color+anim.Frame(t)+reset+"\n\n"+center(Label(inhale), anim.width))
	// Move the cursor back to the top of the frame for the next draw.
	fmt.Fprintf(w, "\r\x1b[%dA", anim.height+1)
}

func breathe(w io.Writer, d time.Duration, width, height int) {
	cleanup := func() { fmt.Fprint(w, leaveAltScreen+showCursor) }

	interrupted := make(chan os.Signal, 1)
	signal.Notify(interrupted, os.Interrupt)
	go func() {
		<-interrupted
		cleanup()
		os.Exit(1)
	}()

	anim := NewAnimation(time.Now().UnixNano(), width, height)
	steps, delay := pacing(d)
	fmt.Fprint(w, enterAltScreen+hideCursor)
	start := time.Now()
	frame := 0
	for _, inhale := range []bool{true, false} {
		for i := 0; i <= steps; i++ {
			t := float64(i) / float64(steps)
			if !inhale {
				t = 1 - t
			}
			render(w, anim, t, inhale)
			frame++
			// Sleeping toward absolute deadlines keeps the total
			// duration accurate regardless of render time.
			time.Sleep(time.Until(start.Add(time.Duration(frame) * delay)))
		}
	}
	// Hold the outro on the label line briefly before restoring the screen.
	fmt.Fprintf(w, "\x1b[%dB\r\x1b[2K%s", height+1, center("carry on, mindfully", width))
	time.Sleep(outroHold)
	cleanup()
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "paush:", err)
	os.Exit(1)
}

func takeBreak(d time.Duration) {
	termW, termH, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		termW, termH = 0, 0
	}
	width, height := frameSize(termW, termH)
	breathe(os.Stdout, d, width, height)
}

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "install":
			if err := install(); err != nil {
				fatal(err)
			}
			return
		case "uninstall":
			if err := uninstall(); err != nil {
				fatal(err)
			}
			return
		case "typo":
			// Invoked through a typo alias; whatever arguments the
			// mistyped command carried are ignored.
			d, err := resolveDuration("", os.Getenv("PAUSH_DURATION"))
			if err != nil {
				fatal(err)
			}
			takeBreak(d)
			return
		}
	}

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), `usage: paush [-t seconds] | install | uninstall

  paush            play a breathing break (duration from -t or PAUSH_DURATION)
  paush install    add typo aliases (sl, gti, cta, ...) to your shell rc
  paush uninstall  remove the typo aliases

`)
		flag.PrintDefaults()
	}
	durationFlag := flag.String("t", "", "animation duration in seconds (default 2, or PAUSH_DURATION)")
	flag.Parse()
	d, err := resolveDuration(*durationFlag, os.Getenv("PAUSH_DURATION"))
	if err != nil {
		fatal(err)
	}
	takeBreak(d)
}
