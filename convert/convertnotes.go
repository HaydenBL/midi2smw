package convert

import (
	"fmt"
	"midi2smw/drumtrack"
	"midi2smw/midi"
)

type midiNote struct {
	Key       uint8
	Velocity  uint8
	StartTime uint32
	Duration  uint32
}

type noteTrack struct {
	Name    string
	Notes   []midiNote
	MaxNote uint8
	MinNote uint8
}

func convertNotes(tracks []midi.Track, drumTrackGroups []drumtrack.Group) []noteTrack {
	var noteTracks = make([]noteTrack, len(tracks))

	for trackIndex, track := range tracks {
		var wallTime uint32
		var notesBeingProcessed []midiNote

		for _, event := range track.Events {
			wallTime += event.DeltaTick
			if event.Event == midi.NoteOn {
				notesBeingProcessed = append(notesBeingProcessed, midiNote{event.Key, event.Velocity, wallTime, 0})
			}
			if event.Event == midi.NoteOff {
				i, note := findNoteIndex(notesBeingProcessed, event.Key)
				if i != -1 {
					note.Duration = wallTime - note.StartTime
					noteTracks[trackIndex].Name = track.Name
					noteTracks[trackIndex].Notes = append(noteTracks[trackIndex].Notes, note)
					noteTracks[trackIndex].MinNote = minUint8(noteTracks[trackIndex].MinNote, note.Key)
					noteTracks[trackIndex].MaxNote = maxUint8(noteTracks[trackIndex].MaxNote, note.Key)
					notesBeingProcessed = deleteAtIndex(notesBeingProcessed, i)
				}
			}
		}
	}

	noteTracks = filterEmptyNoteTracks(noteTracks)
	if len(drumTrackGroups) > 0 {
		noteTracks = splitDrumTracks(noteTracks, drumTrackGroups)
	}

	return noteTracks
}

func filterEmptyNoteTracks(tracks []noteTrack) []noteTrack {
	var nonEmptyTracks []noteTrack
	for _, track := range tracks {
		if len(track.Notes) != 0 {
			nonEmptyTracks = append(nonEmptyTracks, track)
		}
	}

	if len(nonEmptyTracks) < len(tracks) {
		fmt.Printf("Removed %d tracks with no note data\n", len(tracks)-len(nonEmptyTracks))
	}
	return nonEmptyTracks
}
