package convert

import (
	"fmt"
	"midi2smw/midiparse"
)

func Convert(midiTracks []midiparse.MidiTrack) {
	fmt.Println("Converting midi tracks...")

	midiTracks = filterOtherEventTypes(midiTracks)
	midiTracks = filterEmptyTracks(midiTracks)

	noteTracks := convertNotes(midiTracks)
	tracks := createSmwChannelTracks(noteTracks)

	testPrint(tracks[0])
}

// temporary, just to test this thing
func testPrint(smwTrack []SmwNote) {
	lastOctave := smwTrack[0].octave
	fmt.Printf("Start octave: %d\n", lastOctave)
	for _, smwNote := range smwTrack {
		if smwNote.octave > lastOctave {
			fmt.Printf(">")
		} else if smwNote.octave < lastOctave {
			fmt.Printf("<")
		}
		for i, note := range smwNote.lengthValues {
			if i == 0 {
				fmt.Printf("%s%d", smwNote.key, note)
			} else {
				fmt.Printf("^%d", note)
			}
		}
		lastOctave = smwNote.octave
	}
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
