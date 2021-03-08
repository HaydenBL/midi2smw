package convert

import (
	"fmt"
	"midi2smw/midi"
)

func Convert(midiTracks []midi.Track) []SmwTrack {
	fmt.Println("Converting midi tracks...")

	var ticksPer64thNote uint32 = 30 // hardcoding for the track I'm working with, figure this out later

	midiTracks = filterOtherEventTypes(midiTracks)
	midiTracks = filterEmptyTracks(midiTracks)

	noteTracks := convertNotes(midiTracks)
	noteTracks = quantizeNotesOnAllTracks(noteTracks, ticksPer64thNote)

	tracks := createSmwChannelTracksForAllTracks(noteTracks, ticksPer64thNote)

	return tracks
}

func filterEmptyTracks(tracks []midi.Track) []midi.Track {
	var nonEmptyTracks []midi.Track
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

func filterOtherEventTypes(tracks []midi.Track) []midi.Track {
	var filteredTracks []midi.Track
	for i, track := range tracks {
		var filteredEvents []midi.Event
		for j, event := range track.Events {
			if event.Event != midi.Other {
				filteredEvents = append(filteredEvents, event)
			} else {
				// if Other event has a delta time, pass that onto the next event so overall time isn't lost
				if j != len(track.Events)-1 {
					tracks[i].Events[j+1].DeltaTick += event.DeltaTick
				}
			}
		}
		track.Events = filteredEvents
		filteredTracks = append(filteredTracks, track)
	}

	return filteredTracks
}
