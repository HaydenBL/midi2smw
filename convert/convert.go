package convert

import (
	"fmt"
	"midi2smw/midiparse"
)

type MidiNote struct {
	Key       uint8
	Velocity  uint8
	StartTime uint32
	Duration  uint32
}

func Convert(midiTracks []midiparse.MidiTrack) {
	fmt.Println("Converting midi tracks...")

	//midiTracks = filterOtherEventTypes(midiTracks)
	//midiTracks = filterEmptyTracks(midiTracks)

	fmt.Printf("Tracks with note data: %d\n", len(midiTracks))

	noteTracks := convertNotes(midiTracks)

	fmt.Printf("bah, %d", noteTracks[0].MaxNote)

}

func filterEmptyTracks(tracks []midiparse.MidiTrack) []midiparse.MidiTrack {
	var nonEmptyTracks []midiparse.MidiTrack
	for _, track := range tracks {
		if len(track.Events) != 0 {
			nonEmptyTracks = append(nonEmptyTracks, track)
		}
	}

	fmt.Printf("Removed %d tracks without note data\n", len(tracks)-len(nonEmptyTracks))
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
