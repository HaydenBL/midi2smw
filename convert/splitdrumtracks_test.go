package convert

import (
	"midi2smw/convert/drumtrack"
	"reflect"
	"testing"
)

func Test_splitDrumTracks(t *testing.T) {
	tracks := []NoteTrack{
		{
			Name: "Bass Drum",
			Notes: []MidiNote{
				{Key: 0},
				{Key: 1},
				{Key: 4},
				{Key: 1},
				{Key: 5},
				{Key: 3},
				{Key: 3},
			},
		},
		{
			Name: "Snare",
			Notes: []MidiNote{
				{Key: 0},
				{Key: 1},
				{Key: 4},
				{Key: 1},
				{Key: 5},
				{Key: 3},
				{Key: 3},
			},
		},
		{
			Name: "Piano",
			Notes: []MidiNote{
				{Key: 0},
				{Key: 1},
				{Key: 4},
				{Key: 1},
				{Key: 5},
				{Key: 3},
				{Key: 3},
			},
		},
		{
			Name: "High Hat",
			Notes: []MidiNote{
				{Key: 0},
				{Key: 1},
				{Key: 4},
				{Key: 1},
				{Key: 5},
				{Key: 3},
				{Key: 3},
			},
		},
	}

	drumTrackGroups := []drumtrack.Group{
		{
			TrackNumber: 1,
			NoteGroups: [][]uint8{
				{0, 3},
			},
		},
		{
			TrackNumber: 0,
			NoteGroups: [][]uint8{
				{0, 3},
				{1},
			},
		},
		{
			TrackNumber: 3,
			NoteGroups: [][]uint8{
				{3},
				{1},
				{0},
			},
		},
	}

	tracks = splitAllTracks(tracks, drumTrackGroups)

	expectedTracks := []NoteTrack{
		{
			Name: "Bass Drum - Split 1",
			Notes: []MidiNote{
				{Key: 0},
				{Key: 3},
				{Key: 3},
			},
		},
		{
			Name: "Bass Drum - Split 2",
			Notes: []MidiNote{
				{Key: 1},
				{Key: 1},
			},
		},
		{
			Name: "Bass Drum - Split 3",
			Notes: []MidiNote{
				{Key: 4},
				{Key: 5},
			},
		},
		{
			Name: "Snare - Split 1",
			Notes: []MidiNote{
				{Key: 0},
				{Key: 3},
				{Key: 3},
			},
		},
		{
			Name: "Snare - Split 2",
			Notes: []MidiNote{
				{Key: 1},
				{Key: 4},
				{Key: 1},
				{Key: 5},
			},
		},

		{
			Name: "Piano",
			Notes: []MidiNote{
				{Key: 0},
				{Key: 1},
				{Key: 4},
				{Key: 1},
				{Key: 5},
				{Key: 3},
				{Key: 3},
			},
		},
		{
			Name: "High Hat - Split 1",
			Notes: []MidiNote{
				{Key: 3},
				{Key: 3},
			},
		},
		{
			Name: "High Hat - Split 2",
			Notes: []MidiNote{
				{Key: 1},
				{Key: 1},
			},
		},
		{
			Name: "High Hat - Split 3",
			Notes: []MidiNote{
				{Key: 0},
			},
		},
		{
			Name: "High Hat - Split 4",
			Notes: []MidiNote{
				{Key: 4},
				{Key: 5},
			},
		},
	}

	if !reflect.DeepEqual(tracks, expectedTracks) {
		t.Fatalf("Split drum track channel not what was expected")
	}
}
