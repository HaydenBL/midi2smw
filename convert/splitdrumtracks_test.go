package convert

import (
	"fmt"
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
	}

	midiTracks := []MidiTrackWithNoteGroups{
		{NoteGroups: []NoteGroup{
			{[]uint8{0, 3}},
			{[]uint8{1}},
		}},
		{NoteGroups: []NoteGroup{
			{[]uint8{0, 3}},
		}},
		{NoteGroups: []NoteGroup{
			{[]uint8{3}},
			{[]uint8{1}},
			{[]uint8{0}},
		}},
		{NoteGroups: []NoteGroup{}},
	}
	tracks = splitAllTracks(tracks, midiTracks)

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
	}

	if !tracksAreEqual(tracks, expectedTracks) {
		t.Fatalf("Split drum track channel not what was expected")
	}
}

func tracksAreEqual(track1, track2 []NoteTrack) bool {
	if len(track1) != len(track2) {
		fmt.Printf("Lengths of not tracks not equal: %d vs %d\n", len(track1), len(track2))
		return false
	}
	for i := range track1 {
		if track1[i].Name != track2[i].Name {
			fmt.Printf("Track names not equal: %s vs %s\n", track1[i].Name, track2[i].Name)
			return false
		}
		if len(track1[i].Notes) != len(track2[i].Notes) {
			fmt.Printf("Track note lengths not the same: %d vs %d\n", len(track1[i].Notes), len(track2[i].Notes))
			return false
		}
		for noteIndex := range track1[i].Notes {
			if track1[i].Notes[noteIndex] != track2[i].Notes[noteIndex] {
				fmt.Printf("Note indexes note the same: %d vs %d", track1[i].Notes[noteIndex], track2[i].Notes[noteIndex])
				return false
			}
		}
	}
	return true
}
