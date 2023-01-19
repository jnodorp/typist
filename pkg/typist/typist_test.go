package typist_test

import (
	"bytes"
	"errors"
	"io"
	"math"
	"testing"

	. "github.com/jnodorp/typist/pkg/typist"

	"github.com/stretchr/testify/assert"
)

func TestNewWPMZero(t *testing.T) {
	typist, err := New(0, 1)
	assert.EqualError(t, err, "WPM (Words per Minute) must be at least one (was 0)")
	assert.Nil(t, typist)
}

func TestNewNegativeAccuracy(t *testing.T) {
	typist, err := New(math.MaxInt, -0.1)
	assert.EqualError(t, err, "accuracy must be between 0 and 1 (was -0.1)")
	assert.Nil(t, typist)
}

func TestNewAccuracyGreaterOne(t *testing.T) {
	typist, err := New(math.MaxInt, 1.11)
	assert.EqualError(t, err, "accuracy must be between 0 and 1 (was 1.11)")
	assert.Nil(t, typist)
}

func TestNew(t *testing.T) {
	typist, err := New(math.MaxInt, 0.5)
	assert.NoError(t, err)
	assert.Equal(t, 0.5, typist.Accuracy)
	assert.Equal(t, math.MaxInt, typist.WPM)
}

func TestKeystrokeWPMZero(t *testing.T) {
	w := new(bytes.Buffer)
	err := Typist{}.Keystroke(w, 'a')
	assert.EqualError(t, err, "WPM (Words per Minute) must be at least one (was 0)")
	assert.Equal(t, "", w.String())
}

func TestKeystrokeWriterError(t *testing.T) {
	w := new(errorWriter)
	err := Typist{
		WPM: 50,
	}.Keystroke(w, 'a')
	assert.EqualError(t, err, "failed to write rune a: just as expected")
}

func TestKeystroke(t *testing.T) {
	w := new(bytes.Buffer)
	err := Typist{
		WPM: 50,
	}.Keystroke(w, 'a')
	assert.NoError(t, err)
	assert.Equal(t, "a", w.String())
}

func TestTypeNegativeErrorProbability(t *testing.T) {
	w := new(bytes.Buffer)
	err := Typist{
		WPM:      math.MaxUint32,
		Accuracy: -0.1,
	}.Type(w, "a")
	assert.EqualError(t, err, "accuracy must be between 0 and 1 (was -0.1)")
	assert.Equal(t, "", w.String())
}

func TestTypeErrorProbabilityGreaterOne(t *testing.T) {
	w := new(bytes.Buffer)
	err := Typist{
		WPM:      math.MaxUint32,
		Accuracy: 1.11,
	}.Type(w, "a")
	assert.EqualError(t, err, "accuracy must be between 0 and 1 (was 1.11)")
	assert.Equal(t, "", w.String())
}

func TestTypeError(t *testing.T) {
	w := new(bytes.Buffer)
	err := Typist{
		WPM:      math.MaxUint32,
		Accuracy: 0,
	}.Type(w, "ab")
	assert.NoError(t, err)
	assert.Equal(t, "b\bab\n", w.String())
}

func TestTypeWriterError(t *testing.T) {
	w := new(errorWriter)
	err := Typist{
		WPM:      math.MaxUint32,
		Accuracy: 1,
	}.Type(w, "ab")
	assert.EqualError(t, err, "failed to type \"ab\": failed to write rune a: just as expected")
}

var _ io.Writer = (*errorWriter)(nil)

type errorWriter struct{}

func (errorWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("just as expected")
}
