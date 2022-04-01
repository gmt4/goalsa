// Copyright 2015-2016 Cocoon Labs Ltd.
//
// See LICENSE file for terms and conditions.

package alsa

import (
	"testing"

	"github.com/cocoonlife/testify/assert"
)

func TestCapture(t *testing.T) {
	a := assert.New(t)

	c, err := NewCaptureDevice("nonexistent", 1, FormatS16LE, 44100,
		BufferParams{})

	a.Equal(c, (*CaptureDevice)(nil), "capture device is nil")
	a.Error(err, "no device error")

	c, err = NewCaptureDevice("null", 1, FormatS32LE, 0, BufferParams{})

	a.Equal(c, (*CaptureDevice)(nil), "capture device is nil")
	a.Error(err, "bad rate error")

	c, err = NewCaptureDevice("null", 1, FormatS32LE, 44100, BufferParams{})

	a.NoError(err, "created capture device")

	b0 := 3
	samples, err := c.Read(b0)

	a.Error(err, "wrong type error")
	a.Equal(samples, 0, "no samples read")

	b1 := make([]int8, 100)
	samples, err = c.Read(b1)

	a.Error(err, "wrong type error")
	a.Equal(samples, 0, "no samples read")

	b2 := make([]int16, 200)
	samples, err = c.Read(b2)

	a.Error(err, "wrong type error")
	a.Equal(samples, 0, "no samples read")

	b3 := make([]float64, 50)
	samples, err = c.Read(b3)

	a.Error(err, "wrong type error")
	a.Equal(samples, 0, "no samples read")

	b4 := make([]int32, 200)
	samples, err = c.Read(b4)

	a.NoError(err, "read samples ok")
	a.Equal(len(b2), samples, "correct number of samples read")

	c.Close()

	// extras ints

	c, err = NewCaptureDevice("null", 1, FormatS8, 22050, BufferParams{})
	a.NoError(err, "created capture device")

	samples, err = c.Read(b1)

	a.NoError(err, "read samples ok")
	a.Equal(len(b1), samples, "correct number of samples read")
	c.Close()

	c, err = NewCaptureDevice("null", 1, FormatS16LE, 22050, BufferParams{})
	a.NoError(err, "created capture device")

	samples, err = c.Read(b2)

	a.NoError(err, "read samples ok")
	a.Equal(len(b2), samples, "correct number of samples read")
	c.Close()

	samples, err = c.Read(b4)

	a.Error(err, "wrong type error")
	a.Equal(samples, 0, "no samples read")
	c.Close()

	// extras floats

	c, err = NewCaptureDevice("null", 1, FormatFloat64LE, 22050, BufferParams{})
	a.NoError(err, "created capture device")

	b5 := make([]float64, 200)
	samples, err = c.Read(b5)

	a.NoError(err, "read samples ok")
	a.Equal(len(b5), samples, "correct number of samples read")
	c.Close()

}

func TestPlayback(t *testing.T) {
	a := assert.New(t)

	p, err := NewPlaybackDevice("nonexistent", 1, FormatS16LE, 44100,
		BufferParams{})

	a.Equal(p, (*PlaybackDevice)(nil), "playback device is nil")
	a.Error(err, "no device error")

	p, err = NewPlaybackDevice("null", 0, FormatS32LE, 44100,
		BufferParams{})

	a.Equal(p, (*PlaybackDevice)(nil), "playback device is nil")
	a.Error(err, "bad channels error")

	p, err = NewPlaybackDevice("null", 1, FormatS32LE, 44100,
		BufferParams{})

	a.NoError(err, "created playback device")

	b1 := make([]int8, 100)
	frames, err := p.Write(b1)

	a.Error(err, "wrong type error")
	a.Equal(frames, 0, "no frames written")

	b2 := make([]int16, 100)
	frames, err = p.Write(b2)

	a.Error(err, "wrong type error")
	a.Equal(frames, 0, "no frames written")

	b3 := make([]float64, 100)
	frames, err = p.Write(b3)

	a.Error(err, "wrong type error")
	a.Equal(frames, 0, "no frames written")

	b4 := make([]int32, 100)
	frames, err = p.Write(b4)

	a.NoError(err, "buffer written ok")
	a.Equal(frames, len(b4), "correct frames written")

	p.Close()

	// extras  s16

	p, err = NewPlaybackDevice("null", 1, FormatS16LE, 44100, BufferParams{})

	a.NoError(err, "created playback device")

	frames, err = p.Write(b4)

	a.Error(err, "wrong type error")
	a.Equal(frames, 0, "no frames written")

	frames, err = p.Write(b2)

	a.NoError(err, "buffer written ok")
	a.Equal(frames, len(b2), "correct frames written")

	p.Close()
}

func TestPanics(t *testing.T) {
	a := assert.New(t)

	// Capture with Format Unknown
	NewCaptureDevice("null", 1, FormatLast+1, 44100, BufferParams{})

	// Playback with Format Unknown
	NewPlaybackDevice("null", 1, FormatLast+1, 44100, BufferParams{})

	// Capture reading with Unknown format must panic()
	shouldPanic(t, func() {
		c, err0 := NewCaptureDevice("null", 1, FormatS8, 44100, BufferParams{})
		a.NoError(err0, "created capture device")

		c.Format = FormatLast
		b1 := make([]int8, 100)
		c.Read(b1)
	})

	// Plaback writing with Unknown format must panic()
	shouldPanic(t, func() {
		p, err0 := NewPlaybackDevice("null", 1, FormatS8, 44100, BufferParams{})
		a.NoError(err0, "created playback device")

		p.Format = FormatLast
		b1 := make([]int8, 100)
		p.Write(b1)
	})
}

func shouldPanic(t *testing.T, f func()) {
	defer func() { recover() }()
	f()
	t.Errorf("should have panicked")
}
