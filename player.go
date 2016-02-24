package main

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/mjibson/go-dsp/spectral"
	"github.com/mjibson/mog/codec"
	"github.com/mjibson/mog/output"

	// codecs
	_ "github.com/mjibson/mog/codec/flac"
	_ "github.com/mjibson/mog/codec/gme"
	_ "github.com/mjibson/mog/codec/mpa"
	_ "github.com/mjibson/mog/codec/nsf"
	_ "github.com/mjibson/mog/codec/rar"
	_ "github.com/mjibson/mog/codec/vorbis"
	_ "github.com/mjibson/mog/codec/wav"
)

type player struct {
	output output.Output
	song   codec.Song

	sampleRate int
	channels   int

	done chan struct{}

	mu      sync.RWMutex
	samples []float64
}

func (p *player) Start() {
	seekRate := int(p.sampleRate / 10.0) // fps

	po := &spectral.PwelchOptions{
		NFFT:      8,
		Scale_off: true,
	}

	go func() {
		for {
			select {
			case <-p.done:
				return
			default:
			}
			samples, err := p.song.Play(seekRate)
			if err != nil {
				// TODO: Something about this err?
				return
			}
			p.output.Push(samples)

			p.mu.Lock()
			p.samples, _ = spectral.Pwelch(float32To64(samples), float64(seekRate), po)
			p.mu.Unlock()

			if len(samples) < seekRate {
				// Done
				return
			}
		}
	}()
}

func (p *player) Sample() []float64 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.samples
}

func (p *player) Stop() {
	p.done <- struct{}{}
}

func (p *player) Close() {
	p.song.Close()
	p.output.Stop()
}

func SongPlayer(path string) (*player, error) {
	var song codec.Song

	songs, _, err := codec.ByExtension(path, fileReader(path))
	if err != nil {
		return nil, err
	}
	if len(songs) == 0 {
		return nil, fmt.Errorf("no songs detected")
	}
	for _, song = range songs {
		// Get first song
		break
	}
	sampleRate, channels, err := song.Init()
	if err != nil {
		return nil, err
	}

	out, err := output.Get(sampleRate, channels)
	if err != nil {
		return nil, err
	}
	out.Start()

	p := &player{
		output:     out,
		song:       song,
		sampleRate: sampleRate,
		channels:   channels,
		done:       make(chan struct{}),

		mu: sync.RWMutex{},
	}

	return p, nil
}

func fileReader(path string) codec.Reader {
	return func() (io.ReadCloser, int64, error) {
		f, err := os.Open(path)
		if err != nil {
			return nil, 0, err
		}
		fi, err := f.Stat()
		if err != nil {
			f.Close()
			return nil, 0, err
		}
		return f, fi.Size(), nil
	}
}
