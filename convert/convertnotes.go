package convert

import (
	"fmt"
	"midi2smw/midiparse"
)

type midiNote struct {
	Key       uint8
	Velocity  uint8
	StartTime uint32
	Duration  uint32
}

type noteTrack struct {
	Notes   []midiNote
	MaxNote uint8
	MinNote uint8
}

func convertNotes(tracks []midiparse.MidiTrack) []noteTrack {
	var noteTracks = make([]noteTrack, len(tracks))

	for trackIndex, track := range tracks {
		var wallTime uint32
		var notesBeingProcessed []midiNote

		for _, event := range track.Events {
			wallTime += event.DeltaTick
			if event.Event == midiparse.NoteOn {
				notesBeingProcessed = append(notesBeingProcessed, midiNote{event.Key, event.Velocity, wallTime, 0})
			}
			if event.Event == midiparse.NoteOff {
				i, note := findNoteIndex(notesBeingProcessed, event.Key)
				if i != -1 {
					note.Duration = wallTime - note.StartTime
					noteTracks[trackIndex].Notes = append(noteTracks[trackIndex].Notes, note)
					noteTracks[trackIndex].MinNote = min(noteTracks[trackIndex].MinNote, note.Key)
					noteTracks[trackIndex].MaxNote = max(noteTracks[trackIndex].MaxNote, note.Key)
					notesBeingProcessed = deleteAtIndex(notesBeingProcessed, i)
				}
			}
		}
	}

	noteTracks = filterEmptyNoteTracks(noteTracks)

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

func findNoteIndex(notes []midiNote, key uint8) (int, midiNote) {
	for i, note := range notes {
		if note.Key == key {
			return i, note
		}
	}
	return -1, midiNote{}
}

func deleteAtIndex(notes []midiNote, i int) []midiNote {
	copy(notes[i:], notes[i+1:])
	notes[len(notes)-1] = midiNote{}
	notes = notes[:len(notes)-1]
	return notes
}

func min(a, b uint8) uint8 {
	if a < b {
		return a
	}
	return b
}

func max(a, b uint8) uint8 {
	if a > b {
		return a
	}
	return b
}
