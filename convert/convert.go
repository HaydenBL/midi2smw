package convert

import (
	"fmt"
	"midi2smw/midiparse"
)

func Convert(midiTracks []midiparse.MidiTrack) [][]SmwNote {
	fmt.Println("Converting midi tracks...")

	var ticksPer64thNote uint32 = 30 // hardcoding for the track I'm working with, figure this out later

	midiTracks = filterOtherEventTypes(midiTracks)
	midiTracks = filterEmptyTracks(midiTracks)

	noteTracks := convertNotes(midiTracks)
	noteTracks = quantizeNotesOnAllTracks(noteTracks, ticksPer64thNote)

	tracks := createSmwChannelTracksForAllTracks(noteTracks, ticksPer64thNote)

	return tracks
}

func filterEmptyTracks(tracks []midiparse.MidiTrack) []midiparse.MidiTrack {
	var nonEmptyTracks []midiparse.MidiTrack
	for _, track := range tracks {
		if len(track.Events) != 0 {
			nonEmptyTracks = append(nonEmptyTracks, track)
		}
	}

	if len(tracks) > len(nonEmptyTracks) {
		fmt.Printf("Removed %d tracks with no midi event data\n", len(tracks)-len(nonEmptyTracks))
	}
	return nonEmptyTracks
}

func filterOtherEventTypes(tracks []midiparse.MidiTrack) []midiparse.MidiTrack {
	var filteredTracks []midiparse.MidiTrack
	for _, track := range tracks {
		var filteredEvents []midiparse.MidiEvent
		for _, event := range track.Events {
			if event.Event != midiparse.Other {
				filteredEvents = append(filteredEvents, event)
			}
		}
		track.Events = filteredEvents
		filteredTracks = append(filteredTracks, track)
	}

	return filteredTracks
}
