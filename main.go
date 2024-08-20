package main

import (
	"flag"
	"log"
	"math"
	"os"
	"time"

	"github.com/hajimehoshi/oto/v2"
)

func main() {
	fs := flag.NewFlagSet("metronome", flag.ExitOnError)
	bpm := fs.Float64("bpm", 90.0, "beats per minute")
	fs.Parse(os.Args[1:])

	beep, sampleRate := makeBeep()

	context, ready, err := oto.NewContext(sampleRate, 1, oto.FormatSignedInt16LE)
	if err != nil {
		log.Fatal(err)
	}
	<-ready

	if err := context.Err(); err != nil {
		log.Fatal(err)
	}

	interval := time.Duration(60.0 / *bpm * float64(time.Second))

	r := newIntervalReader(beep, interval)

	p := context.NewPlayer(r)
	defer p.Close()

	p.(oto.BufferSizeSetter).SetBufferSize(len(beep))

	p.Play()

	select {}
}

func makeBeep() ([]byte, int) {
	const (
		sampleRate = 44100
		frequency  = 1000.0
		duration   = 0.1 // duration in seconds
	)

	// Generate a 16-bit sine wave
	numSamples := int(sampleRate * duration)
	buffer := make([]byte, numSamples*2) // 2 bytes per sample for 16-bit audio
	amplitude := 32767.0                 // maximum amplitude for 16-bit audio

	for i := 0; i < numSamples; i++ {
		t := float64(i) / sampleRate
		// Generate sine wave value
		sample := int16(amplitude * math.Sin(2*math.Pi*frequency*t))
		buffer[2*i] = byte(sample & 0xff)
		buffer[2*i+1] = byte(sample >> 8)
	}

	return buffer, sampleRate
}

type intervalReader struct {
	buffer  []byte
	initial chan struct{}
	tick    *time.Ticker
	pos     int
}

func newIntervalReader(buffer []byte, interval time.Duration) *intervalReader {
	ticker := time.NewTicker(interval)

	initial := make(chan struct{})
	close(initial)
	return &intervalReader{
		buffer:  buffer,
		initial: initial,
		tick:    ticker,
	}
}

func (r *intervalReader) Read(p []byte) (int, error) {
	if r.pos == 0 {
		select {
		case <-r.initial:
			r.initial = nil
		case <-r.tick.C:
		}
	}
	n := copy(p, r.buffer[r.pos:])
	r.pos += n
	if r.pos >= len(r.buffer) {
		r.pos = 0
	}
	return n, nil
}
