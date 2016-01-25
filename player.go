package main

import (
	"io"
	"os"

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

const seekRate = 4096

func PlayPath(p string) error {
	var song codec.Song
	songs, _, err := codec.ByExtension(p, fileReader(p))
	if err != nil {
		return err
	}
	for _, song = range songs {
		// Get first song
		break
	}
	sampleRate, channels, err := song.Init()
	if err != nil {
		return err
	}
	defer song.Close()

	out, err := output.Get(sampleRate, channels)
	if err != nil {
		return err
	}

	for {
		samples, err := song.Play(seekRate)
		if err != nil {
			return err
		}
		out.Push(samples)
		if len(samples) < seekRate {
			// Done
			break
		}
	}

	return nil
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
