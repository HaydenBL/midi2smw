package convert

import (
	"midi2smw/drumtrack"
	"reflect"
	"testing"
)

func Test_splitDrumTracks(t *testing.T) {
	tracks := []noteTrack{
		{
			Name: "Bass Drum",
			Notes: []midiNote{
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
			Notes: []midiNote{
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
			Notes: []midiNote{
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
			Notes: []midiNote{
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

	tracks = splitDrumTracks(tracks, drumTrackGroups)

	expectedTracks := []noteTrack{
		{
			Name: "Bass Drum - Split 1",
			Notes: []midiNote{
				{Key: 0},
				{Key: 3},
				{Key: 3},
			},
		},
		{
			Name: "Bass Drum - Split 2",
			Notes: []midiNote{
				{Key: 1},
				{Key: 1},
			},
		},
		{
			Name: "Bass Drum - Split 3",
			Notes: []midiNote{
				{Key: 4},
				{Key: 5},
			},
		},
		{
			Name: "Snare - Split 1",
			Notes: []midiNote{
				{Key: 0},
				{Key: 3},
				{Key: 3},
			},
		},
		{
			Name: "Snare - Split 2",
			Notes: []midiNote{
				{Key: 1},
				{Key: 4},
				{Key: 1},
				{Key: 5},
			},
		},

		{
			Name: "Piano",
			Notes: []midiNote{
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
			Notes: []midiNote{
				{Key: 3},
				{Key: 3},
			},
		},
		{
			Name: "High Hat - Split 2",
			Notes: []midiNote{
				{Key: 1},
				{Key: 1},
			},
		},
		{
			Name: "High Hat - Split 3",
			Notes: []midiNote{
				{Key: 0},
			},
		},
		{
			Name: "High Hat - Split 4",
			Notes: []midiNote{
				{Key: 4},
				{Key: 5},
			},
		},
	}

	if !reflect.DeepEqual(tracks, expectedTracks) {
		t.Fatalf("Split drum track channel not what was expected")
	}
}
