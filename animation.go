// ABOUTME: Generates randomized ASCII breathing animations as text frames.
// ABOUTME: Each seed picks a shape, character gradient, and color; Frame(t) renders one breath phase.
package main

import (
	"math"
	"math/rand"
	"strings"
)

// shapeFunc returns the normalized distance of a point from the shape's
// center; a cell is drawn when the distance is within the current radius.
type shapeFunc func(dx, dy float64) float64

var shapes = []shapeFunc{
	// circle
	func(dx, dy float64) float64 {
		return math.Hypot(dx, dy)
	},
	// diamond
	func(dx, dy float64) float64 {
		return math.Abs(dx) + math.Abs(dy)
	},
	// square
	func(dx, dy float64) float64 {
		return math.Max(math.Abs(dx), math.Abs(dy))
	},
	// flower
	func(dx, dy float64) float64 {
		r := math.Hypot(dx, dy)
		theta := math.Atan2(dy, dx)
		return r / (0.75 + 0.25*math.Abs(math.Sin(3*theta)))
	},
	// starburst
	func(dx, dy float64) float64 {
		r := math.Hypot(dx, dy)
		theta := math.Atan2(dy, dx)
		return r / (0.6 + 0.4*math.Pow(math.Abs(math.Cos(4*theta)), 0.5))
	},
}

// motionFunc maps a cell to its position on the charset gradient for breath
// phase t in [0,1]; ok is false when the cell stays blank. dist is the cell's
// shape distance, dy its vertical position in [-1,1], and theta its angle.
type motionFunc func(dist, dy, theta, t, noise float64) (pos float64, ok bool)

var motions = []motionFunc{
	// expand: the shape grows on inhale and shrinks on exhale
	func(dist, dy, theta, t, noise float64) (float64, bool) {
		r := 0.08 + 0.92*t
		d := dist + noise
		return d / r, d <= r
	},
	// ripple: rings travel outward on inhale and inward on exhale
	func(dist, dy, theta, t, noise float64) (float64, bool) {
		r := 0.3 + 0.7*t
		d := dist + noise
		return d / r, d <= r && math.Sin(2*math.Pi*(2.5*d-1.5*t)) > 0
	},
	// fill: the shape fills from the bottom on inhale and drains on exhale
	func(dist, dy, theta, t, noise float64) (float64, bool) {
		d := dist + noise
		return d, d <= 1 && dy >= 1-2*t
	},
	// spiral: arms sweep around the center as the shape breathes
	func(dist, dy, theta, t, noise float64) (float64, bool) {
		r := 0.2 + 0.8*t
		d := dist + noise
		return d / r, d <= r && math.Sin(3*theta+8*d-4*t) > 0.2
	},
	// ring: a pulse band travels from the center to the edge and back
	func(dist, dy, theta, t, noise float64) (float64, bool) {
		const halfWidth = 0.18
		d := math.Abs(dist + noise - (0.15 + 0.75*t))
		return d / halfWidth, d <= halfWidth
	},
}

// charsets are ordered dense (center) to faint (edge).
var charsets = [][]rune{
	[]rune("@%#*+=-:."),
	[]rune("█▓▒░·"),
	[]rune("◉◎○∘·"),
	[]rune("✦✧*+·"),
	[]rune("&$o=~-."),
}

// colors are ANSI 256-color foreground codes in calm tones.
var colors = []string{
	"\x1b[38;5;117m", // sky blue
	"\x1b[38;5;121m", // mint
	"\x1b[38;5;183m", // lavender
	"\x1b[38;5;216m", // peach
	"\x1b[38;5;159m", // pale cyan
}

type Animation struct {
	width, height int
	shape         shapeFunc
	motion        motionFunc
	charset       []rune
	Color         string
	noise         []float64
}

func NewAnimation(seed int64, width, height int) *Animation {
	rng := rand.New(rand.NewSource(seed))
	a := &Animation{
		width:   width,
		height:  height,
		shape:   shapes[rng.Intn(len(shapes))],
		motion:  motions[rng.Intn(len(motions))],
		charset: charsets[rng.Intn(len(charsets))],
		Color:   colors[rng.Intn(len(colors))],
		noise:   make([]float64, width*height),
	}
	for i := range a.noise {
		a.noise[i] = rng.Float64() * 0.12
	}
	return a
}

// Frame renders the animation at breath phase t in [0,1], where 0 is fully
// exhaled and 1 is fully inhaled. Lines are padded to the full width.
func (a *Animation) Frame(t float64) string {
	var b strings.Builder
	cx := float64(a.width-1) / 2
	cy := float64(a.height-1) / 2
	for y := 0; y < a.height; y++ {
		for x := 0; x < a.width; x++ {
			// Coordinates are normalized to [-1,1] on each axis, so the
			// shape stretches to fill the frame regardless of cell aspect.
			dx := (float64(x) - cx) / cx
			dy := (float64(y) - cy) / cy
			dist := a.shape(dx, dy)
			pos, ok := a.motion(dist, dy, math.Atan2(dy, dx), t, a.noise[y*a.width+x])
			if ok {
				idx := int(pos * float64(len(a.charset)))
				if idx >= len(a.charset) {
					idx = len(a.charset) - 1
				}
				if idx < 0 {
					idx = 0
				}
				b.WriteRune(a.charset[idx])
			} else {
				b.WriteRune(' ')
			}
		}
		if y < a.height-1 {
			b.WriteByte('\n')
		}
	}
	return b.String()
}

// Label returns the breathing guidance text for the current phase.
func Label(inhale bool) string {
	if inhale {
		return "breathe in"
	}
	return "breathe out"
}
