package hls

import (
	"bytes"
	"embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed test_data/*
var f embed.FS

func TestParser(t *testing.T) {
	type testCase struct {
		input    string
		expected *Manifest
	}

	for name, test := range map[string]testCase{
		"vod playlist": {
			input: "test_data/vod_playlist.m3u8",
			expected: &Manifest{
				PlaylistType:   "VOD",
				TargetDuration: 10,
				Version:        4,
				MediaSequence:  0,
				Segments: []*Segment{
					{
						Duration: 10.0,
						URL:      "http://example.com/movie1/fileSequenceA.ts",
					},
					{
						Duration: 10.0,
						URL:      "http://example.com/movie1/fileSequenceB.ts",
					},
					{
						Duration: 10.0,
						URL:      "http://example.com/movie1/fileSequenceC.ts",
					},
					{
						Duration: 9.0,
						URL:      "http://example.com/movie1/fileSequenceD.ts",
					},
				},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			data, readErr := f.ReadFile(test.input)
			assert.NoError(t, readErr)

			reader := bytes.NewReader(data)

			parser := NewParser()
			m, parseErr := parser.Parse(reader)

			assert.NoError(t, parseErr)
			assert.Equal(t, test.expected, m)
		})
	}
}
