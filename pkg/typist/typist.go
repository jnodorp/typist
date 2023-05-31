package typist

import (
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// Typist to emulate typing.
type Typist struct {
	// WPM (Words Per Minute) to type. A word is defined as 5 keystrokes.
	WPM int

	// Accuracy is the share of correct keystrokes. Must be between 0 (error on every keystroke) and 1 (no errors).
	Accuracy float64
}

// New typist with the given characteristics (WPM and accuracy).
func New(wpm int, a float64) (*Typist, error) {
	if wpm < 1 {
		return nil, fmt.Errorf("WPM (Words per Minute) must be at least one (was %d)", wpm)
	}

	if a < 0 || a > 1 {
		return nil, fmt.Errorf("accuracy must be between 0 and 1 (was %s)", strconv.FormatFloat(a, 'f', -1, 64))
	}

	return &Typist{
		WPM:      wpm,
		Accuracy: a,
	}, nil
}

// Type a string into an io.Writer.
func (t Typist) Type(w io.Writer, s string) error {
	if t.Accuracy < 0 || t.Accuracy > 1 {
		return fmt.Errorf("accuracy must be between 0 and 1 (was %s)", strconv.FormatFloat(t.Accuracy, 'f', -1, 64))
	}

	// Append line break, if missing.
	if !strings.HasSuffix(s, "\n") {
		s = fmt.Sprintf("%s\n", s)
	}

	for i, r := range s {
		// Occasionally skip a character and correct the "mistake".
		if t.Accuracy < rand.Float64() {
			// Do not error if there is no next character or when around a disallowed rune.
			if i+1 < len(s) && errorAllowed(r, rune(s[i+1])) {
				t.Keystroke(w, rune(s[i+1]))
				t.Keystroke(w, '\b')
			}
		}

		if err := t.Keystroke(w, r); err != nil {
			return fmt.Errorf("failed to type %q: %w", strings.TrimSpace(s), err)
		}
	}

	return nil
}

// Keystroke simulation. Write a rune to an io.Writer after a random delay.
func (t Typist) Keystroke(w io.Writer, r rune) error {
	if t.WPM < 1 {
		return fmt.Errorf("WPM (Words per Minute) must be at least one (was %d)", t.WPM)
	}

	time.Sleep(randomDelay(t.WPM))
	bytes := make([]byte, utf8.RuneLen(r))
	utf8.EncodeRune(bytes, r)
	_, err := w.Write(bytes)
	if err != nil {
		return fmt.Errorf("failed to write rune %c: %w", r, err)
	}

	return nil
}

// randomDelay between keystrokes.
func randomDelay(wpm int) time.Duration {
	// Compute average delay between two keystrokes from WPM.
	avgDelay := (1 / (float64(wpm) * 5)) * float64(time.Minute)

	// Assume delay between keystrokes to be normally distributed.
	mu := avgDelay
	sigma := avgDelay / 3
	return time.Duration(rand.NormFloat64()*sigma + mu)
}

// errorAllowed returns false if one of the provided runes results in terminals when deleted via \b. For example, a tab
// will produce multiple spaces. A subsequent \b will remove only one of these spaces.
func errorAllowed(r ...rune) bool {
	for _, v := range r {
		switch v {
		case '\n', '\t':
			return false
		}
	}

	return true
}
