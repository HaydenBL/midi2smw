package convert

import (
	"fmt"
	"midi2smw/midi"
	"sort"
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

func blah(noteTracks []noteTrack) []noteTrack {
	blahed := make([]noteTrack, 0)
	for _, track := range noteTracks {
		track := blahblah(track)
		blahed = append(blahed, track)
	}
	return blahed
}

func blahblah(track noteTrack) noteTrack {
	newNoteTrack := make([]midiNote, 0)
	var sameNotes []midiNote
	chain := false
	for i := range track.Notes {
		if i == len(track.Notes)-1 {
			if chain {
				sameNotes = append(sameNotes, track.Notes[i])
				nthNote := getNthNote(sameNotes)
				if nthNote != nil {
					newNoteTrack = append(newNoteTrack, *nthNote)
				}
			}
			continue
		}
		if track.Notes[i].StartTime == track.Notes[i+1].StartTime {
			chain = true
			sameNotes = append(sameNotes, track.Notes[i])
		} else if chain {
			sameNotes = append(sameNotes, track.Notes[i])

			nthNote := getNthNote(sameNotes)
			if nthNote != nil {
				newNoteTrack = append(newNoteTrack, *nthNote)
			}
			chain = false
			sameNotes = make([]midiNote, 0)
		} else {
			newNoteTrack = append(newNoteTrack, track.Notes[i])
		}
	}
	track.Notes = newNoteTrack
	return track
}

func getNthNote(notes []midiNote) *midiNote {
	const n = 3
	if !(len(notes) > n) {
		return nil
	}
	sort.Slice(notes, func(i, j int) bool {
		return notes[i].Key > notes[j].Key
	})
	return &notes[n]
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
	for _, track := range tracks {
		var filteredEvents []midi.Event
		for _, event := range track.Events {
			if event.Event != midi.Other {
				filteredEvents = append(filteredEvents, event)
			}
		}
		track.Events = filteredEvents
		filteredTracks = append(filteredTracks, track)
	}

	return filteredTracks
}
